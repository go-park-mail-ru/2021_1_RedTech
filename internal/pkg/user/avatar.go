package user

import (
	"Redioteka/internal/pkg/session"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

const (
	root    = "./img"
	urlRoot = "/static"
	path    = "/users/"
)

func createFile(dir, name string) (*os.File, error) {
	_, err := os.ReadDir(root + dir)
	if err != nil {
		err = os.MkdirAll(root+dir, 0777)
		if err != nil {
			return nil, err
		}
	}
	file, err := os.Create(root + dir + name)
	return file, err
}

//Avatar - handler for uploading user avatar
func (api *Handler) Avatar(w http.ResponseWriter, r *http.Request) {
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

	r.ParseMultipartForm(5 * 1024 * 1025)
	uploaded, handler, err := r.FormFile("avatar")
	if err != nil {
		log.Printf("Error while uploading file: %s", err)
		http.Error(w, `{"error":"server"}`, http.StatusInternalServerError)
		return
	}
	defer uploaded.Close()

	filename := string(securecookie.GenerateRandomKey(32)) + filepath.Ext(handler.Filename)
	file, err := createFile(path, filename)
	if err != nil {
		log.Printf("error while creating file: %s", err)
		http.Error(w, `{"error":"server"}`, http.StatusInternalServerError)
		return
	}
	defer file.Close()
	filename = urlRoot + path + filename

	_, err = io.Copy(file, uploaded)
	if err != nil {
		log.Printf("error while writing in file: %s", err)
		http.Error(w, `{"error":"server"}`, http.StatusInternalServerError)
		return
	}

	user := data.getByID(userID)
	user.Avatar = filename
	fmt.Fprintf(w, `{"user_avatar":"%s"}`, filename)
}