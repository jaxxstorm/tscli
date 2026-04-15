package oauth

import (
	"context"
	"fmt"
	"os"
	"time"

	"golang.org/x/oauth2/clientcredentials"
)

// TokenResponse represents the response from the OAuth token exchange
type TokenResponse struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresIn   int       `json:"expires_in"`
	ExpiresAt   time.Time `json:"expires_at"`
}

// ExchangeClientCredentials exchanges OAuth client credentials for an access token
func ExchangeClientCredentials(ctx context.Context, clientID, clientSecret string) (*TokenResponse, error) {
	tokenURL := os.Getenv("TSCLI_OAUTH_TOKEN_URL")
	if tokenURL == "" {
		tokenURL = "https://api.tailscale.com/api/v2/oauth/token"
	}
	return ExchangeClientCredentialsAtURL(ctx, clientID, clientSecret, tokenURL)
}

func ExchangeClientCredentialsAtURL(ctx context.Context, clientID, clientSecret, tokenURL string) (*TokenResponse, error) {

	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     tokenURL,
	}

	token, err := config.Token(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange credentials: %w", err)
	}

	// Convert oauth2.Token to our TokenResponse format
	expiresIn := 0
	expiresAt := time.Time{}
	if !token.Expiry.IsZero() {
		expiresIn = int(time.Until(token.Expiry).Seconds())
		// Ensure we don't return negative values if token is already expired
		if expiresIn < 0 {
			expiresIn = 0
		}
		expiresAt = token.Expiry
	}

	return &TokenResponse{
		AccessToken: token.AccessToken,
		TokenType:   token.TokenType,
		ExpiresIn:   expiresIn,
		ExpiresAt:   expiresAt,
	}, nil
}
