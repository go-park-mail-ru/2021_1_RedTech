package grpc

import (
	"Redioteka/internal/app/microservices/auth/delivery/grpc/proto"
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/utils/log"
	"Redioteka/internal/pkg/utils/session"
	"context"
	"fmt"
	"time"
)

type authorizationHandler struct {
	userUsecase    domain.UserUsecase
	sessionManager session.SessionManager
}

func NewAuthorizationHandler(uucase domain.UserUsecase, manager session.SessionManager) authorizationHandler {
	return authorizationHandler{
		userUsecase:    uucase,
		sessionManager: manager,
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
		Email:         credentials.Email,
		InputPassword: credentials.Password,
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
		Username:             credentials.Username,
		Email:                credentials.Email,
		InputPassword:        credentials.Password,
		ConfirmInputPassword: credentials.ConfirmPassword,
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
		UserID: uint(credentials.UserId),
	}
	err := handler.sessionManager.Create(sess)
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

	err = handler.sessionManager.Delete(sess)
	if err != nil {
		return nil, err
	}
	return &proto.DeleteSessionInfo{}, nil
}

func (handler *authorizationHandler) CheckSession(context.Context, *proto.Session) (*proto.CheckSessionInfo, error) {
	return nil, nil
}

func (handler *authorizationHandler) mustEmbedUnimplementedAuthorizationServer() {
}
