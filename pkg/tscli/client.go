// pkg/tscli/client.go
//
// Thin wrapper around tailscale-client-go that:
//
//   - picks up tailnet / api-key / debug from Viper
//   - logs every HTTP request & response when --debug or TSCLI_DEBUG=1 is set
package tscli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jaxxstorm/tscli/pkg/oauth"
	"github.com/jaxxstorm/tscli/pkg/version"
	"github.com/spf13/viper"
	tsapi "tailscale.com/client/tailscale/v2"
)

const (
	defaultBaseURL     = "https://api.tailscale.com"
	defaultContentType = "application/json"
)

// getUserAgent returns the default user agent string, derived from the
// build version (or, if unset, this process's local git repository state).
// Callers that want a fixed value — e.g. a library consumer who should not
// leak their own repo's git metadata — should set the "user-agent" viper
// key (flag/env/config) instead of relying on this default.
func getUserAgent() string {
	return fmt.Sprintf("tscli/%s (Go client)", version.GetVersion())
}

func configuredUserAgent() string {
	userAgent := viper.GetString("user-agent")
	if userAgent == "" {
		userAgent = getUserAgent()
	}
	return userAgent
}

func New() (*tsapi.Client, error) {
	tailnet := viper.GetString("tailnet")
	apiKey := viper.GetString("api-key")
	oauthClientID := viper.GetString("oauth-client-id")
	oauthClientSecret := viper.GetString("oauth-client-secret")
	baseURL := viper.GetString("base-url")
	if tailnet == "" {
		return nil, fmt.Errorf("tailnet is required")
	}
	if apiKey == "" && (oauthClientID == "" || oauthClientSecret == "") {
		return nil, fmt.Errorf("either api-key or both oauth-client-id and oauth-client-secret are required")
	}

	userAgent := configuredUserAgent()

	// Create a custom transport that ensures UserAgent is always set
	transport := &userAgentTransport{
		rt:        http.DefaultTransport,
		userAgent: userAgent,
	}

	// Wrap with debug logging if enabled
	if viper.GetBool("debug") {
		transport.rt = &logTransport{rt: transport.rt}
	}

	httpClient := &http.Client{
		Transport: transport,
		// Tailscale's API never legitimately redirects. Refusing to follow
		// redirects prevents oauthBearerTransport's Authorization header
		// (re-attached on every RoundTrip call, including the request the
		// stdlib http.Client builds for a followed redirect) from ever being
		// sent to a host other than the one this client was configured for.
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	client := &tsapi.Client{
		Tailnet:   tailnet,
		UserAgent: userAgent,
		HTTP:      httpClient,
	}
	resolvedBaseURL, err := resolveConfiguredBaseURL(baseURL)
	if err != nil {
		return nil, err
	}
	if baseURL != "" {
		client.BaseURL = resolvedBaseURL
	}

	if apiKey != "" {
		client.APIKey = apiKey
	} else {
		tokenURL, err := oauthTokenURL(resolvedBaseURL)
		if err != nil {
			return nil, err
		}
		httpClient.Transport = &oauthBearerTransport{
			rt:           httpClient.Transport,
			clientID:     oauthClientID,
			clientSecret: oauthClientSecret,
			tokenURL:     tokenURL,
		}
	}

	return client, nil
}

// ExchangeOAuthClientCredentials exchanges OAuth client credentials via
// oauth.ExchangeClientCredentials, passing the token endpoint derived from
// the same tailnet/base-url configuration New() reads from viper (flag,
// environment, or config file) as the default. TSCLI_OAUTH_TOKEN_URL still
// overrides that derived default; see oauth.ExchangeClientCredentials for
// the validation applied to both.
func ExchangeOAuthClientCredentials(ctx context.Context, clientID, clientSecret string) (*oauth.TokenResponse, error) {
	baseURL, err := resolveBaseURL(nil)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return oauth.ExchangeClientCredentials(ctx, clientID, clientSecret, joinAPIV2Path(baseURL, "/oauth/token").String())
}

// oauthTokenURL returns the OAuth client-credentials token endpoint: the
// TSCLI_OAUTH_TOKEN_URL environment variable if it is set and passes
// oauth.ValidateTokenURL, otherwise baseURL's "/api/v2/oauth/token" path
// (preserving any path prefix baseURL has, e.g. behind a reverse proxy).
func oauthTokenURL(baseURL *url.URL) (string, error) {
	return oauth.ResolveTokenURL(joinAPIV2Path(baseURL, "/oauth/token").String())
}

type oauthBearerTransport struct {
	rt           http.RoundTripper
	clientID     string
	clientSecret string
	tokenURL     string

	mu          sync.Mutex
	accessToken string
	expiresAt   time.Time
}

func (t *oauthBearerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	accessToken, err := t.token(req.Context())
	if err != nil {
		return nil, err
	}
	clone := req.Clone(req.Context())
	clone.Header = req.Header.Clone()
	clone.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := t.rt.RoundTrip(clone)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusUnauthorized {
		return resp, nil
	}
	resp.Body.Close()
	if err := t.invalidateToken(); err != nil {
		return nil, err
	}
	accessToken, err = t.token(req.Context())
	if err != nil {
		return nil, err
	}
	retry := req.Clone(req.Context())
	retry.Header = req.Header.Clone()
	retry.Header.Set("Authorization", "Bearer "+accessToken)
	return t.rt.RoundTrip(retry)
}

func (t *oauthBearerTransport) token(ctx context.Context) (string, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.accessToken != "" && (t.expiresAt.IsZero() || time.Now().Before(t.expiresAt)) {
		return t.accessToken, nil
	}
	resp, err := oauth.ExchangeClientCredentialsAtURL(ctx, t.clientID, t.clientSecret, t.tokenURL)
	if err != nil {
		return "", fmt.Errorf("exchange oauth credentials: %w", err)
	}
	t.accessToken = resp.AccessToken
	t.expiresAt = time.Time{}
	if !resp.ExpiresAt.IsZero() {
		t.expiresAt = resp.ExpiresAt.Add(-30 * time.Second)
	}
	return t.accessToken, nil
}

func (t *oauthBearerTransport) invalidateToken() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.accessToken = ""
	t.expiresAt = time.Time{}
	return nil
}

