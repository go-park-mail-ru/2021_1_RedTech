package http

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/utils/session"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (handler *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)

	// Первый аргумент в парсинге беззнаковых чисел - база системы счисления, второй -
	// количество бит, которые он занимает. Четырех миллиардов пользователей нам хватит
	userId64, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		log.Printf("Error while getting user: %s", err)
		http.Error(w, `{"error":"bad id"}`, http.StatusBadRequest)
		return
	}
	userId := uint(userId64)

	isCurrent := false
	sess, err := getSession(r)
	if err == nil && session.Manager.Check(sess) == nil && sess.UserID == userId {
		isCurrent = true
	}

	userToSend, err := handler.UHandler.GetById(userId)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	if isCurrent {
		userToSend = userToSend.Private()
	} else {
		userToSend = userToSend.Public()
	}

	if err := json.NewEncoder(w).Encode(userToSend); err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
}

func (handler *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	sess, err := getSession(r)
	if err != nil || session.Manager.Check(sess) != nil {
		http.Error(w, `{"message":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// we use gorilla sessions, we can't do it another way
	userToSend := domain.User{
		ID: sess.UserID,
	}

	if err := json.NewEncoder(w).Encode(userToSend); err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
}
