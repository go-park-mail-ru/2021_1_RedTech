package user

import (
	"Redioteka/internal/pkg/session"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type avatar struct {
	Data []byte `json:"avatar"`
}

const path = ".img/users/"

func (api *Handler) Avatar(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userID, err := session.Check(r)
	if userID == 0 || err != nil {
		log.Printf("Error while getting session: %s", err)
		http.Error(w, `{"error":"can't find user'"}`, http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	uploaded := new(avatar)
	err = decoder.Decode(uploaded)
	if err != nil {
		log.Printf("error while unmarshalling JSON: %s", err)
		http.Error(w, `{"error":"bad image"}`, http.StatusBadRequest)
		return
	}

	user := data.getByID(userID)
	hash := sha256.Sum256(append(uploaded.Data, []byte(user.Email)...))
	filename := path + string(hash[:])
	file, fileErr := os.Create(filename)
	if fileErr != nil {
		log.Printf("error while creating file: %s", fileErr)
		http.Error(w, `{"error":"server"}`, http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, fileErr = file.Write(uploaded.Data)
	if fileErr != nil {
		log.Printf("error while writing in file: %s", fileErr)
		http.Error(w, `{"error":"server"}`, http.StatusInternalServerError)
		return
	}

	user.Avatar = filename
	fmt.Fprintf(w, `{"user_avatar":"%s"}`, filename)
}
