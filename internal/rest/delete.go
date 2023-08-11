package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// https://barpav.github.io/msg-api-spec/#/files/delete_files__id_
func (s *Service) deleteFile(w http.ResponseWriter, r *http.Request) {
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

	err = s.storage.DeleteFile(ctx, id)

	if err != nil {
		logAndReturnErrorWithIssue(w, r, err, "Failed to delete file")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
