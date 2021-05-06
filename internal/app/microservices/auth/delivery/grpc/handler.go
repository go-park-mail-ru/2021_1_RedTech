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

func (handler *authorizationHandler) SignIn(context.Context, *proto.SignInCredentials) (*proto.User, error) {
	handler.userUsecase.Signup
}
func (handler *authorizationHandler) SignUp(context.Context, *proto.SignupCredentials) (*proto.User, error) {
	return nil, nil
}

// Session things
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
