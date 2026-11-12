package oidctest

import (
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/babakapaa/lodestar/server/internal/platform/demoidp"
)

type Server struct {
	URL          string
	ClientID     string
	ClientSecret string

	mu   sync.Mutex
	user User
}

func New(t *testing.T) *Server {
	t.Helper()
	s := &Server{
		ClientID:     "lodestar",
		ClientSecret: "test-secret",
		user:         User{Subject: "student-1", Email: "student@uni.example", Name: "Айгерим"},
	}
	provider, err := demoidp.New(demoidp.Config{
		ClientID:     s.ClientID,
		ClientSecret: s.ClientSecret,
		Pick:         s.current,
	})
	if err != nil {
		t.Fatalf("demo provider: %v", err)
	}
	srv := httptest.NewServer(provider.Handler())
	t.Cleanup(srv.Close)
	provider.SetIssuer(srv.URL)
	s.URL = srv.URL
	return s
}

func (s *Server) SignIn(u User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.user = u
}

func (s *Server) current() demoidp.User {
	s.mu.Lock()
	defer s.mu.Unlock()
	return demoidp.User{Subject: s.user.Subject, Email: s.user.Email, Name: s.user.Name}
}
