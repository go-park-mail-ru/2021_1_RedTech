package grpc_usecase

import (
	"Redioteka/internal/pkg/authorization/delivery/grpc/proto"
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/user"
	"Redioteka/internal/pkg/utils/session"
	"context"
	"io"
)

type grpcUserUsecase struct {
	userRepo       domain.UserRepository
	authService    proto.AuthorizationClient
	sessionManager session.SessionManager
	avatarRepo     domain.AvatarRepository
}

func NewGrpcUserUsecase(u domain.UserRepository, a domain.AvatarRepository,
	s proto.AuthorizationClient, sm session.SessionManager) domain.UserUsecase {
	return &grpcUserUsecase{
		userRepo:       u,
		avatarRepo:     a,
		authService:    s,
		sessionManager: sm,
	}
}

func (g grpcUserUsecase) GetById(id uint) (domain.User, error) {
	foundUser, err := g.authService.GetById(context.Background(), &proto.UserId{
		Id: uint64(id),
	})
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		ID:           uint(foundUser.GetId()),
		Username:     foundUser.GetUsername(),
		Email:        foundUser.GetEmail(),
		Avatar:       foundUser.GetAvatar(),
		IsSubscriber: foundUser.GetIsSubscriber(),
	}, nil
}

func (g grpcUserUsecase) Signup(u *domain.User) (domain.User, error) {
	foundUser, err := g.authService.Signup(context.Background(), &proto.SignupCredentials{
		Username:        u.Username,
		Email:           u.Email,
		Password:        u.InputPassword,
		ConfirmPassword: u.ConfirmInputPassword,
	})
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		ID:       uint(foundUser.GetId()),
		Username: foundUser.GetUsername(),
		Email:    foundUser.GetEmail(),
		Avatar:   foundUser.GetAvatar(),
	}, nil
}

func (g grpcUserUsecase) Login(u *domain.User) (domain.User, error) {
	foundUser, err := g.authService.Login(context.Background(), &proto.LoginCredentials{
		Email:    u.Email,
		Password: u.InputPassword,
	})
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		ID:       uint(foundUser.GetId()),
		Username: foundUser.GetUsername(),
		Email:    foundUser.GetEmail(),
		Avatar:   foundUser.GetAvatar(),
	}, nil
}

func (g grpcUserUsecase) Logout(sess *session.Session) error {
	resSession, err := g.authService.DeleteSession(context.Background(), session.ToProto(sess))
	if err != nil {
		return err
	}
	return session.FromProto(sess, resSession)
}

func (g grpcUserUsecase) Update(u *domain.User) error {
	_, err := g.authService.Update(context.Background(), &proto.UpdateInfo{
		UserId:   uint64(u.ID),
		Username: u.Username,
		Email:    u.Email,
		Avatar:   u.Avatar,
	})
	if err != nil {
		return err
	}
	return nil
}

func (g grpcUserUsecase) Delete(id uint) error {
	_, err := g.authService.Delete(context.Background(), &proto.UserId{
		Id: uint64(id),
	})
	if err != nil {
		return err
	}
	return nil
}

func (g grpcUserUsecase) GetFavourites(id uint, sess *session.Session) ([]domain.Movie, error) {
	err := g.sessionManager.Check(sess)
	if err != nil {
		return nil, user.UnauthorizedError
	}
	if sess.UserID != id {
		return nil, user.InvalidCredentials
	}

	return g.userRepo.GetFavouritesByID(id)
}

func (g grpcUserUsecase) UploadAvatar(reader io.Reader, path, ext string) (string, error) {
	return g.avatarRepo.UploadAvatar(reader, path, ext)
}
