package session

import (
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

const secondsInDay = 86400

var storeKey = securecookie.GenerateRandomKey(32)
var store = sessions.NewCookieStore(storeKey)

//Create - func for creating session-cookie of user
func Create(w http.ResponseWriter, r *http.Request, userID uint) error {
	session, err := store.Get(r, "session_id")
	if err != nil {
		return err
	}

	session.Options = &sessions.Options{
		MaxAge:   secondsInDay,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}
	session.Values["id"] = userID

	log.Printf("Setting user id %v", userID)

	err = session.Save(r, w)
	if err != nil {
		return err
	}
	return nil
}

//Delete - func for deleting session-cookie of user
func Delete(w http.ResponseWriter, r *http.Request, userID uint) error {
	session, err := store.Get(r, "session_id")
	if err != nil {
		return err
	}

	session.Options = &sessions.Options{
		MaxAge:   -secondsInDay,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}
	if id, exist := session.Values["id"]; exist == true && id == userID {
		err := session.Save(r, w)
		if err != nil {
			return err
		}
		session.Values["id"] = uint(0)
	}
	return nil
}

//Check - func for checking session-cookie
func Check(r *http.Request) (uint, error) {
	session, err := store.Get(r, "session_id")
	if err != nil {
		return 0, err
	}

	log.Printf("Checking user from session with values: %v", session.Values)
	user, exist := session.Values["id"]
	if exist != true {
		return 0, errors.New("User does not exist")
	}
	return user.(uint), nil
}
