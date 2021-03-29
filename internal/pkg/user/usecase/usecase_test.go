package usecase

import (
	"Redioteka/internal/pkg/domain"
	mock2 "Redioteka/internal/pkg/user/repository/mock"
	"crypto/sha256"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserUsecase_GetById(t *testing.T) {
	retUser := domain.User{
		ID:       123,
		Username: "user",
		Email:    "mail@mail.ru",
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mock2.NewMockUserRepository(ctrl)
	repoMock.EXPECT().GetById(retUser.ID).Times(1).Return(retUser, nil)

	uc := NewUserUsecase(repoMock)

	user, err := uc.GetById(retUser.ID)
	require.NoError(t, err)
	require.Equal(t, user, retUser)
}

func TestUserUsecase_Signup(t *testing.T) {
	signupForm := &domain.User{
		Username:             "user",
		Email:                "mail@mail.ru",
		InputPassword:        "password1",
		ConfirmInputPassword: "password1",
	}

	retUser := domain.User{
		ID:       1,
		Username: "user",
		Email:    "mail@mail.ru",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock := mock2.NewMockUserRepository(ctrl)

	userRepoMock.EXPECT().Store(signupForm).Times(1).Return(retUser.ID, nil)
	userRepoMock.EXPECT().GetById(retUser.ID).Times(1).Return(retUser, nil)

	uc := NewUserUsecase(userRepoMock)
	user, err := uc.Signup(signupForm)
	require.NoError(t, err)
	require.Equal(t, user, retUser)
}

func TestUserUsecase_Login(t *testing.T) {
	loginForm := &domain.User{
		Email:         "mail@mail.ru",
		InputPassword: "password1",
	}

	retUser := domain.User{
		ID:       1,
		Username: "user",
		Email:    "mail@mail.ru",
		Password: sha256.Sum256([]byte("password1")),
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock := mock2.NewMockUserRepository(ctrl)

	userRepoMock.EXPECT().GetByEmail(retUser.Email).Times(1).Return(retUser, nil)

	uc := NewUserUsecase(userRepoMock)

	user, err := uc.Login(loginForm)
	require.NoError(t, err)
	require.Equal(t, user, retUser)
}

func TestUserUsecase_Update(t *testing.T) {
	userUpdateForm := &domain.User{
		ID:    1,
		Email: "new_mail@mail.ru",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock := mock2.NewMockUserRepository(ctrl)

	userRepoMock.EXPECT().Update(userUpdateForm).Times(1).Return(nil)

	uc := NewUserUsecase(userRepoMock)

	err := uc.Update(userUpdateForm)
	require.NoError(t, err)
}
