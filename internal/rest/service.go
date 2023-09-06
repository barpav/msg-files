package rest

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/barpav/msg-files/internal/rest/models"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type Service struct {
	Shutdown chan struct{}
	cfg      *config
	server   *http.Server
	auth     Authenticator
	storage  Storage
}

type Authenticator interface {
	ValidateSession(ctx context.Context, key, ip, agent string) (userId string, err error)
}

//go:generate mockery --name Storage
type Storage interface {
	AllocateNewFile(ctx context.Context, info *models.AllocatedFile) (id string, err error)
	AllocatedFileInfo(ctx context.Context, id string) (info *models.AllocatedFile, err error)
	UploadFileContent(id string, content io.Reader) error
	FileSize(ctx context.Context, id string) (size int, err error)
	DownloadFile(id string, stream io.Writer) error
	DeleteFile(ctx context.Context, id string) error
	MarkAsUnused(ctx context.Context, fileId string) error
}

func (s *Service) Start(auth Authenticator, storage Storage) {
	s.cfg = &config{}
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
	err = s.server.Shutdown(ctx)

	if err != nil {
		err = fmt.Errorf("failed to stop HTTP service: %w", err)
	}

	return err
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
