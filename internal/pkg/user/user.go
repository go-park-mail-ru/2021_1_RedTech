package user

import (
	"net/http"
	"sync"
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

type usersData struct {
	sync.Mutex
	users map[uint]*User
}

func (data *usersData) addUser(u *User) {
	data.Lock()
	data.users[u.ID] = u
	data.Unlock()
}

func (data *usersData) getByID(id uint) *User {
	data.Lock()
	user, exists := data.users[id]
	data.Unlock()
	if exists {
		return user
	}
	return nil
}

func (data *usersData) getByEmail(email string) *User {
	var result *User
	data.Lock()
	for _, user := range data.users {
		if user.Email == email {
			result = user
			break
		}
	}
	data.Unlock()
	return result
}

var data = usersData{
	users: make(map[uint]*User),
}

func (api *Handler) Avatar(w http.ResponseWriter, r *http.Request) {
}


