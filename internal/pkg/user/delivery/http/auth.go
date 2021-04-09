package http

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/user"
	"Redioteka/internal/pkg/utils/jsonerrors"
	"Redioteka/internal/pkg/utils/session"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

//Signup - handler for user registration
func (handler *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	signupForm := &domain.User{}
	if err := decoder.Decode(signupForm); err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		http.Error(w, jsonerrors.JSONDecode, http.StatusBadRequest)
		return
	}

	createdUser, sess, err := handler.UUsecase.Signup(signupForm)
	if err != nil {
		log.Printf("Signup error")
		http.Error(w, jsonerrors.JSONMessage("signup"), user.CodeFromError(err))
		return
	}

	session.SetSession(w, sess)

	if err = json.NewEncoder(w).Encode(createdUser); err != nil {
		log.Printf("error while marshalling JSON: %s", err)
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
		log.Printf("error while unmarshalling JSON: %s", err)
		http.Error(w, jsonerrors.JSONDecode, http.StatusBadRequest)
		return
	}

	loggedUser, sess, err := handler.UUsecase.Login(loginForm)
	if err != nil {
		log.Printf("error while login: %s", err)
		http.Error(w, jsonerrors.JSONMessage("login"), user.CodeFromError(err))
		return
	}

	session.SetSession(w, sess)

	err = json.NewEncoder(w).Encode(loggedUser)
	if err != nil {
		log.Printf("error while marshalling JSON: %s", err)
		http.Error(w, jsonerrors.Session, http.StatusInternalServerError)
		return
	}
}

//Logout - handler for user logout with session deleting
func (handler *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	sess, err := session.GetSession(r)
	if err != nil {
		http.Error(w, jsonerrors.Session, http.StatusBadRequest)
		return
	}

	sess, err = handler.UUsecase.Logout(sess)
	if err != nil {
		log.Printf("error while logout user: %s", err)
		http.Error(w, jsonerrors.JSONMessage("session deletion"), http.StatusInternalServerError)
		return
	}

	session.SetSession(w, sess)
	fmt.Fprint(w, jsonerrors.JSONMessage("OK"))
}
