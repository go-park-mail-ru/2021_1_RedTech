package session

import (
	"Redioteka/internal/pkg/utils/log"
	"fmt"
	"net/http"
	"time"

	"github.com/tarantool/go-tarantool"
)

const sessKeyLen = 32

type Session struct {
	UserID           uint
	Cookie           string
	CookieExpiration time.Time
}

type SessionManager interface {
	Create(*Session) error
	Check(*Session) error
	Delete(*Session) error
}

func getSessionManager() SessionManager {
	tarantoolAddress := "tarantool:5555"
	opts := tarantool.Opts{User: "redtech", Pass: "netflix"}
	conn, err := tarantool.Connect(tarantoolAddress, opts)
	if err != nil {
		log.Log.Warn(fmt.Sprintf("tarantool connection refused: %s - using map", err))
		return NewSessionMap()
	}
	log.Log.Info("Successful connect to tarantool")
	return NewSessionTarantool(conn)
}

func Destruct() {
	switch Manager.(type) {
	case *SessionTarantool:
		tarantoolManager, ok := Manager.(*SessionTarantool)
		if !ok {
			log.Log.Warn("Cannot cast to SessionTarantool")
		}
		err := tarantoolManager.tConn.Close()
		if err != nil {
			log.Log.Warn("Tarantool connection closing failed")
		}
		log.Log.Info("Tarantool connection closed")
	default:
		log.Log.Info("Nothing to be done")
	}
}

func SetSession(w http.ResponseWriter, s *Session) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    s.Cookie,
		Expires:  s.CookieExpiration,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})
}

func GetSession(r *http.Request) (*Session, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Error while getting session cookie: %s", err))
		return &Session{}, err
	}
	return &Session{Cookie: cookie.Value}, nil
}

var Manager = getSessionManager()
