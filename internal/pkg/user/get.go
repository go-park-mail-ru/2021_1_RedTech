package user

import (
	"encoding/json"
	"log"
	"net/http"
)

type userGetForm struct {
	Username string `json:"username"`
}

func (form userGetForm) isEmpty() bool {
	return form.Username == ""
}

func (api *Handler) Get(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	userForm := new(userGetForm)
	if err := decoder.Decode(userForm); err != nil {
		log.Printf("Error while unmarshalling JSON: %s", err)
		http.Error(w, `{"error":"bad form"}`, http.StatusBadRequest)
		return
	}
	if userForm.isEmpty() {
		log.Print("Empty form")
		http.Error(w, `{"error":"empty form"}`, http.StatusBadRequest)
		return
	}
	if err := getUser(userForm); err != nil {
		log.Printf("Error while getting user: %s", err)
		http.Error(w, `{"error":"bad form"}`, http.StatusBadRequest)
		return
	}
	// TODO session handling
}

func getUser(form *userGetForm) (err error) {
	// TODO data handling
	return
}

func (api *Handler) Me(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if err := getMe(); err != nil {
		log.Printf("Error while getting user: %s", err)
		http.Error(w, `{"error":"bad form"}`, http.StatusBadRequest)
		return
	}
	// TODO session handling
}


func getMe() (err error) {
	// TODO data handling
	return
}

