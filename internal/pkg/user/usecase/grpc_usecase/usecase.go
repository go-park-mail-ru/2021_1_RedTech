package grpc_usecase

import (
	"Redioteka/internal/pkg/authorization/delivery/grpc/proto"
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/utils/session"
	"io"
)

type grpcUserUsecase struct {
	userRepo    domain.UserRepository
	authService proto.AuthorizationClient
	avatarRepo  domain.AvatarRepository
}

func (g grpcUserUsecase) GetById(id uint) (domain.User, error) {
	panic("implement me")
}

func (g grpcUserUsecase) Signup(u *domain.User) (domain.User, error) {
	panic("implement me")
}

func (g grpcUserUsecase) Login(u *domain.User) (domain.User, error) {
	panic("implement me")
}

func (g grpcUserUsecase) Logout(sess *session.Session) error {
	panic("implement me")
}

func (g grpcUserUsecase) Update(u *domain.User) error {
	panic("implement me")
}

func (g grpcUserUsecase) Delete(id uint) error {
	panic("implement me")
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
