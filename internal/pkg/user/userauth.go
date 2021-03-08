package user

import (
	"Redioteka/internal/pkg/session"
	"crypto/sha256"
	"encoding/json"
	"log"
	"net/http"
)

type userSignupForm struct {
	Login         string `json:"username"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	PasswordCheck string `json:"confirm_password"`
}

func (form *userSignupForm) isEmpty() bool {
	return form.Login == "" || form.Email == "" || form.Password == "" || form.PasswordCheck == ""
}

type userLoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (form *userLoginForm) isEmpty() bool {
	return form.Email == "" || form.Password == ""
}

//Login - handler for user authorization
func (api *Handler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	userForm := new(userLoginForm)
	err := decoder.Decode(userForm)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		http.Error(w, `{"error":"bad form"}`, http.StatusBadRequest)
		return
	}
	if userForm.isEmpty() {
		log.Printf("Empty form of login")
		http.Error(w, `{"error":"Empty login or password"}`, http.StatusBadRequest)
		return
	}

	key := userForm.Email
	data.Lock()
	user, exists := data.users[key]
	data.Unlock()

	if exists != true || user.Password != sha256.Sum256([]byte(userForm.Password)) {
		log.Printf("This user does not exist")
		http.Error(w, `{"error":"Wrong login or password"}`, http.StatusBadRequest)
		return
	}

	err = session.Create(w, r, user.ID)
	if err != nil {
		log.Printf("error while creating session cookie: %s", err)
		http.Error(w, `{"error":"server"}`, http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(data.users[key])
	if err != nil {
		log.Printf("error while marshalling JSON: %s", err)
		http.Error(w, `{"error":"server"}`, http.StatusInternalServerError)
		return
	}
}

//Signup - handler for user registration
func (api *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	userForm := new(userSignupForm)
	err := decoder.Decode(userForm)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		http.Error(w, `{"error":"bad form"}`, http.StatusBadRequest)
		return
	}
	if userForm.isEmpty() {
		log.Printf("Empty form of signup")
		http.Error(w, `{"error":"Empty fields in form"}`, http.StatusBadRequest)
		return
	}

	if userForm.Password != userForm.PasswordCheck {
		log.Printf("Passwords do not match")
		http.Error(w, `{"error":"Passwords do not match"}`, http.StatusBadRequest)
		return
	}

	key := userForm.Email
	data.Lock()
	if _, exists := data.users[key]; exists == true {
		data.Unlock()
		log.Printf("This user already exists")
		http.Error(w, `{"error":"Wrong username or password"}`, http.StatusBadRequest)
		return
	}

	id := len(data.users) + 1
	data.users[key] = &User{
		Username: userForm.Login,
		Password: sha256.Sum256([]byte(userForm.Password)),
		Email:    userForm.Email,
		ID:       uint(id),
	}
	data.Unlock()

	err = session.Create(w, r, uint(id))
	if err != nil {
		log.Printf("error while creating session cookie: %s", err)
		http.Error(w, `{"error":"server"}`, http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(data.users[key])
	if err != nil {
		log.Printf("error while marshalling JSON: %s", err)
		http.Error(w, `{"error":"server"}`, http.StatusInternalServerError)
		return
	}
}

func (api *Handler) Logout(w http.ResponseWriter, r *http.Request) {
}
