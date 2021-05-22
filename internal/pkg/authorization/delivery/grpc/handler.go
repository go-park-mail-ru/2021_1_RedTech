package grpc

import (
	proto2 "Redioteka/internal/pkg/authorization/delivery/grpc/proto"
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/utils/log"
	"Redioteka/internal/pkg/utils/session"
	"context"
	"fmt"
	"time"
)

type authorizationHandler struct {
	proto2.UnimplementedAuthorizationServer
	authUsecase    domain.AuthorizationUsecase
	sessionManager session.SessionManager
}

func NewAuthorizationHandler(uucase domain.AuthorizationUsecase, sManager session.SessionManager) proto2.AuthorizationServer {
	return &authorizationHandler{
		authUsecase:    uucase,
		sessionManager: sManager,
	}
}

func (handler *authorizationHandler) Delete(ctx context.Context, id *proto2.UserId) (*proto2.EmptyMessage, error) {
	err := handler.authUsecase.Delete(uint(id.GetId()))
	if err != nil {
		return nil, err
	}
	return &proto2.EmptyMessage{}, nil
}

func toProto(user domain.User) *proto2.User {
	return &proto2.User{
		Id:       uint64(user.ID),
		Username: user.Username,
		Email:    user.Email,
		Avatar:   user.Avatar,
	}
}

func (handler *authorizationHandler) Login(ctx context.Context, credentials *proto2.LoginCredentials) (*proto2.User, error) {
	user, err := handler.authUsecase.Login(&domain.User{
		Email:         credentials.GetEmail(),
		InputPassword: credentials.GetPassword(),
	})
	if err != nil {
		return nil, err
	}
	return toProto(user), nil
}

func (handler *authorizationHandler) Signup(ctx context.Context, credentials *proto2.SignupCredentials) (*proto2.User, error) {
	user, err := handler.authUsecase.Signup(&domain.User{
		Username:             credentials.GetUsername(),
		Email:                credentials.GetEmail(),
		InputPassword:        credentials.GetPassword(),
		ConfirmInputPassword: credentials.GetConfirmPassword(),
	})
	if err != nil {
		return nil, err
	}
	return toProto(user), nil
}

func (handler *authorizationHandler) Update(ctx context.Context, info *proto2.UpdateInfo) (*proto2.EmptyMessage, error) {
	err := handler.authUsecase.Update(&domain.User{
		ID:       uint(info.GetUserId()),
		Email:    info.GetEmail(),
		Avatar:   info.GetAvatar(),
		Username: info.GetUsername(),
	})
	if err != nil {
		return nil, err
	}
	return &proto2.EmptyMessage{}, nil
}

func (handler *authorizationHandler) GetById(ctx context.Context, userId *proto2.UserId) (*proto2.User, error) {
	res, err := handler.authUsecase.GetById(uint(userId.Id))
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Error while getting with id: %v", userId.Id))
		return nil, err
	}

	return &proto2.User{
		Id:       uint64(res.ID),
		Username: res.Username,
		Email:    res.Email,
		Avatar:   res.Avatar,
	}, nil
}

func parseProtoUser(user *proto2.User) *domain.User {
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

func parseProtoSession(protoSession *proto2.Session) (*session.Session, error) {
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

func (handler *authorizationHandler) CreateSession(ctx context.Context, credentials *proto2.CreateSessionParams) (*proto2.Session, error) {
	sess := &session.Session{
		UserID: uint(credentials.GetUserId()),
	}
	err := handler.sessionManager.Create(sess)
	if err != nil {
		return nil, err
	}
	return &proto2.Session{
		UserId:           uint64(sess.UserID),
		Cookie:           sess.Cookie,
		CookieExpiration: sess.CookieExpiration.Format(time.RFC3339),
	}, nil
}

func (handler *authorizationHandler) DeleteSession(ctx context.Context, credentials *proto2.Session) (*proto2.DeleteSessionInfo, error) {
	sess, err := parseProtoSession(credentials)
	if err != nil {
		return nil, err
	}

	err = handler.sessionManager.Delete(sess)
	if err != nil {
		return nil, err
	}
	return &proto2.DeleteSessionInfo{}, nil
}

func (handler *authorizationHandler) CheckSession(ctx context.Context, credentials *proto2.Session) (*proto2.CheckSessionInfo, error) {
	sess, err := parseProtoSession(credentials)
	if err != nil {
		return nil, err
	}

	err = handler.sessionManager.Check(sess)
	if err != nil {
		return nil, err
	}
	return &proto2.CheckSessionInfo{}, nil
}
