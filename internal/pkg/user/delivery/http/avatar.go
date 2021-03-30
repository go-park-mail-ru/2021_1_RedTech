package http

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/user"
	"Redioteka/internal/pkg/utils/fileutils"
	"Redioteka/internal/pkg/utils/jsonerrors"
	"Redioteka/internal/pkg/utils/session"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

const (
	root    = "./img"
	urlRoot = "/static"
	path    = "/users/"
)

//Avatar - handler for uploading user avatar
func (handler *UserHandler) Avatar(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	idString := vars["id"]
	urlID, err := strconv.Atoi(idString)
	if err != nil {
		log.Print("Id is not a number")
		http.Error(w, jsonerrors.URLParams, http.StatusBadRequest)
		return
	}

	userID, err := session.Check(r)
	if userID == 0 || err != nil {
		log.Printf("Error while getting session: %s", err)
		http.Error(w, jsonerrors.Session, http.StatusUnauthorized)
		return
	}
	if uint(urlID) != userID {
		log.Print("User try update wrong avatar")
		http.Error(w, jsonerrors.JSONMessage("wrong id"), http.StatusForbidden)
		return
	}

	r.ParseMultipartForm(5 * 1024 * 1024)
	uploaded, header, err := r.FormFile("avatar")
	if err != nil {
		log.Printf("error while file parsing file: %s", err)
		http.Error(w, jsonerrors.JSONMessage("file parsing"), http.StatusBadRequest)
	}
	defer uploaded.Close()

	filename, err := fileutils.UploadFile(uploaded, root, path, urlRoot, filepath.Ext(header.Filename))
	if err != nil {
		log.Printf("Upload error: %s", err)
		http.Error(w, jsonerrors.JSONMessage("upload"), http.StatusInternalServerError)
	}

	err = handler.UUsecase.Update(&domain.User{
		ID:     userID,
		Avatar: filename,
	})
	if err != nil {
		log.Printf("error while updating user: %s", err)
		http.Error(w, jsonerrors.JSONMessage("update fail"), user.CodeFromError(err))
		return
	}
	if err := json.NewEncoder(w).Encode(map[string]string{"user_avatar": filename}); err != nil {
		log.Printf("Path encodind error: %s", err)
		http.Error(w, jsonerrors.JSONEncode, http.StatusInternalServerError)
		return

	}
}