// Do performs an HTTP call on top of an existing *tsapi.Client.  Useful for
// endpoints not yet covered by the SDK.  When “debug” is on, full request /
// response dumps are printed to stderr.
func Do(
	ctx context.Context,
	c *tsapi.Client,
	method, path string,
	body any,
	out any,
) (http.Header, error) {
	base, err := resolveBaseURL(c.BaseURL)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}

	u.Path = strings.ReplaceAll(u.Path, "{tailnet}", url.PathEscape(c.Tailnet))

	full := joinAPIV2Path(base, u.Path)
	full.RawQuery = u.RawQuery

	var rdr io.Reader
	if body != nil {
		switch v := body.(type) {
		case []byte:
			rdr = bytes.NewReader(v)
		case string:
			rdr = strings.NewReader(v)
		default:
			b, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("marshal body: %w", err)
			}
			rdr = bytes.NewReader(b)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, full.String(), rdr)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Accept", defaultContentType)
	if body != nil {
		req.Header.Set("Content-Type", defaultContentType)
	}
	if c.APIKey != "" {
		req.SetBasicAuth(c.APIKey, "")
	}

	return doRequest(c.HTTP, req, method, path, out)
}

func DoBearer(
	ctx context.Context,
	method, path string,
	accessToken string,
	body any,
	out any,
) (http.Header, error) {
	base, err := resolveBaseURL(nil)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}

	full := joinAPIV2Path(base, u.Path)
	full.RawQuery = u.RawQuery

	var rdr io.Reader
	if body != nil {
		switch v := body.(type) {
		case []byte:
			rdr = bytes.NewReader(v)
		case string:
			rdr = strings.NewReader(v)
		default:
			b, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("marshal body: %w", err)
			}
			rdr = bytes.NewReader(b)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, full.String(), rdr)
	if err != nil {
		return nil, err
	}
	userAgent := configuredUserAgent()
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", defaultContentType)
	if body != nil {
		req.Header.Set("Content-Type", defaultContentType)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	transport := &userAgentTransport{
		rt:        http.DefaultTransport,
		userAgent: userAgent,
	}

	return doRequest(&http.Client{Transport: transport}, req, method, path, out)
}

