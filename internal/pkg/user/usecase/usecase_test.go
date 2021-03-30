package usecase

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/user"
	"Redioteka/internal/pkg/user/repository/mock"
	"crypto/sha256"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

type getByIdTestCase struct {
	ID      uint
	outUser domain.User
	outErr  error
}

var getByIdTests = []getByIdTestCase{
	{
		ID: 1,
		outUser: domain.User{
			ID:       1,
			Username: "user",
			Email:    "mail@mail.ru",
		},
		outErr: nil,
	},
	{
		ID:      2,
		outUser: domain.User{},
		outErr:  user.NotFoundError,
	},
}

func TestUserUsecase_GetById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mock.NewMockUserRepository(ctrl)
	uc := NewUserUsecase(repoMock)

	for _, test := range getByIdTests {
		t.Run(fmt.Sprintf("ID: %v", test.ID),
			func(t *testing.T) {
				repoMock.EXPECT().GetById(test.ID).Times(1).Return(test.outUser, test.outErr)
				currentUser, currentErr := uc.GetById(test.ID)
				if currentErr != nil {
					require.Equal(t, currentErr, test.outErr)
				} else {
					require.Equal(t, currentUser, test.outUser)
				}
			})
	}
}

type signupTestCase struct {
	inUser  *domain.User
	outId   uint
	outUser domain.User
	outErr  error
}

var signupTests = []signupTestCase{
	{
		inUser: &domain.User{
			Username:             "user",
			Email:                "mail@mail.ru",
			InputPassword:        "password1",
			ConfirmInputPassword: "password1",
		},
		outId: 1,
		outUser: domain.User{
			ID:       1,
			Username: "user",
			Email:    "mail@mail.ru",
		},
		outErr: nil,
	},
	{
		inUser: &domain.User{
			Username:             "user",
			Email:                "mail@mail.ru",
			InputPassword:        "password1",
			ConfirmInputPassword: "not_password1",
		},
		outErr: user.InvalidCredentials,
	},
	{
		inUser: &domain.User{
			Username:             "user",
			Email:                "mail@mail.ru",
			InputPassword:        "password1",
			ConfirmInputPassword: "password1",
		},
		outId:  1,
		outErr: user.NotFoundError,
	},
	{
		inUser: &domain.User{
			Username:             "already_added_user",
			Email:                "already_added_email@mail.ru",
			InputPassword:        "password1",
			ConfirmInputPassword: "password1",
		},
		outErr: user.AlreadyAddedError,
	},
}

func TestUserUsecase_Signup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock := mock.NewMockUserRepository(ctrl)
	uc := NewUserUsecase(userRepoMock)

	for testId, test := range signupTests {
		t.Run(fmt.Sprintln(testId), func(t *testing.T) {
			if test.outErr == user.AlreadyAddedError {
				userRepoMock.EXPECT().Store(test.inUser).Times(1).Return(test.outId, test.outErr)
			} else if test.outErr == user.NotFoundError {
				userRepoMock.EXPECT().Store(test.inUser).Times(1).Return(test.outId, nil)
				userRepoMock.EXPECT().GetById(test.outId).Times(1).Return(test.outUser, test.outErr)
			} else if test.outErr == nil {
				userRepoMock.EXPECT().Store(test.inUser).Times(1).Return(test.outId, nil)
				userRepoMock.EXPECT().GetById(test.outId).Times(1).Return(test.outUser, nil)
			}
			currentUser, currentErr := uc.Signup(test.inUser)
			if test.outErr != nil {
				require.Equal(t, currentErr, test.outErr)
			} else {
				require.Equal(t, currentUser, test.outUser)
			}
		})
	}
}

type loginTestCase struct {
	inUser  *domain.User
	outUser domain.User
	outErr  error
}

var loginTests = []loginTestCase{
	{
		inUser: &domain.User{
			Email:         "mail@mail.ru",
			InputPassword: "password1",
		},
		outUser: domain.User{
			ID:       1,
			Username: "user",
			Email:    "mail@mail.ru",
			Password: sha256.Sum256([]byte("password1")),
		},
		outErr: nil,
	},
	{
		inUser: &domain.User{
			Email:         "mail@mail.ru",
			InputPassword: "",
		},
		outErr: user.InvalidForm,
	},
	{
		inUser: &domain.User{
			Email:         "",
			InputPassword: "password",
		},
		outErr: user.InvalidForm,
	},
	{
		inUser: &domain.User{
			Email:         "not_found_mail@mail.ru",
			InputPassword: "password1",
		},
		outErr: user.NotFoundError,
	},
	{
		inUser: &domain.User{
			Email:         "mail@mail.ru",
			InputPassword: "wrong_password",
		},
		outUser: domain.User{
			ID:       1,
			Username: "user",
			Email:    "mail@mail.ru",
			Password: sha256.Sum256([]byte("password1")),
		},
		outErr: user.InvalidCredentials,
	},
}

func TestUserUsecase_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock := mock.NewMockUserRepository(ctrl)
	uc := NewUserUsecase(userRepoMock)

	for testId, test := range loginTests {
		t.Run(fmt.Sprintln(testId, test.outErr), func(t *testing.T) {
			if test.outErr == user.NotFoundError {
				userRepoMock.EXPECT().GetByEmail(test.inUser.Email).Times(1).Return(test.outUser, test.outErr)
			} else if test.outErr == nil || test.outErr == user.InvalidCredentials {
				userRepoMock.EXPECT().GetByEmail(test.inUser.Email).Times(1).Return(test.outUser, nil)
			}
			currentUser, currentErr := uc.Login(test.inUser)
			if test.outErr != nil {
				require.Equal(t, currentErr, test.outErr)
			} else {
				require.Equal(t, currentUser, test.outUser)
			}
		})
	}
}

type updateTestCase struct {
	inUpdate *domain.User
	outErr   error
}

var updateTests = []updateTestCase{
	{
		inUpdate: &domain.User{
			Email:    "new_mail@mail.ru",
			Username: "new_user",
		},
		outErr: nil,
	},
	{
		inUpdate: &domain.User{
			Email:    "",
			Username: "",
			Avatar:   "",
		},
		outErr: user.InvalidUpdateError,
	},
}

func TestUserUsecase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepoMock := mock.NewMockUserRepository(ctrl)
	uc := NewUserUsecase(userRepoMock)

	for testId, test := range updateTests {
		t.Run(fmt.Sprintln(testId, test.outErr), func(t *testing.T) {
			if test.outErr == nil {
				userRepoMock.EXPECT().Update(test.inUpdate).Times(1).Return(nil)
			}
			currentErr := uc.Update(test.inUpdate)
			require.Equal(t, currentErr, test.outErr)
		})
	}
}
