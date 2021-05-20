package grpc

import (
	"Redioteka/internal/app/auth/delivery/grpc/proto"
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/utils/log"
	"Redioteka/internal/pkg/utils/session"
	"context"
	"fmt"
	"time"
)

type authorizationHandler struct {
	proto.UnimplementedAuthorizationServer
	userUsecase    domain.UserUsecase
}

func NewAuthorizationHandler(uucase domain.UserUsecase) proto.AuthorizationServer {
	return &authorizationHandler{
		userUsecase: uucase,
	}
}

func (handler *authorizationHandler) GetById(ctx context.Context, userId *proto.UserId) (*proto.User, error) {
	res, err := handler.userUsecase.GetById(uint(userId.Id))
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Error while getting with id: %v", userId.Id))
		return nil, err
	}

	return &proto.User{
		Id:       uint32(res.ID),
		Username: res.Username,
		Email:    res.Username,
		Avatar:   res.Avatar,
	}, nil
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

func (handler *authorizationHandler) SignIn(ctx context.Context, credentials *proto.SignInCredentials) (*proto.User, error) {
	user, err := handler.userUsecase.Login(&domain.User{
		Email:         credentials.GetEmail(),
		InputPassword: credentials.GetPassword(),
	})
	if err != nil {
		return nil, err
	}
	return &proto.User{
		Id:       uint32(user.ID),
		Username: user.Username,
		Email:    user.Email,
		Avatar:   user.Avatar,
	}, err
}

func (handler *authorizationHandler) SignUp(ctx context.Context, credentials *proto.SignupCredentials) (*proto.User, error) {
	user, err := handler.userUsecase.Signup(&domain.User{
		Username:             credentials.GetUsername(),
		Email:                credentials.GetEmail(),
		InputPassword:        credentials.GetPassword(),
		ConfirmInputPassword: credentials.GetConfirmPassword(),
	})
	if err != nil {
		return nil, err
	}
	return &proto.User{
		Id:       uint32(user.ID),
		Username: user.Username,
		Email:    user.Email,
		Avatar:   user.Avatar,
	}, nil
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
