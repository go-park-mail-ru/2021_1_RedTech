package http

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/utils/session"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type userSignupForm struct {
	Login         string `json:"username"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	PasswordCheck string `json:"confirm_password"`
}

func (form *userSignupForm) isValid() bool {
	return form.Login != "" && form.Email != "" && form.Password != "" && form.PasswordCheck == form.Password
}

type userLoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (form *userLoginForm) isValid() bool {
	return form.Email != "" && form.Password != ""
}

//Signup - handler for user registration
func (handler *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	userForm := new(userSignupForm)
	if err := decoder.Decode(userForm); err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		http.Error(w, `{"error":"bad form"}`, http.StatusBadRequest)
		return
	}
	if !userForm.isValid() {
		log.Printf("Invalid signup form")
		http.Error(w, `{"error":"Empty fields in form"}`, http.StatusBadRequest)
		return
	}

	id, err := handler.UHandler.Signup(&domain.User{
		Username: userForm.Login,
		Email:    userForm.Email,
		Password: sha256.Sum256([]byte(userForm.Password)),
	})

	if err != nil {
		log.Printf("Registration error")
		http.Error(w, `{"error":"register"}`, http.StatusInternalServerError)
		return
	}

	err = session.Create(w, r, id)
	if err != nil {
		log.Printf("error while creating session cookie: %s", err)
		http.Error(w, `{"error":"server"}`, http.StatusInternalServerError)
		return
	}

	newUser, err := handler.UHandler.GetById(id)
	if err != nil {
		log.Printf("Registration error")
		http.Error(w, `{"error":"register"}`, http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(newUser); err != nil {
		log.Printf("error while marshalling JSON: %s", err)
		http.Error(w, `{"error":"server"}`, http.StatusInternalServerError)
		return
	}
}

//Login - handler for user authorization
func (handler *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	userForm := new(userLoginForm)
	err := decoder.Decode(userForm)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		http.Error(w, `{"error":"bad form"}`, http.StatusBadRequest)
		return
	}

	if !userForm.isValid() {
		log.Printf("Invalid login form")
		http.Error(w, `{"error":"Invalid login form"}`, http.StatusBadRequest)
		return
	}

	id, err := handler.UHandler.Login(&domain.User{
		Email:    userForm.Email,
		Password: sha256.Sum256([]byte(userForm.Password)),
	})

	err = session.Create(w, r, id)
	if err != nil {
		log.Printf("error while creating session cookie: %s", err)
		http.Error(w, `{"error":"server"}`, http.StatusInternalServerError)
		return
	}

	user, err := handler.UHandler.GetById(id)
	if err != nil {
		log.Printf("Login error")
		http.Error(w, `{"error":"login"}`, http.StatusInternalServerError)
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
