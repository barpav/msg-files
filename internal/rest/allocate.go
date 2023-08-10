package rest

import (
	"net/http"

	"github.com/barpav/msg-files/internal/rest/models"
)

// https://barpav.github.io/msg-api-spec/#/files/post_files
func (s *Service) allocateNewFile(w http.ResponseWriter, r *http.Request) {
	var err error
	newFile := models.AllocatedFile{}

	switch r.Header.Get("Content-Type") {
	case "application/vnd.newPrivateFile.v1+json":
		fileDesc := models.NewPrivateFileV1{}
		err = fileDesc.Deserialize(r.Body)

		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		newFile.Name, newFile.Mime, newFile.Access = fileDesc.Name, fileDesc.Mime, fileDesc.Access
	case "application/vnd.newPublicFile.v1+json":
		fileDesc := models.NewPublicFileV1{}
		err = fileDesc.Deserialize(r.Body)

		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		newFile.Name, newFile.Mime = fileDesc.Name, fileDesc.Mime
	default:
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	newFile.Owner = authenticatedUser(r)

	var id string
	id, err = s.storage.AllocateNewFile(r.Context(), &newFile)

	if err != nil {
		logAndReturnErrorWithIssue(w, r, err, "Failed to allocate new file")
		return
	}

	w.Header().Set("Location", "/"+id)
	w.WriteHeader(http.StatusCreated)
}
