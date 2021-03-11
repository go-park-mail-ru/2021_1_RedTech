package user

import (
	"Redioteka/internal/pkg/session"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type userUpdate struct {
	Email              string `json:"email,omitempty"`
	Username           string `json:"username,omitempty"`
	NewPassword        string `json:"new_password,omitempty"`
	ConfirmNewPassword string `json:"confirm_new_password,omitempty"`
	// OldPassword        string `json:"password"`
}

func (update userUpdate) isValid() bool {
	return !(update.Email == "" && update.Username == "" && update.NewPassword == "" && update.ConfirmNewPassword == "")
}

func (update userUpdate) updateUser(u *User) error {
	if !update.isValid() {
		log.Printf("Form validity error")
		return errors.New("invalid user update JSON")
	}

	if update.Email != u.Email && update.Email != "" {
		u.Email = update.Email
	}

	if update.Username != u.Username && update.Username != "" {
		u.Username = update.Username
	}

	/*
		if !passwordValid(u, update.OldPassword) {
			log.Printf("Error while updating user: passowrd doesn't pass")
			return errors.New("wrong password")
		}
	*/
	return nil
}

func passwordValid(u *User, password string) bool {
	return u.Password == sha256.Sum256([]byte(password))
}

func updateCurrentUser(r *http.Request, update userUpdate) error {
	user, err := getCurrentUser(r)
	if err != nil {
		log.Printf("Error while getting user")
		return err
	}

	if err := update.updateUser(user); err != nil {
		log.Printf("Error while updating user")
		return err
	}
	return nil
}

func getCurrentUser(r *http.Request) (user *User, err error) {
	currentUserId, err := session.Check(r)
	if err != nil {
		log.Printf("Can't find session")
		return
	}
	user = data.getByID(currentUserId)
	if user == nil {
		log.Printf("Can't find user with id %d", currentUserId)
		err = errors.New("error while accessing data")
		return
	}
	return
}

func (api *Handler) Update(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	userUpdates := userUpdate{}

	if err := decoder.Decode(&userUpdates); err != nil {
		log.Printf("Error while unmarshalling JSON")
		http.Error(w, `{"error":"bad form"}`, http.StatusBadRequest)
		return
	}

	if !userUpdates.isValid() {
		log.Printf("Invalid form")
		http.Error(w, `{"error":"Invalid form"}`, http.StatusBadRequest)
		return
	}

	if err := updateCurrentUser(r, userUpdates); err != nil {
		log.Printf("Error while updating user")
		http.Error(w, `{"error":"error while updating user"}`, http.StatusBadRequest)
		return
	}

	if err := sendCurrentUser(w, r); err != nil {
		log.Printf("Error while sending updated user")
		http.Error(w, `{"error":"error while sending user"}`, http.StatusBadRequest)
		return
	}

}
