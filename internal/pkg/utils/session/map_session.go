package session

import (
	"Redioteka/internal/pkg/utils/randstring"
	"errors"
	"log"
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
	log.Printf("Session for user %d created", sess.UserID)

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
		log.Print("Bad cookie")
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

	log.Print("Successful delete session:", sess)
	sess.CookieExpiration = time.Now().AddDate(0, 0, -1)
	return nil
}
