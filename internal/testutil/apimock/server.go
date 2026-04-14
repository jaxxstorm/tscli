package apimock

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

type Request struct {
	Method string
	Path   string
	Query  string
	Body   string
	Header http.Header
}

type Route struct {
	Method       string
	PathContains string
	Status       int
	Body         string
}

type Server struct {
	t      *testing.T
	Server *httptest.Server

	mu       sync.Mutex
	routes   []Route
	requests []Request
}

func New(t *testing.T) *Server {
	t.Helper()

	s := &Server{t: t}
	s.Server = httptest.NewServer(http.HandlerFunc(s.handle))
	t.Cleanup(func() {
		s.Server.Close()
	})
	return s
}

func (s *Server) URL() string {
	return s.Server.URL
}

func (s *Server) AddJSON(method, pathContains string, status int, payload any) {
	s.t.Helper()
	b, err := json.Marshal(payload)
	if err != nil {
		s.t.Fatalf("marshal payload: %v", err)
	}
	s.AddRaw(method, pathContains, status, string(b))
}

func (s *Server) AddRaw(method, pathContains string, status int, body string) {
	s.t.Helper()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.routes = append(s.routes, Route{
		Method:       strings.ToUpper(method),
		PathContains: pathContains,
		Status:       status,
		Body:         body,
	})
}

func (s *Server) Requests() []Request {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Request, len(s.requests))
	copy(out, s.requests)
	return out
}

func (s *Server) handle(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	_ = r.Body.Close()

	s.mu.Lock()
	s.requests = append(s.requests, Request{
		Method: r.Method,
		Path:   r.URL.Path,
		Query:  r.URL.RawQuery,
		Body:   string(body),
		Header: r.Header.Clone(),
	})
	routes := make([]Route, len(s.routes))
	copy(routes, s.routes)
	s.mu.Unlock()

	for _, route := range routes {
		if route.Method != strings.ToUpper(r.Method) {
			continue
		}
		if route.PathContains != "" && !strings.Contains(r.URL.Path, route.PathContains) {
			continue
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(route.Status)
		_, _ = io.WriteString(w, route.Body)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	_, _ = io.WriteString(w, fmt.Sprintf(`{"message":"no mock route for %s %s"}`, r.Method, r.URL.Path))
}
