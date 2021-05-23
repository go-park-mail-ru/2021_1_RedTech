package session

import (
	"Redioteka/internal/pkg/authorization/delivery/grpc/proto"
	"context"
	"google.golang.org/grpc"
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

func (sm *GrpcSession) handleSession(sess *Session, handle func(ctx context.Context, in *proto.Session, opts ...grpc.CallOption) (*proto.Session, error)) error {
	resSession, err := handle(context.Background(), &proto.Session{
		UserId:           uint64(sess.UserID),
		Cookie:           sess.Cookie,
		CookieExpiration: sess.CookieExpiration.Format(time.RFC3339),
	})

	if err != nil {
		return err
	}

	sess.UserID = uint(resSession.GetUserId())
	sess.Cookie = resSession.GetCookie()
	parsedTime, err := time.Parse(time.RFC3339, resSession.CookieExpiration)
	if err != nil {
		return err
	}
	sess.CookieExpiration = parsedTime
	return nil
}

func (sm *GrpcSession) Create(sess *Session) error {
	return sm.handleSession(sess, sm.authService.CreateSession)
}

func (sm *GrpcSession) Check(sess *Session) error {
	return sm.handleSession(sess, sm.authService.CheckSession)
}

func (sm *GrpcSession) Delete(sess *Session) error {
	return sm.handleSession(sess, sm.authService.DeleteSession)
}
