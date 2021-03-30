package session

import (
	"log"
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
	tarantoolAddress := "127.0.0.1:5555"
	opts := tarantool.Opts{User: "redtech", Pass: "netflix"}
	conn, err := tarantool.Connect(tarantoolAddress, opts)
	if err != nil {
		log.Printf("tarantool connection refused: %s - using map", err.Error())
		return NewSessionMap()
	}
	return NewSessionTarantool(conn)
}

var Manager = getSessionManager()
