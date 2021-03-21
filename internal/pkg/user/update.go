package user

import (
	"Redioteka/internal/pkg/session"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func (user User) isUpdateValid() bool {
	return !(user.Email == "" && user.Username == "")
}

func (user User) updateUser(userToUpdate *User) error {
	if !user.isUpdateValid() {
		log.Printf("Form validity error")
		return errors.New("invalid user update JSON")
	}

	if user.Email != userToUpdate.Email && user.Email != "" {
		userToUpdate.Email = user.Email
	}

	if user.Username != userToUpdate.Username && user.Username != "" {
		userToUpdate.Username = user.Username
	}
	return nil
}

func updateCurrentUser(r *http.Request, update *User) error {
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
	userUpdate := &User{}

	if err := decoder.Decode(userUpdate); err != nil {
		log.Printf("Error while unmarshalling JSON")
		http.Error(w, `{"error":"bad form"}`, http.StatusBadRequest)
		return
	}

	if err := updateCurrentUser(r, userUpdate); err != nil {
		log.Printf("Error while updating user")
		http.Error(w, `{"error":"error while updating user"}`, http.StatusBadRequest)
		return
	}

	userId, err := getCurrentId(r)
	if err != nil {
		log.Printf("Error while getting current user") http.Error(w, `{"error":"error while updating user"}`, http.StatusBadRequest)
	}

	if err := sendByID(userId, false, w); err != nil {
		log.Printf("Error while sending updated user")
		http.Error(w, `{"error":"error while sending user"}`, http.StatusBadRequest)
		return
	}
}
