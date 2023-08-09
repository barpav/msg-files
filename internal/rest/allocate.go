package rest

import (
	"net/http"

	"github.com/barpav/msg-files/internal/rest/models"
)

// https://barpav.github.io/msg-api-spec/#/files/post_files
func (s *Service) allocateNewFile(w http.ResponseWriter, r *http.Request) {
	var (
		name, mime string
		access     []string
		err        error
	)

	switch r.Header.Get("Content-Type") {
	case "application/vnd.newPrivateFile.v1+json":
		fileDesc := models.NewPrivateFileV1{}
		err = fileDesc.Deserialize(r.Body)

		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		name, mime, access = fileDesc.Name, fileDesc.Mime, fileDesc.Access
	case "application/vnd.newPublicFile.v1+json":
		fileDesc := models.NewPublicFileV1{}
		err = fileDesc.Deserialize(r.Body)

		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		name, mime = fileDesc.Name, fileDesc.Mime
	default:
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	var id string
	id, err = s.storage.AllocateNewFile(r.Context(), name, mime, access)

	if err != nil {
		logAndReturnErrorWithIssue(w, r, err, "Failed to allocate new file")
		return
	}

	w.Header().Set("Location", "/"+id)
	w.WriteHeader(http.StatusCreated)
}
