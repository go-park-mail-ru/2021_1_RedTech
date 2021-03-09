package user

import (
	"Redioteka/internal/pkg/session"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type userUpdateForm struct {
	Email              string `json:"email"`
	Username           string `json:"username"`
	NewPassword        string `json:"new_password"`
	ConfirmNewPassword string `json:"confirm_new_password"`
	OldPassword        string `json:"password"`
}

func (form userUpdateForm) hasUpdates() bool {
	return !(form.Email == "" && form.Username == "" && form.NewPassword == "" && form.ConfirmNewPassword == "")
}

func (form userUpdateForm) isValid() bool {
	// todo проверка уникальности мыла и ника
	return true
}

func (form userUpdateForm) updateUser(u *User) error {
	if !form.isValid() {
		log.Printf("Form validity error")
		return errors.New("invalid user update JSON")
	}

	if form.Email != u.Email && form.Email != "" {
		u.Email = form.Email
	}

	if form.Username != u.Username && form.Username != "" {
		u.Username = form.Email
	}

	if !passwordValid(u, form.OldPassword) {
		log.Printf("Error while updating user: passowrd doesn't pass")
		return errors.New("wrong password")
	}
	return nil
}

func passwordValid(u *User, password string) bool {
	return u.Password == sha256.Sum256([]byte(password))
}

func updateCurrentUser(r *http.Request, form userUpdateForm) error {
	user, err := getCurrentUser(r)
	if err != nil {
		log.Printf("Error while getting user")
		return err
	}

	if err := form.updateUser(user); err != nil {
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
	userUpdates := userUpdateForm{}

	if err := decoder.Decode(&userUpdates); err != nil {
		log.Printf("Error while unmarshalling JSON")
		http.Error(w, `{"error": "bad form"}`, http.StatusBadRequest)
		return
	}

	if !userUpdates.hasUpdates() {
		log.Printf("Empty form")
		http.Error(w, `{"error": "empty form"}`, http.StatusBadRequest)
		return
	}

	if err := updateCurrentUser(r, userUpdates); err != nil {
		log.Printf("Error while updating user")
		http.Error(w, `{"error": "error while updating user"}`, http.StatusBadRequest)
		return
	}
}
