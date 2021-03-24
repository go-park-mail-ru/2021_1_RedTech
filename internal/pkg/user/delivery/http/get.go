package http

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/user"
	"Redioteka/internal/pkg/utils/jsonerrors"
	"Redioteka/internal/pkg/utils/session"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func getCurrentId(r *http.Request) (uint, error) {
	userId, err := session.Check(r)
	if err != nil {
		log.Printf("Error while getting session: %s", err)
		return 0, user.UnauthorizedError
	}
	return userId, nil
}

func (handler *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)

	// Первый аргумент в парсинге беззнаковых чисел - база системы счисления, второй -
	// количество бит, которые он занимает. Четырех миллиардов пользователей нам хватит
	userId64, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		log.Printf("Error while getting user: %s", err)
		http.Error(w, jsonerrors.JSONMessage("var parse"), user.CodeFromError(err))
		return
	}
	userId := uint(userId64)

	isCurrent := false
	currentId, err := getCurrentId(r)
	if err == nil && currentId == userId {
		isCurrent = true
	}

	userToSend, err := handler.UUsecase.GetById(userId)
	if err != nil {
		http.Error(w, jsonerrors.JSONMessage("get"), user.CodeFromError(err))
		return
	}

	if isCurrent {
		userToSend = userToSend.Private()
	} else {
		userToSend = userToSend.Public()
	}

	if err := json.NewEncoder(w).Encode(userToSend); err != nil {
		http.Error(w, jsonerrors.JSONMessage("json encode"), http.StatusInternalServerError)
		return
	}
}

func (handler *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userId, err := getCurrentId(r)
	if err != nil {
		http.Error(w, jsonerrors.JSONMessage("unauthorized"), user.CodeFromError(err))
		return
	}

	// we use gorilla sessions, we can't do it another way
	userToSend := domain.User{
		ID: userId,
	}

	if err := json.NewEncoder(w).Encode(userToSend); err != nil {
		http.Error(w, jsonerrors.JSONMessage("json encode"), http.StatusInternalServerError)
		return
	}
}
