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
	login         string `json:"username"`
	email         string `json:"email"`
	password      string `json:"password"`
	passwordCheck string `json:"repeated_password"`
}

type userLoginForm struct {
	email    string `json:"email"`
	password string `json:"password"`
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

	data.Lock()
	key := userForm.email
	user, exists := data.users[key]
	if exists != true {
		log.Printf("This user does not exist")
		http.Error(w, `{"error":"Wrong username or password"}`, 400)
		return
	}
	if user.Password != sha256.Sum256([]byte(userForm.password)) {
		log.Printf("The user password does not match")
		http.Error(w, `{"error":"Wrong username or password"}`, 400)
		return
	}

	session := sha256.Sum256(append([]byte(key), byte(rand.Int())))
	data.sessions[key] = session
	data.Unlock()

	err = json.NewEncoder(w).Encode(data.users[key])
	if err != nil {
		log.Printf("error while marshalling JSON: %s", err)
		http.Error(w, `{"error":"server"}`, 500)
		return
	}
	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   string(session[:]),
		Expires: time.Now().Add(24 * time.Hour),
	}
	http.SetCookie(w, cookie)
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

	if userForm.password != userForm.passwordCheck {
		log.Printf("Passwords do not match")
		http.Error(w, `{"error":"Passwords do not match"}`, 400)
		return
	}

	data.Lock()
	key := userForm.email
	if _, exists := data.users[key]; exists == true {
		log.Printf("This user already exists")
		http.Error(w, `{"error":"Wrong username or password"}`, 400)
		return
	}

	id := len(data.users) + 1
	session := sha256.Sum256(append([]byte(key), byte(rand.Int())))
	data.users[key] = &User{
		Username: userForm.login,
		Password: sha256.Sum256([]byte(userForm.password)),
		Email:    userForm.email,
		ID:       uint(id),
	}
	data.sessions[key] = session
	data.Unlock()

	err = json.NewEncoder(w).Encode(data.users[key])
	if err != nil {
		log.Printf("error while marshalling JSON: %s", err)
		http.Error(w, `{"error":"server"}`, 500)
		return
	}
	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   string(session[:]),
		Expires: time.Now().Add(24 * time.Hour),
	}
	http.SetCookie(w, cookie)
}

func (api *Handler) Logout(w http.ResponseWriter, r *http.Request) {
}

func (api *Handler) Get(w http.ResponseWriter, r *http.Request) {
}

func (api *Handler) Me(w http.ResponseWriter, r *http.Request) {
}

func (api *Handler) Avatar(w http.ResponseWriter, r *http.Request) {
}
