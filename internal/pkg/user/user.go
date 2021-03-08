package user

import (
	"crypto/sha256"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

const hashLen = 32

type User struct {
	ID       uint          `json:"id"`
	Email    string        `json:"email"`
	Username string        `json:"username"`
	Password [hashLen]byte `json:"-"`
}

type Handler struct {
}

type userSignupForm struct {
	Login         string `json:"username"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	PasswordCheck string `json:"repeated_password"`
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

type usersData struct {
	sync.Mutex
	users    map[string]*User
	sessions map[string][hashLen]byte
}

var data = usersData{
	users:    make(map[string]*User),
	sessions: make(map[string][hashLen]byte),
}

//Login - handler for user authorization
func (api *Handler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	userForm := new(userLoginForm)
	err := decoder.Decode(userForm)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		http.Error(w, `{"error":"bad form"}`, 400)
		return
	}
	if userForm.isEmpty() {
		log.Printf("Empty form of login")
		http.Error(w, `{"error":"Empty login or password"}`, 400)
		return
	}

	data.Lock()
	key := userForm.Email
	user, exists := data.users[key]
	if exists != true || user.Password != sha256.Sum256([]byte(userForm.Password)) {
		log.Printf("This user does not exist")
		http.Error(w, `{"error":"Wrong login or password"}`, 400)
		data.Unlock()
	}

	session := sha256.Sum256(append([]byte(key), byte(rand.Int())))
	data.sessions[key] = session
	data.Unlock()

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   string(session[:]),
		Expires: time.Now().Add(24 * time.Hour),
	}
	http.SetCookie(w, cookie)

	err = json.NewEncoder(w).Encode(data.users[key])
	if err != nil {
		log.Printf("error while marshalling JSON: %s", err)
		http.Error(w, `{"error":"server"}`, 500)
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
		http.Error(w, `{"error":"bad form"}`, 400)
		return
	}
	if userForm.isEmpty() {
		log.Printf("Empty form of signup")
		http.Error(w, `{"error":"Empty fields in form"}`, 400)
		return
	}

	if userForm.Password != userForm.PasswordCheck {
		log.Printf("Passwords do not match")
		http.Error(w, `{"error":"Passwords do not match"}`, 400)
		return
	}

	data.Lock()
	key := userForm.Email
	if _, exists := data.users[key]; exists == true {
		log.Printf("This user already exists")
		http.Error(w, `{"error":"Wrong username or password"}`, 400)
		data.Unlock()
		return
	}

	id := len(data.users) + 1
	session := sha256.Sum256(append([]byte(key), byte(rand.Int())))
	data.users[key] = &User{
		Username: userForm.Login,
		Password: sha256.Sum256([]byte(userForm.Password)),
		Email:    userForm.Email,
		ID:       uint(id),
	}
	data.sessions[key] = session
	data.Unlock()

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   string(session[:]),
		Expires: time.Now().Add(24 * time.Hour),
	}
	http.SetCookie(w, cookie)

	err = json.NewEncoder(w).Encode(data.users[key])
	if err != nil {
		log.Printf("error while marshalling JSON: %s", err)
		http.Error(w, `{"error":"server"}`, 500)
		return
	}
}

func (api *Handler) Logout(w http.ResponseWriter, r *http.Request) {
}

func (api *Handler) Get(w http.ResponseWriter, r *http.Request) {
}

func (api *Handler) Me(w http.ResponseWriter, r *http.Request) {
}

func (api *Handler) Avatar(w http.ResponseWriter, r *http.Request) {
}
