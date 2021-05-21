package grpc

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/user/delivery/grpc/proto"
	"Redioteka/internal/pkg/utils/log"
	"Redioteka/internal/pkg/utils/session"
	"context"
	"fmt"
	"time"
)

type authorizationHandler struct {
	proto.UnimplementedAuthorizationServer
	authUsecase domain.AuthorizationUsecase
}

func NewAuthorizationHandler(uucase domain.AuthorizationUsecase) proto.AuthorizationServer {
	return &authorizationHandler{
		authUsecase: uucase,
	}
}

func (handler *authorizationHandler) GetById(ctx context.Context, userId *proto.UserId) (*proto.User, error) {
	res, err := handler.authUsecase.GetById(uint(userId.Id))
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Error while getting with id: %v", userId.Id))
		return nil, err
	}

	return &proto.User{
		Id:       uint32(res.ID),
		Username: res.Username,
		Email:    res.Email,
		Avatar:   res.Avatar,
	}, nil
}

func (handler *authorizationHandler) GetByEmail(ctx context.Context, email *proto.Email) (*proto.User, error) {
	res, err := handler.authUsecase.GetByEmail(email.GetEmail())
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Error while getting with email: %v", email.GetEmail()))
		return nil, err
	}

	return &proto.User{
		Id:       uint32(res.ID),
		Username: res.Username,
		Email:    res.Email,
		Avatar:   res.Avatar,
	}, nil
}

func parseProtoUser(user *proto.User) *domain.User {
	if user == nil {
		return nil
	}
	return &domain.User{
		ID:       uint(user.GetId()),
		Username: user.GetUsername(),
		Email:    user.GetEmail(),
		Avatar:   user.GetAvatar(),
	}
}

func (handler *authorizationHandler) Update(ctx context.Context, user *proto.User) (*proto.Nothing, error) {
	err := handler.authUsecase.Update(parseProtoUser(user))
	if err != nil {
		return nil, err
	}
	return &proto.Nothing{}, nil
}

func (handler *authorizationHandler) Store(ctx context.Context, user *proto.User) (*proto.UserId, error) {
	id, err := handler.authUsecase.Store(parseProtoUser(user))
	if err != nil {
		return nil, err
	}
	return &proto.UserId{
		Id: uint64(id),
	}, nil
}

func (handler *authorizationHandler) Delete(ctx context.Context, id *proto.UserId) (*proto.Nothing, error) {
	err := handler.authUsecase.Delete(uint(id.GetId()))
	if err != nil {
		return nil, err
	}
	return &proto.Nothing{}, nil
}
func parseProtoSession(protoSession *proto.Session) (*session.Session, error) {
	parsedTime, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", protoSession.CookieExpiration)
	if err != nil {
		return nil, err
	}

	sess := &session.Session{
		UserID:           uint(protoSession.UserId),
		Cookie:           protoSession.Cookie,
		CookieExpiration: parsedTime,
	}
	return sess, nil
}

func (handler *authorizationHandler) CreateSession(ctx context.Context, credentials *proto.CreateSessionParams) (*proto.Session, error) {
	sess := &session.Session{
		UserID: uint(credentials.GetUserId()),
	}
	err := session.Manager.Create(sess)
	if err != nil {
		return nil, err
	}
	return &proto.Session{
		UserId:           uint32(sess.UserID),
		Cookie:           sess.Cookie,
		CookieExpiration: sess.CookieExpiration.String(),
	}, nil
}

func (handler *authorizationHandler) DeleteSession(ctx context.Context, credentials *proto.Session) (*proto.DeleteSessionInfo, error) {
	sess, err := parseProtoSession(credentials)
	if err != nil {
		return nil, err
	}

	err = session.Manager.Delete(sess)
	if err != nil {
		return nil, err
	}
	return &proto.DeleteSessionInfo{}, nil
}

func (handler *authorizationHandler) CheckSession(ctx context.Context, credentials *proto.Session) (*proto.CheckSessionInfo, error) {
	sess, err := parseProtoSession(credentials)
	if err != nil {
		return nil, err
	}

	err = session.Manager.Check(sess)
	if err != nil {
		return nil, err
	}
	return &proto.CheckSessionInfo{}, nil
}

func (handler *authorizationHandler) mustEmbedUnimplementedAuthorizationServer() {
}
