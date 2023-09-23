package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// https://barpav.github.io/msg-api-spec/#/files/post_files__id_
func (s *Service) uploadNewFileContent(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/octet-stream" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	ctx := r.Context()
	id := chi.URLParam(r, "id")

	info, err := s.storage.AllocatedFileInfo(ctx, id)

	if err != nil {
		logAndReturnErrorWithIssue(w, r, err, "Failed to receive allocated file info")
		return
	}

	user := authenticatedUser(r)

	if info == nil || !info.HasAccess(user) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if user != info.Owner {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var fileSize int
	fileSize, err = s.storage.FileSize(ctx, id)

	if err != nil {
		logAndReturnErrorWithIssue(w, r, err, "Failed to receive file upload status")
		return
	}

	if fileSize != 0 { // already uploaded
		w.WriteHeader(http.StatusConflict)
		return
	}

	err = s.storage.UploadFileContent(id, r.Body)

	if err != nil {
		logAndReturnErrorWithIssue(w, r, err, "Failed to upload file content")
		return
	}
}
