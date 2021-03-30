package http

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/utils/fileutils"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
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
		http.Error(w, `{"error":"server"}`, http.StatusInternalServerError)
		return
	}

	sess, err := getSession(r)
	if err != nil {
		http.Error(w, `{"error":"can't find user'"}`, http.StatusBadRequest)
		return
	}
	if uint(urlID) != sess.UserID {
		log.Print("User try update wrong avatar")
		http.Error(w, `{"error":"wrong user"}`, http.StatusForbidden)
		return
	}

	r.ParseMultipartForm(5 * 1024 * 1024)
	uploaded, header, err := r.FormFile("avatar")
	if err != nil {
		log.Printf("error while file parsing file: %s", err)
		http.Error(w, `{"error":"server"}`, http.StatusForbidden)
	}
	defer uploaded.Close()

	filename, err := fileutils.UploadFile(uploaded, root, path, urlRoot, filepath.Ext(header.Filename))
	if err != nil {
		log.Printf("Upload error: %s", err)
		http.Error(w, `{"error":"server"}`, http.StatusForbidden)
	}

	err = handler.UHandler.Update(&domain.User{
		ID:     sess.UserID,
		Avatar: filename,
	})
	if err != nil {
		log.Printf("error while updating user: %s", err)
		http.Error(w, `{"error":"server"}`, http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, `{"user_avatar":"%s"}`, filename)
}
