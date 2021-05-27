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

func ToProto(sess *Session) *proto.Session {
	return &proto.Session{
		UserId:           uint64(sess.UserID),
		Cookie:           sess.Cookie,
		CookieExpiration: sess.CookieExpiration.Format(time.RFC3339),
	}
}

func FromProto(sess *Session, protoSession *proto.Session) error {
	sess.UserID = uint(protoSession.GetUserId())
	sess.Cookie = protoSession.GetCookie()
	parsedTime, err := time.Parse(time.RFC3339, protoSession.CookieExpiration)
	if err != nil {
		return err
	}
	sess.CookieExpiration = parsedTime
	return nil
}

func NewGrpcSession(aClient proto.AuthorizationClient) SessionManager {
	return &GrpcSession{
		authService: aClient,
	}
}

func (sm *GrpcSession) handleSession(sess *Session, handle func(ctx context.Context, in *proto.Session, opts ...grpc.CallOption) (*proto.Session, error)) error {
	resSession, err := handle(context.Background(), ToProto(sess))
	if err != nil {
		return err
	}
	return FromProto(sess, resSession)
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
