package session

import (
	"Redioteka/internal/pkg/utils/randstring"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/tarantool/go-tarantool"
)

type SessionTarantool struct {
	tConn *tarantool.Connection
}

func NewSessionTarantool(conn *tarantool.Connection) *SessionTarantool {
	return &SessionTarantool{
		tConn: conn,
	}
}

func (sm *SessionTarantool) Create(sess *Session) error {
	cookieValue := randstring.RandString(sessKeyLen)
	expiration := time.Now().AddDate(0, 0, 1)

	_, err := sm.tConn.Insert("session", []interface{}{cookieValue, sess.UserID, expiration.Unix()})
	if err != nil {
		log.Printf("Error while creating session for user %d: %s", sess.UserID, err.Error())
		return err
	}

	sess.Cookie = cookieValue
	sess.CookieExpiration = expiration
	log.Printf("Session for user %d created", sess.UserID)
	return nil
}

func (sm *SessionTarantool) Check(sess *Session) error {
	resp, err := sm.tConn.Select("session", "primary", 0, 1, tarantool.IterEq, []interface{}{sess.Cookie})
	if err != nil {
		log.Print("Cannot check session:", err)
		return err
	}

	data := resp.Data[0]
	log.Print("data:", data)
	if data == nil {
		return errors.New("Getting no data from session store while session check")
	}
	sessionDataSlice, ok := data.([]interface{})
	if !ok {
		return fmt.Errorf("Cannot cast session data: %v", sessionDataSlice)
	}
	if sessionDataSlice[0] == nil || sessionDataSlice[1] == nil || sessionDataSlice[2] == nil {
		return errors.New("Getting no data from session store while session check")
	}

	cookie, ok := sessionDataSlice[0].(string)
	if !ok {
		return fmt.Errorf("Cannot cast session data: %v", sessionDataSlice[0])
	}
	expire, ok := sessionDataSlice[2].(uint64)
	if !ok {
		return fmt.Errorf("Cannot cast session data: %v", sessionDataSlice[2])
	}

	if cookie != sess.Cookie || time.Now().Sub(time.Unix(int64(expire), 0)) > 0 {
		log.Print("Bad cookie")
		return errors.New("Cookie value does not match or already expired")
	}

	sess.UserID = uint(sessionDataSlice[1].(uint64))
	sess.CookieExpiration = time.Unix(int64(expire), 0)
	return nil
}

func (sm *SessionTarantool) Delete(sess *Session) error {
	resp, err := sm.tConn.Delete("session", "primary", []interface{}{sess.Cookie})
	if err != nil {
		log.Print("Cannot delete session:", err)
		return err
	}

	log.Print("Successful delete session:", resp.Data)
	sess.CookieExpiration = time.Now().AddDate(0, 0, -1)
	return nil
}
