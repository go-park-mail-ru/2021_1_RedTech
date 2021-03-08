package session

import (
	"errors"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

const secondsInDay = 86400

var storeKey = securecookie.GenerateRandomKey(10)
var store = sessions.NewCookieStore(storeKey)

//Create - func for creating session-cookie of user
func Create(w http.ResponseWriter, r *http.Request, user string) error {
	session, err := store.Get(r, "session_id")
	if err != nil {
		return err
	}

	store.MaxAge(secondsInDay)
	key := string(securecookie.GenerateRandomKey(32))
	session.Values[key] = user
	err = session.Save(r, w)
	if err != nil {
		return err
	}
	return nil
}

//Delete - func for deleting session-cookie of user
func Delete(w http.ResponseWriter, r *http.Request, cookie string) error {
	session, err := store.Get(r, "session_id")
	if err != nil {
		return err
	}

	store.MaxAge(-secondsInDay)
	if _, exist := session.Values[cookie]; exist == true {
		err := session.Save(r, w)
		if err != nil {
			return err
		}
	}
	return nil
}

//Check - func for checking session-cookie
func Check(r *http.Request, cookie string) (string, error) {
	session, err := store.Get(r, "session_id")
	if err != nil {
		return "", err
	}

	user, exist := session.Values[cookie]
	if exist != true {
		return "", errors.New("User does not exist")
	}
	return user.(string), nil
}
