package http

import (
	"Redioteka/internal/pkg/domain"
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
		http.Error(w, `{"error":"bad form"}`, http.StatusBadRequest)
		return
	}

	user, err := handler.UHandler.Signup(signupForm)
	if err != nil {
		log.Printf("Signup error")
		http.Error(w, `{"error":"signup"}`, http.StatusInternalServerError)
		return
	}

	err = session.Create(w, r, user.ID)
	if err != nil {
		log.Printf("error while creating session cookie: %s", err)
		http.Error(w, `{"error":"server"}`, http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("error while marshalling JSON: %s", err)
		http.Error(w, `{"error":"server"}`, http.StatusInternalServerError)
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
		http.Error(w, `{"error":"bad form"}`, http.StatusBadRequest)
		return
	}

	user, err := handler.UHandler.Login(loginForm)
	if err != nil {
		log.Printf("error while login: %s", err)
		http.Error(w, `{"error":"login"}`, http.StatusInternalServerError)
		return
	}

	err = session.Create(w, r, user.ID)
	if err != nil {
		log.Printf("error while creating session cookie: %s", err)
		http.Error(w, `{"error":"server"}`, http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Printf("error while marshalling JSON: %s", err)
		http.Error(w, `{"error":"server"}`, http.StatusInternalServerError)
		return
	}
}

//Logout - handler for user logout with session deleting
func (handler *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userID, err := session.Check(r)
	if err != nil {
		log.Printf("error while logout user: %s", err)
		http.Error(w, `{"error":"user not found"}`, http.StatusBadRequest)
		return
	}

	err = session.Delete(w, r, userID)
	if err != nil {
		log.Printf("error while deleting session cookie: %s", err)
		http.Error(w, `{"error":"server"}`, http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, `{"status":"OK"}`)
}
