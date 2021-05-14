package grpc

import (
	"Redioteka/internal/app/microservices/auth/delivery/grpc/proto"
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/utils/log"
	"Redioteka/internal/pkg/utils/session"
	"context"
	"fmt"
)

type authorizationHandler struct {
	userUsecase    domain.UserUsecase
	sessionManager session.SessionManager
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
		Username: credentials.Username,
		Email: credentials.Email,
		InputPassword: credentials.Password,
		ConfirmInputPassword: credentials.ConfirmPassword,
	})
	if err != nil {
		return nil, err
	}
	return &proto.User{
		Id: uint32(user.ID),
		Username: user.Username,
		Email: user.Email,
		Avatar: user.Avatar,
	}, nil
}

func (handler *authorizationHandler) CreateSession(context.Context, *proto.CreateSessionParams) (*proto.SessionId, error) {
	return nil, nil
}

func (handler *authorizationHandler) DeleteSession(context.Context, *proto.SessionId) (*proto.DeleteSessionInfo, error) {
	return nil, nil
}

func (handler *authorizationHandler) CheckSession(context.Context, *proto.SessionId) (*proto.CheckSessionInfo, error) {
	return nil, nil
}

func (handler *authorizationHandler) mustEmbedUnimplementedAuthorizationServer() {
}
