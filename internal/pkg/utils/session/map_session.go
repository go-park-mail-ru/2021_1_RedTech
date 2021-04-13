package session

import (
	"Redioteka/internal/pkg/utils/log"
	"Redioteka/internal/pkg/utils/randstring"
	"errors"
	"fmt"
	"sync"
	"time"
)

type SessionMap struct {
	sync.Mutex
	store map[string]*Session
}

func NewSessionMap() *SessionMap {
	return &SessionMap{
		store: map[string]*Session{},
	}
}

func (sm *SessionMap) Create(sess *Session) error {
	cookieValue := randstring.RandString(sessKeyLen)
	expiration := time.Now().AddDate(0, 0, 1)

	sess.Cookie = cookieValue
	sess.CookieExpiration = expiration
	log.Log.Info(fmt.Sprintf("Session for user %d created", sess.UserID))

	sm.Lock()
	sm.store[cookieValue] = sess
	sm.Unlock()
	return nil
}

func (sm *SessionMap) Check(sess *Session) error {
	sm.Lock()
	s, exist := sm.store[sess.Cookie]
	sm.Unlock()

	if !exist || time.Now().Sub(s.CookieExpiration) > 0 {
		log.Log.Warn("Bad cookie")
		return errors.New("Cookie value does not match or already expired")
	}

	sess.UserID = s.UserID
	sess.CookieExpiration = s.CookieExpiration
	return nil
}

func (sm *SessionMap) Delete(sess *Session) error {
	sm.Lock()
	delete(sm.store, sess.Cookie)
	sm.Unlock()

	log.Log.Info(fmt.Sprint("Successful delete session:", sess))
	sess.CookieExpiration = time.Now().AddDate(0, 0, -1)
	return nil
}
