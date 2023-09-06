package rest

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/barpav/msg-files/internal/rest/mocks"
	"github.com/barpav/msg-files/internal/rest/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_getFile(t *testing.T) {
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
			name: "File successfully received (200)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "/test-id", nil)
					return r
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("AllocatedFileInfo", mock.Anything, mock.Anything).Return(
						&models.AllocatedFile{
							Name: "test.jpg",
							Mime: "image/jpeg",
						},
						nil)
					s.On("FileSize", mock.Anything, mock.Anything).Return(1024, nil)
					s.On("DownloadFile", mock.Anything, mock.Anything).Return(nil)
					return s
				}(),
			},
			wantHeaders: map[string]string{
				"Content-Type":        "image/jpeg",
				"Content-Disposition": `attachment; filename="test.jpg"`,
				"Content-Length":      "1024",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "File info successfully received (200)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("HEAD", "/test-id", nil)
					return r
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("AllocatedFileInfo", mock.Anything, mock.Anything).Return(
						&models.AllocatedFile{
							Name: "test.jpg",
							Mime: "image/jpeg",
						},
						nil)
					s.On("FileSize", mock.Anything, mock.Anything).Return(1024, nil)
					return s
				}(),
			},
			wantHeaders: map[string]string{
				"Content-Type":        "image/jpeg",
				"Content-Disposition": `attachment; filename="test.jpg"`,
				"Content-Length":      "1024",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "File not found (404)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "/test-id", nil)
					return r
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("AllocatedFileInfo", mock.Anything, mock.Anything).Return(nil, nil)
					return s
				}(),
			},
			wantHeaders: map[string]string{},
			wantStatus:  http.StatusNotFound,
		},
		{
			name: "Server-side issue (500)",
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "/test-id", nil)
					r.Header.Set("request-id", "test-request-id")
					return r
				}(),
			},
			testService: testService{
				storage: func() *mocks.Storage {
					s := mocks.NewStorage(t)
					s.On("AllocatedFileInfo", mock.Anything, mock.Anything).Return(
						&models.AllocatedFile{
							Name: "test.jpg",
							Mime: "image/jpeg",
						},
						nil)
					s.On("FileSize", mock.Anything, mock.Anything).Return(1024, nil)
					s.On("DownloadFile", mock.Anything, mock.Anything).Return(errors.New("test error"))
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
			s.getFile(tt.args.w, tt.args.r)

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
