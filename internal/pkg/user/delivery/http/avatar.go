package http

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/utils/fileutils"
	"Redioteka/internal/pkg/utils/session"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
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
		http.Error(w, `{"error":"server"}`, http.StatusInternalServerError)
		return
	}

	userID, err := session.Check(r)
	if userID == 0 || err != nil {
		log.Printf("Error while getting session: %s", err)
		http.Error(w, `{"error":"can't find user'"}`, http.StatusBadRequest)
		return
	}
	if uint(urlID) != userID {
		log.Print("User try update wrong avatar")
		http.Error(w, `{"error":"wrong user"}`, http.StatusForbidden)
		return
	}

	filename, err := fileutils.UploadFromRequest(r, root, path, urlRoot)
	if err != nil {
		log.Printf("Upload error: %s", err)
		http.Error(w, `{"error":"server"}`, http.StatusForbidden)
	}

	err = handler.UHandler.Update(&domain.User{
		ID:     userID,
		Avatar: filename,
	})
	if err != nil {
		log.Printf("error while updating user: %s", err)
		http.Error(w, `{"error":"server"}`, http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, `{"user_avatar":"%s"}`, filename)
}
