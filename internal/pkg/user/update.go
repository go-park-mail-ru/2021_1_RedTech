package user

import (
	"encoding/json"
	"net/http"
)

func (api *Handler) Update(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	userToUpdate := User{}
	if err := decoder.Decode(&userToUpdate); err != nil {
		_, _ = w.Write([]byte("{}"))
		return
	}
	updateUser(userToUpdate)
}

func updateUser(user User) {
}
