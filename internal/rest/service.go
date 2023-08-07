package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type Service struct {
	Shutdown chan struct{}
	cfg      *Config
	server   *http.Server
	auth     Authenticator
	storage  Storage
}

type Authenticator interface {
	ValidateSession(ctx context.Context, key, ip, agent string) (userId string, err error)
}

type Storage interface {
}

func (s *Service) Start(auth Authenticator, storage Storage) {
	s.cfg = &Config{}
	s.cfg.Read()

	s.auth, s.storage = auth, storage

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%s", s.cfg.port),
		Handler: s.operations(),
	}

	s.Shutdown = make(chan struct{}, 1)

	go func() {
		err := s.server.ListenAndServe()

		if err != http.ErrServerClosed {
			log.Err(err).Msg("HTTP server crashed.")
		}

		s.Shutdown <- struct{}{}
	}()
}

func (s *Service) Stop(ctx context.Context) (err error) {
	return s.server.Shutdown(ctx)
}

// Specification: https://barpav.github.io/msg-api-spec/#/files
func (s *Service) operations() *chi.Mux {
	ops := chi.NewRouter()

	ops.Use(s.traceInternalServerError)
	ops.Use(s.authenticate)

	// Public endpoint is the concern of the api gateway
	ops.Post("/", s.allocateNewFile)
	ops.Post("/{id}", s.uploadNewFileContent)
	ops.Get("/{id}", s.getFile)
	ops.Head("/{id}", s.getFile)
	ops.Delete("/{id}", s.deleteFile)

	return ops
}