func resolveBaseURL(current *url.URL) (*url.URL, error) {
	if current != nil {
		if err := validateBaseURL(current); err != nil {
			return nil, err
		}
		return current, nil
	}

	return resolveConfiguredBaseURL(viper.GetString("base-url"))
}

func resolveConfiguredBaseURL(baseURL string) (*url.URL, error) {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return parseBaseURL(baseURL)
}

// joinAPIV2Path appends endpointPath below the API v2 root. A configured
// base URL may either identify a proxy prefix or the API root itself, as in
// the server URL published by the OpenAPI document.
func joinAPIV2Path(baseURL *url.URL, endpointPath string) *url.URL {
	apiBaseURL := baseURL
	if !strings.HasSuffix(strings.TrimRight(baseURL.Path, "/"), "/api/v2") {
		apiBaseURL = baseURL.JoinPath("api/v2")
	}
	return apiBaseURL.JoinPath(endpointPath)
}

func parseBaseURL(raw string) (*url.URL, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("invalid base-url: %w", err)
	}
	if err := validateBaseURL(parsed); err != nil {
		return nil, err
	}
	return parsed, nil
}

func validateBaseURL(u *url.URL) error {
	if _, err := oauth.ValidateTokenURL(u.String()); err != nil {
		return fmt.Errorf("invalid base-url: %w", err)
	}
	return nil
}

func doRequest(httpc *http.Client, req *http.Request, method string, path string, out any) (http.Header, error) {

	// dump request information if debug is enabled
	if viper.GetBool("debug") {
		if dump, _ := httputil.DumpRequestOut(req, true); len(dump) > 0 {
			os.Stderr.Write(dump)
		}
	}

	if httpc == nil {
		httpc = http.DefaultClient
	}

	res, err := httpc.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return res.Header, err
	}

	// dump response information if debug is enabled
	if viper.GetBool("debug") {
		if dump, _ := httputil.DumpResponse(res, false); len(dump) > 0 {
			os.Stderr.Write(dump)
		}
		if len(data) < 4_096 {
			os.Stderr.Write(data)
			fmt.Fprintln(os.Stderr)
		}
	}

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return res.Header, fmt.Errorf(
			"tailscale API %s %s -> %d: %s",
			method, path, res.StatusCode, strings.TrimSpace(string(data)),
		)
	}

	if out == nil || len(data) == 0 {
		return res.Header, nil
	}
	if raw, ok := out.(*[]byte); ok {
		*raw = append((*raw)[:0], data...)
		return res.Header, nil
	}
	if raw, ok := out.(*json.RawMessage); ok {
		*raw = append((*raw)[:0], data...)
		return res.Header, nil
	}
	if err := json.Unmarshal(data, out); err != nil {
		return res.Header, fmt.Errorf("decode response: %w", err)
	}
	return res.Header, nil
}

type logTransport struct{ rt http.RoundTripper }

func (t *logTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if dump, _ := httputil.DumpRequestOut(req, true); len(dump) > 0 {
		os.Stderr.Write(dump)
	}
	resp, err := t.rt.RoundTrip(req)
	if err != nil {
		return resp, err
	}
	if dump, _ := httputil.DumpResponse(resp, false); len(dump) > 0 {
		os.Stderr.Write(dump)
	}
	return resp, nil
}

// userAgentTransport wraps an http.RoundTripper to ensure UserAgent is always set
type userAgentTransport struct {
	rt        http.RoundTripper
	userAgent string
}

func (t *userAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Always set our custom user agent
	req.Header.Set("User-Agent", t.userAgent)
	return t.rt.RoundTrip(req)
}
