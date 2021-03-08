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

func (api *Handler) Update(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	userToUpdate := userUpdateForm{}
	if err := decoder.Decode(&userToUpdate); err != nil {
		log.Printf("error while unmarshalling JSON")
		http.Error(w, `{"error": "bad form"}`, http.StatusBadRequest)
		return
	}
	updateUser(userToUpdate)
	// TODO handle session
}

func updateUser(userUpdateForm) {
	// TODO handle data
}
