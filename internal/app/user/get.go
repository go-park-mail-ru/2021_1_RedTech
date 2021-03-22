package user

import (
	"Redioteka/internal/pkg/utils/session"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func (api *Handler) Get(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	// Первый аргумент в парсинге беззнаковых чисел - база системы счисления, второй -
	// количество бит, которые он занимает. Четырех миллиардов пользователей нам хватит
	userId64, err := strconv.ParseUint(vars["id"], 10, 64)
	userId := uint(userId64)

	if err != nil {
		log.Printf("Error while getting user: %s", err)
		http.Error(w, `{"error":"bad id"}`, http.StatusBadRequest)
		return
	}

	isCurrent := false
	currentId, err := getCurrentId(r)
	if err == nil && currentId == userId {
		isCurrent = true
	}

	if err := sendByID(userId, !isCurrent, w); err != nil {
		log.Printf("Error while finding user: %s", err)
		http.Error(w, `{"error":"server can't send user'"}`, http.StatusBadRequest)
		return
	}
}

func (api *Handler) Me(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userId, err := getCurrentId(r)
	if err != nil {
		log.Printf("Error while getting id: %s", err)
		http.Error(w, `{"error":"can't find user"}`, http.StatusBadRequest)
		return
	}

	userToSend := &User{
		ID: userId,
	}

	if err := json.NewEncoder(w).Encode(userToSend); err != nil {
		fmt.Errorf("error while marshalling JSON: %s", err)
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}
}

func getCurrentId(r *http.Request) (uint, error) {
	userId, err := session.Check(r)
	if err != nil {
		log.Printf("Error while getting session: %s", err)
		return 0, errors.New("can't find user")
	}
	return userId, nil
}

func sendByID(userId uint, isPublic bool, w http.ResponseWriter) error {
	user := data.getByID(userId)

	if user == nil {
		log.Printf("Can't find user with id %d", userId)
		return errors.New("can't find user")
	}

	var userToSend *User
	if isPublic {
		userToSend = user.Public()
	} else {
		userToSend = user.Private()
	}

	if err := json.NewEncoder(w).Encode(userToSend); err != nil {
		return fmt.Errorf("error while marshalling JSON: %s", err)
	}
	return nil
}
