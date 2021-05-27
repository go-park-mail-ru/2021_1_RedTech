package http

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/user"
	"Redioteka/internal/pkg/utils/jsonerrors"
	"Redioteka/internal/pkg/utils/log"
	"Redioteka/internal/pkg/utils/session"
	"encoding/json"
	"fmt"
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
		log.Log.Warn(fmt.Sprintf("Cannot get user id: %s", err))
		http.Error(w, jsonerrors.URLParams, user.CodeFromError(err))
		return
	}
	userId := uint(userId64)

	isCurrent := false
	sess, err := session.GetSession(r)
	if err == nil && handler.SessionManager.Check(sess) == nil && sess.UserID == userId {
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
		log.Log.Error(err)
		http.Error(w, jsonerrors.JSONEncode, http.StatusInternalServerError)
		return
	}
}

func (handler *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	sess, err := session.GetSession(r)
	if err != nil || handler.SessionManager.Check(sess) != nil {
		log.Log.Info("User is unauthorized")
		http.Error(w, jsonerrors.JSONMessage("unauthorized"), user.CodeFromError(user.UnauthorizedError))
		return
	}

	userToSend := domain.User{
		ID: sess.UserID,
	}

	if err := json.NewEncoder(w).Encode(userToSend); err != nil {
		log.Log.Error(err)
		http.Error(w, jsonerrors.JSONEncode, http.StatusInternalServerError)
		return
	}
}
