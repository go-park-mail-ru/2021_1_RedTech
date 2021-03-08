package user

import (
	"log"
	"net/http"
)

func (api *Handler) Get(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	userId := 10
	// todo get Id
	if err := getUser(userId); err != nil {
		log.Printf("Error while getting user: %s", err)
		http.Error(w, `{"error":"bad form"}`, http.StatusBadRequest)
		return
	}
	// TODO session handling
}

func getUser(userId int) (err error) {
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
