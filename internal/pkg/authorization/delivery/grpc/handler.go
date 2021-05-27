package grpc

import (
	"Redioteka/internal/pkg/authorization/delivery/grpc/proto"
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/utils/log"
	"Redioteka/internal/pkg/utils/session"
	"context"
	"fmt"
	"time"
)

type authorizationHandler struct {
	proto.UnimplementedAuthorizationServer
	authUsecase    domain.AuthorizationUsecase
	sessionManager session.SessionManager
}

func NewAuthorizationHandler(uucase domain.AuthorizationUsecase, sManager session.SessionManager) proto.AuthorizationServer {
	return &authorizationHandler{
		authUsecase:    uucase,
		sessionManager: sManager,
	}
}

func (handler *authorizationHandler) Delete(ctx context.Context, id *proto.UserId) (*proto.EmptyMessage, error) {
	err := handler.authUsecase.Delete(uint(id.GetId()))
	if err != nil {
		return nil, err
	}
	return &proto.EmptyMessage{}, nil
}

func toProto(user domain.User) *proto.User {
	return &proto.User{
		Id:           uint64(user.ID),
		Username:     user.Username,
		Email:        user.Email,
		Avatar:       user.Avatar,
		IsSubscriber: user.IsSubscriber,
	}
}

func (handler *authorizationHandler) Login(ctx context.Context, credentials *proto.LoginCredentials) (*proto.User, error) {
	user, err := handler.authUsecase.Login(&domain.User{
		Email:         credentials.GetEmail(),
		InputPassword: credentials.GetPassword(),
	})
	if err != nil {
		return nil, err
	}
	return toProto(user), nil
}

func (handler *authorizationHandler) Signup(ctx context.Context, credentials *proto.SignupCredentials) (*proto.User, error) {
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

func (handler *authorizationHandler) Update(ctx context.Context, info *proto.UpdateInfo) (*proto.EmptyMessage, error) {
	err := handler.authUsecase.Update(&domain.User{
		ID:       uint(info.GetUserId()),
		Email:    info.GetEmail(),
		Avatar:   info.GetAvatar(),
		Username: info.GetUsername(),
	})
	if err != nil {
		return nil, err
	}
	return &proto.EmptyMessage{}, nil
}

func (handler *authorizationHandler) GetById(ctx context.Context, userId *proto.UserId) (*proto.User, error) {
	res, err := handler.authUsecase.GetById(uint(userId.Id))
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Error while getting with id: %v", userId.Id))
		return nil, err
	}
	log.Log.Info(fmt.Sprint("GetById: ", userId.Id))

	protoUser := &proto.User{
		Id:           uint64(res.ID),
		Username:     res.Username,
		Email:        res.Email,
		Avatar:       res.Avatar,
		IsSubscriber: res.IsSubscriber,
	}
	log.Log.Info(fmt.Sprint("GetById: ", protoUser))
	return protoUser, nil
}

func parseProtoSession(protoSession *proto.Session) (*session.Session, error) {
	parsedTime, err := time.Parse(time.RFC3339, protoSession.CookieExpiration)
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

func (handler *authorizationHandler) CreateSession(ctx context.Context, credentials *proto.Session) (*proto.Session, error) {
	sess := &session.Session{
		UserID: uint(credentials.GetUserId()),
		// we need only user id for creation
	}
	err := handler.sessionManager.Create(sess)
	if err != nil {
		return nil, err
	}
	return &proto.Session{
		UserId:           uint64(sess.UserID),
		Cookie:           sess.Cookie,
		CookieExpiration: sess.CookieExpiration.Format(time.RFC3339),
	}, nil
}

func (handler *authorizationHandler) DeleteSession(ctx context.Context, credentials *proto.Session) (*proto.Session, error) {
	sess, err := parseProtoSession(credentials)
	if err != nil {
		return nil, err
	}

	err = handler.sessionManager.Delete(sess)
	if err != nil {
		return nil, err
	}
	return &proto.Session{
		UserId:           credentials.GetUserId(),
		Cookie:           credentials.GetCookie(),
		CookieExpiration: credentials.GetCookieExpiration(),
	}, nil
}

func (handler *authorizationHandler) CheckSession(ctx context.Context, credentials *proto.Session) (*proto.Session, error) {
	sess, err := parseProtoSession(credentials)
	if err != nil {
		return nil, err
	}

	err = handler.sessionManager.Check(sess)
	if err != nil {
		return nil, err
	}
	return &proto.Session{
		UserId:           uint64(sess.UserID),
		Cookie:           credentials.GetCookie(),
		CookieExpiration: credentials.GetCookieExpiration(),
	}, nil
}
