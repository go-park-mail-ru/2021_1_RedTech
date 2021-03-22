package http

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/utils/fileutils"
	"Redioteka/internal/pkg/utils/randstring"
	"Redioteka/internal/pkg/utils/session"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

const (
	root    = "./img"
	urlRoot = "/static"
	path    = "/users/"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

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

	r.ParseMultipartForm(5 * 1024 * 1024)
	uploaded, header, err := r.FormFile("avatar")
	if err != nil {
		log.Printf("Error while uploading file: %s", err)
		http.Error(w, `{"error":"server"}`, http.StatusInternalServerError)
		return
	}
	defer uploaded.Close()

	filename := randstring.RandString(32) + filepath.Ext(header.Filename)
	log.Print("avatar name ", filename)
	file, err := fileutils.CreateFile(root, path, filename)
	if err != nil {
		log.Printf("error while creating file: %s", err)
		http.Error(w, `{"error":"server"}`, http.StatusInternalServerError)
		return
	}
	defer file.Close()

	filename = urlRoot + path + filename
	log.Print("avatar name ", filename)
	_, err = io.Copy(file, uploaded)
	if err != nil {
		log.Printf("error while writing in file: %s", err)
		http.Error(w, `{"error":"server"}`, http.StatusInternalServerError)
		return
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
