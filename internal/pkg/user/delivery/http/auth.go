package http

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/user"
	"Redioteka/internal/pkg/utils/jsonerrors"
	"Redioteka/internal/pkg/utils/log"
	"Redioteka/internal/pkg/utils/session"
	"Redioteka/internal/pkg/utils/xss"
	"encoding/json"
	"fmt"
	"net/http"
)

//Signup - handler for user registration
func (handler *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	signupForm := &domain.User{}
	if err := decoder.Decode(signupForm); err != nil {
		log.Log.Warn(fmt.Sprintf("Cannot parse signup form: %s", err))
		http.Error(w, jsonerrors.JSONDecode, http.StatusBadRequest)
		return
	}
	xss.SanitizeUser(signupForm)

	createdUser, err := handler.UUsecase.Signup(signupForm)
	if err != nil {
		http.Error(w, jsonerrors.JSONMessage("signup"), user.CodeFromError(err))
		return
	}

	var sess *session.Session
	err := handler.SessionManager.Create(sess)
	session.SetSession(w, sess)

	if err = json.NewEncoder(w).Encode(createdUser); err != nil {
		log.Log.Error(err)
		http.Error(w, jsonerrors.JSONEncode, http.StatusInternalServerError)
		return
	}
}

//Login - handler for user authorization
func (handler *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	loginForm := &domain.User{}
	err := decoder.Decode(loginForm)
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Cannot parse login form: %s", err))
		http.Error(w, jsonerrors.JSONDecode, http.StatusBadRequest)
		return
	}

	loggedUser, err := handler.UUsecase.Login(loginForm)
	if err != nil {
		http.Error(w, jsonerrors.JSONMessage("login"), user.CodeFromError(err))
		return
	}

	var sess *session.Session
	err = handler.SessionManager.Create(sess)
	if err != nil {
		log.Log.Error(err)
		http.Error(w, jsonerrors.Session, http.StatusInternalServerError)
		return
	}
	session.SetSession(w, sess)

	err = json.NewEncoder(w).Encode(loggedUser)
	if err != nil {
		log.Log.Error(err)
		http.Error(w, jsonerrors.Session, http.StatusInternalServerError)
		return
	}
}

//Logout - handler for user logout with session deleting
func (handler *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	sess, err := session.GetSession(r)
	if err != nil {
		log.Log.Warn("Cannot logout - unauthorized")
		http.Error(w, jsonerrors.Session, http.StatusBadRequest)
		return
	}

	err = handler.UUsecase.Logout(sess)
	if err != nil {
		http.Error(w, jsonerrors.JSONMessage("session deletion"), http.StatusInternalServerError)
		return
	}

	session.SetSession(w, sess)
	w.WriteHeader(http.StatusOK)
}
