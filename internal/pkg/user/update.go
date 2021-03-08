package user

import (
	"encoding/json"
	"log"
	"net/http"
)

type userUpdateForm struct {
	Email           string `json:"email"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func (form userUpdateForm) isEmpty() bool {
	return form.Email == "" && form.Username == "" && form.Password == "" && form.ConfirmPassword == ""
}

func (api *Handler) Update(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	userToUpdate := userUpdateForm{}
	if err := decoder.Decode(&userToUpdate); err != nil {
		log.Printf("Error while unmarshalling JSON")
		http.Error(w, `{"error": "bad form"}`, http.StatusBadRequest)
		return
	}
	if userToUpdate.isEmpty() {
		log.Printf("Empty form")
		http.Error(w, `{"error": "empty form"}`, http.StatusBadRequest)
		return
	}

	if err := updateUser(userToUpdate); err != nil {
		log.Printf("Error while updating user")
		http.Error(w, `{"error": "error while updating user"}`, http.StatusBadRequest)
		return
	}
	// TODO handle session

	if err := updateUser(userToUpdate); err != nil {
		log.Printf("Error while updating user")
		http.Error(w, `{"error": "error while updating user"}`, http.StatusBadRequest)
		return
	}
	// TODO handle session
}

func updateUser(userUpdateForm) (err error ) {
	// TODO handle data
	return
}
