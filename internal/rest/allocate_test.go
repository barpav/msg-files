package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/barpav/msg-files/internal/rest/mocks"
	"github.com/barpav/msg-files/internal/rest/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_allocateNewFile(t *testing.T) {
	type testService struct {
		storage Storage
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name        string
		testService testService
		args        args
		wantHeaders map[string]string
		wantStatus  int
	}{
		{
			name: "New private file allocated successfully (201)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					fileDesc := models.NewPrivateFileV1{
						Name:   "test.jpg",
						Mime:   "image/jpg",
						Access: []string{"john", "bob", "alice"},
					}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(fileDesc)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("POST", "/", &buf)
					r.Header.Set("Content-Type", "application/vnd.newPrivateFile.v1+json")
					return r
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("AllocateNewFile", mock.Anything, mock.Anything).Return("test-file-id", nil)
					s.On("MarkAsUnused", mock.Anything, "test-file-id").Return(nil)
					return s
				}(),
			},
			wantHeaders: map[string]string{
				"Location": "/test-file-id",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "New public file allocated successfully (201)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					fileDesc := models.NewPublicFileV1{
						Name: "test.jpg",
						Mime: "image/jpg",
					}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(fileDesc)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("POST", "/", &buf)
					r.Header.Set("Content-Type", "application/vnd.newPublicFile.v1+json")
					return r
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("AllocateNewFile", mock.Anything, mock.Anything).Return("test-file-id", nil)
					s.On("MarkAsUnused", mock.Anything, "test-file-id").Return(nil)
					return s
				}(),
			},
			wantHeaders: map[string]string{
				"Location": "/test-file-id",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "Incorrect private file description (400)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					fileDesc := models.NewPrivateFileV1{
						Name: "test.jpg",
						Mime: "image/jpg",
					}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(fileDesc)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("POST", "/", &buf)
					r.Header.Set("Content-Type", "application/vnd.newPrivateFile.v1+json")
					return r
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusBadRequest,
		},
		{
			name: "Incorrect public file description (400)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					fileDesc := models.NewPublicFileV1{
						Mime: "image/jpg",
					}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(fileDesc)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("POST", "/", &buf)
					r.Header.Set("Content-Type", "application/vnd.newPublicFile.v1+json")
					return r
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusBadRequest,
		},
		{
			name: "Unsupported file description (415)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					fileDesc := models.NewPublicFileV1{
						Mime: "image/jpg",
					}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(fileDesc)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("POST", "/", &buf)
					r.Header.Set("Content-Type", "application/json")
					return r
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusUnsupportedMediaType,
		},
		{
			name: "Server-side issue (500)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					fileDesc := models.NewPrivateFileV1{
						Name:   "test.jpg",
						Mime:   "image/jpg",
						Access: []string{"john", "bob", "alice"},
					}
					var buf bytes.Buffer
					err := json.NewEncoder(&buf).Encode(fileDesc)
					if err != nil {
						log.Fatal(err)
					}
					r := httptest.NewRequest("POST", "/", &buf)
					r.Header.Set("Content-Type", "application/vnd.newPrivateFile.v1+json")
					r.Header.Set("request-id", "test-request-id")
					return r
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("AllocateNewFile", mock.Anything, mock.Anything).Return("", errors.New("test error"))
					return s
				}(),
			},
			wantHeaders: map[string]string{
				"issue": "test-request-id",
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				storage: tt.testService.storage,
			}
			s.allocateNewFile(tt.args.w, tt.args.r)

			for k, v := range tt.wantHeaders {
				require.Equal(t, v, func() string {
					h := tt.args.w.Result().Header
					if h == nil {
						return ""
					}
					v := h[k]
					if len(v) == 0 {
						return ""
					}
					return v[0]
				}())
			}

			require.Equal(t, tt.wantStatus, tt.args.w.Code)
		})
	}
}
