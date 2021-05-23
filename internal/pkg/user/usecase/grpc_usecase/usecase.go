package grpc_usecase

import (
	"Redioteka/internal/pkg/authorization/delivery/grpc/proto"
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/utils/session"
	"context"
	"io"
)

type grpcUserUsecase struct {
	userRepo    domain.UserRepository
	authService proto.AuthorizationClient
	avatarRepo  domain.AvatarRepository
}

func (g grpcUserUsecase) GetById(id uint) (domain.User, error) {
	user, err := g.authService.GetById(context.Background(), &proto.UserId{
		Id: uint64(id),
	})
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		ID:       uint(user.GetId()),
		Username: user.GetUsername(),
		Email:    user.GetEmail(),
		Avatar:   user.GetAvatar(),
	}, nil
}

func (g grpcUserUsecase) Signup(u *domain.User) (domain.User, error) {
	user, err := g.authService.Signup(context.Background(), &proto.SignupCredentials{
		Username:        u.Username,
		Email:           u.Email,
		Password:        u.InputPassword,
		ConfirmPassword: u.ConfirmInputPassword,
	})
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		ID:       uint(user.GetId()),
		Username: user.GetUsername(),
		Email:    user.GetEmail(),
		Avatar:   user.GetAvatar(),
	}, nil
}

func (g grpcUserUsecase) Login(u *domain.User) (domain.User, error) {
	user, err := g.authService.Login(context.Background(), &proto.LoginCredentials{
		Email:    u.Email,
		Password: u.InputPassword,
	})
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		ID:       uint(user.GetId()),
		Username: user.GetUsername(),
		Email:    user.GetEmail(),
		Avatar:   user.GetAvatar(),
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
	panic("implement me")
}

func (g grpcUserUsecase) UploadAvatar(reader io.Reader, path, ext string) (string, error) {
	panic("implement me")
}

func NewGrpcUserUsecase(u domain.UserRepository, a domain.AvatarRepository, s proto.AuthorizationClient) domain.UserUsecase {
	return &grpcUserUsecase{
		userRepo:    u,
		avatarRepo:  a,
		authService: s,
	}
}
