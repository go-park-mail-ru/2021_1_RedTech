package session

import (
	"Redioteka/internal/pkg/authorization/delivery/grpc/proto"
	"context"
	"time"
)

type GrpcSession struct {
	authService proto.AuthorizationClient
}

func NewGrpcSession(aClient proto.AuthorizationClient) SessionManager {
	return &GrpcSession{
		authService: aClient,
	}
}

func (sm *GrpcSession) Create(sess *Session) error {
	resSession, err := sm.authService.CreateSession(context.Background(), &proto.CreateSessionParams{
		UserId: uint32(sess.UserID),
	})
	if err != nil {
		return err
	}
	sess.Cookie = resSession.Cookie
	parsedTime, err := time.Parse(time.RFC3339, resSession.CookieExpiration)
	sess.CookieExpiration = parsedTime
	return nil
}

func (sm *GrpcSession) Check(sess *Session) error {
	return nil
}

func (sm *GrpcSession) Delete(sess *Session) error {
	return nil
}
