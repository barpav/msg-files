package rest

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// https://barpav.github.io/msg-api-spec/#/files/get_files__id_
// https://barpav.github.io/msg-api-spec/#/files/head_files__id_
func (s *Service) getFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	info, err := s.storage.AllocatedFileInfo(ctx, id)

	if err != nil {
		logAndReturnErrorWithIssue(w, r, err, "Failed to receive allocated file info")
		return
	}

	if info == nil || !info.HasAccess(authenticatedUser(r)) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var fileSize int
	fileSize, err = s.storage.FileSize(ctx, id)

	if err != nil {
		logAndReturnErrorWithIssue(w, r, err, "Failed to receive file size")
		return
	}

	w.Header().Set("Content-Type", info.Mime)
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, info.Name))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileSize))

	if r.Method == http.MethodGet {
		err = s.storage.DownloadFile(id, w)

		if err != nil {
			logAndReturnErrorWithIssue(w, r, err, "Failed to download file data")
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
