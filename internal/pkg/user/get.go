package user

import (
	"Redioteka/internal/pkg/session"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type userGet struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (api *Handler) Get(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	// Первый аргумент в парсинге беззнаковых чисел - база системы счисления, второй -
	// количество бит, которые он занимает. Четырех миллиардов пользователей нам хватит
	userId, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		log.Printf("Error while getting user: %s", err)
		http.Error(w, `{"error":"bad id"}`, http.StatusBadRequest)
		return
	}

	if err := sendUser(uint(userId), w); err != nil {
		log.Printf("Error while finding user: %s", err)
		http.Error(w, `{"error":"server can't send user'"}`, http.StatusBadRequest)
		return
	}
}

func (api *Handler) Me(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userId, err := session.Check(r)
	if err != nil {
		log.Printf("Error while getting session: %s", err)
		http.Error(w, `{"error":"can't find user'"}`, http.StatusBadRequest)
		return
	}

	if err := sendUser(userId, w); err != nil {
		log.Printf("Error while finding user: %s", err)
		http.Error(w, `{"error":"server can't send user'"}`, http.StatusBadRequest)
		return
	}
}

func sendUser(userId uint, w http.ResponseWriter) error {
	user := data.getByID(userId)

	if user == nil {
		log.Printf("Can't find user with id %d", userId)
		return errors.New("can't find user")
	}

	userToSend := &userGet{
		Username: user.Username,
		Email:    user.Email,
	}
	if err := json.NewEncoder(w).Encode(userToSend); err != nil {
		return fmt.Errorf("error while marshalling JSON: %s", err)
	}
	return nil
}
