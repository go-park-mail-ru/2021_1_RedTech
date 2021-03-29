package http

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/user"
	mock2 "Redioteka/internal/pkg/user/usecase/mock"
	"Redioteka/internal/pkg/utils/jsonerrors"
	"Redioteka/internal/pkg/utils/session"
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

var usersTestData = map[uint]domain.User{
	1: {
		ID:       1,
		Email:    "mail1@mail.ru",
		Username: "user1",
		Password: sha256.Sum256([]byte("pass1")),
	},
	2: {
		ID:       2,
		Email:    "mail2@mail.ru",
		Username: "user2",
		Password: sha256.Sum256([]byte("pass2")),
	},
}

var userTestErrors = map[uint]error{
	3: user.InvalidCredentials,
}

type TestCaseGet struct {
	ID      uint
	outJSON string
	status  int
}

var testCaseGet = []TestCaseGet{
	{
		ID:      1,
		outJSON: `{"id":1,"email":"mail1@mail.ru","username":"user1"}`,
		status:  http.StatusOK,
	},
	{
		ID:      2,
		outJSON: `{"id":2,"username":"user2"}`,
		status:  http.StatusOK,
	},
	{
		ID:      3,
		outJSON: jsonerrors.JSONMessage("get"),
		status:  http.StatusNotAcceptable,
	},
}

func TestUserHandler_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uCaseMock := mock2.NewMockUserUsecase(ctrl)
	for uid, value := range usersTestData {
		uCaseMock.EXPECT().GetById(uid).Times(1).Return(value, nil)
	}
	for uid, errValue := range userTestErrors {
		uCaseMock.EXPECT().GetById(uid).Times(1).Return(domain.User{}, errValue)
	}
	handler := &UserHandler{
		UUsecase: uCaseMock,
	}
	for _, test := range testCaseGet {
		t.Run(fmt.Sprintf("IN: %v, OUT: %v, CODE: %v", test.ID, test.outJSON, test.status),
			func(t *testing.T) {
				test.outJSON += "\n"
				r := httptest.NewRequest("GET", "/api/users/"+strconv.Itoa(int(test.ID)), nil)
				r = mux.SetURLVars(r, map[string]string{"id": strconv.Itoa(int(test.ID))})
				w := httptest.NewRecorder()

				err := session.Create(w, r, 1)
				defer session.Delete(w, r, 1)
				require.NoError(t, err)

				handler.Get(w, r)
				current := TestCaseGet{
					ID:      test.ID,
					outJSON: w.Body.String(),
					status:  w.Code,
				}
				require.Equal(t, test, current)
			})
	}
}

var meTests = []TestCaseGet{
	{
		// 1 is authorized
		ID:      1,
		outJSON: `{"id":1}`,
		status:  http.StatusOK,
	},
	{
		// 2 is unauthorized
		ID:      2,
		outJSON: jsonerrors.JSONMessage("unauthorized"),
		status:  http.StatusUnauthorized,
	},
}

func TestUserHandler_Me(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uCaseMock := mock2.NewMockUserUsecase(ctrl)
	handler := &UserHandler{
		UUsecase: uCaseMock,
	}
	for _, test := range meTests {
		t.Run(fmt.Sprintf("IN: %v, OUT: %v, CODE: %v", test.ID, test.outJSON, test.status),
			func(t *testing.T) {
				test.outJSON += "\n"
				r := httptest.NewRequest("GET", "/api/users/"+strconv.Itoa(int(test.ID)), nil)
				r = mux.SetURLVars(r, map[string]string{"id": strconv.Itoa(int(test.ID))})
				w := httptest.NewRecorder()

				if test.ID == 1 {
					err := session.Create(w, r, 1)
					require.NoError(t, err)
					defer session.Delete(w, r, 1)
				}

				handler.Me(w, r)
				current := TestCaseGet{
					ID:      test.ID,
					outJSON: w.Body.String(),
					status:  w.Code,
				}
				require.Equal(t, test, current)
			})
	}
}

type userTestCase struct {
	inJSON  string
	inUser  domain.User
	outUser domain.User
	outJSON string
	status  int
}

var updateTests = []userTestCase{
	{
		inJSON: `{"email":"newmail1@mail.ru"}` + "\n",
		inUser: domain.User{
			ID:    1,
			Email: "newmail1@mail.ru",
		},
		outUser: domain.User{
			ID:       1,
			Email:    "newmail1@mail.ru",
			Username: "user1",
		},
		outJSON: `{"id":1,"email":"newmail1@mail.ru","username":"user1"}`,
		status:  http.StatusOK,
	},
	{
		inJSON: `{"username":"new_user1"}` + "\n",
		inUser: domain.User{
			ID:       1,
			Username: "new_user1",
		},
		outUser: domain.User{
			ID:       1,
			Username: "new_user1",
			Email:    "mail1@mail.ru",
		},
		outJSON: `{"id":1,"email":"mail1@mail.ru","username":"new_user1"}`,
		status:  http.StatusOK,
	},
	{
		inJSON: `{"email":"newmail1@mail.ru","username":"new_user1"}` + "\n",
		inUser: domain.User{
			ID:       1,
			Email:    "newmail1@mail.ru",
			Username: "new_user1",
		},
		outUser: domain.User{
			ID:       1,
			Email:    "newmail1@mail.ru",
			Username: "new_user1",
		},
		outJSON: `{"id":1,"email":"newmail1@mail.ru","username":"new_user1"}`,
		status:  http.StatusOK,
	},
	{
		inJSON: `{"email":"","username":""}` + "\n",
		inUser: domain.User{
			ID:       1,
			Email:    "",
			Username: "",
		},
		outUser: domain.User{},
		outJSON: `{"message":"invalid update"}`,
		status:  http.StatusBadRequest,
	},
	{
		inJSON: `{}` + "\n",
		inUser: domain.User{
			ID: 1,
		},
		outUser: domain.User{},
		outJSON: `{"message":"invalid update"}`,
		status:  http.StatusBadRequest,
	},
}

func TestUserHandler_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uCaseMock := mock2.NewMockUserUsecase(ctrl)
	handler := &UserHandler{
		UUsecase: uCaseMock,
	}
	for _, test := range updateTests {
		t.Run(fmt.Sprintf("IN: %v, OUT: %v, CODE: %v", test.inJSON, test.outJSON, test.status),
			func(t *testing.T) {
				test.outJSON += "\n"
				body := bytes.NewReader([]byte(test.inJSON))
				r := httptest.NewRequest("PATCH", "/api/users/1", body)
				r = mux.SetURLVars(r, map[string]string{"id": "1"})
				w := httptest.NewRecorder()

				err := session.Create(w, r, 1)
				require.NoError(t, err)
				defer session.Delete(w, r, 1)
				if test.status == http.StatusOK {
					uCaseMock.EXPECT().Update(&test.inUser).Times(1).Return(nil)
					uCaseMock.EXPECT().GetById(uint(1)).Times(1).Return(test.outUser, nil)
				} else {
					uCaseMock.EXPECT().Update(&test.inUser).Times(1).Return(user.InvalidUpdateError)
				}
				handler.Update(w, r)
				current := userTestCase{
					inJSON:  test.inJSON,
					inUser:  test.inUser,
					outUser: test.outUser,
					outJSON: w.Body.String(),
					status:  w.Code,
				}
				require.Equal(t, test, current)
			})
	}
}

var signupTests = []userTestCase{
	{
		inJSON:  "",
		inUser:  domain.User{},
		outUser: domain.User{},
		outJSON: `{"message":"json decode"}`,
		status:  http.StatusBadRequest,
	},
	{
		inJSON:  `{"bad_field":"bad_value"}`,
		inUser:  domain.User{},
		outUser: domain.User{},
		outJSON: `{"message":"signup"}`,
		status:  http.StatusNotAcceptable,
	},
	{
		inJSON: `{"username":"a","email":"ya@mail.ru","password":"pass","confirm_password":"not_pass"}`,
		inUser: domain.User{
			Username:             "a",
			Email:                "ya@mail.ru",
			InputPassword:        "pass",
			ConfirmInputPassword: "not_pass",
		},
		outUser: domain.User{},
		outJSON: `{"message":"signup"}`,
		status:  http.StatusNotAcceptable,
	},
	{
		inJSON:  "{}",
		inUser:  domain.User{},
		outUser: domain.User{},
		outJSON: `{"message":"signup"}`,
		status:  http.StatusNotAcceptable},
	{
		inJSON: `{"username":"good_user","email":"gmail@mail.ru","password":"pass","confirm_password":"pass"}`,
		inUser: domain.User{
			Email:                "gmail@mail.ru",
			Username:             "good_user",
			InputPassword:        "pass",
			ConfirmInputPassword: "pass",
		},
		outUser: domain.User{
			ID:       1,
			Email:    "gmail@mail.ru",
			Username: "good_user",
		},
		outJSON: `{"id":1,"email":"gmail@mail.ru","username":"good_user"}`,
		status:  http.StatusOK,
	},
}

func TestUserHandler_Signup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uCaseMock := mock2.NewMockUserUsecase(ctrl)
	handler := &UserHandler{
		UUsecase: uCaseMock,
	}
	for _, test := range signupTests {
		t.Run(fmt.Sprintf("IN: %v, OUT: %v, CODE: %v", test.inJSON, test.outJSON, test.status),
			func(t *testing.T) {
				test.outJSON += "\n"
				body := bytes.NewReader([]byte(test.inJSON))
				r := httptest.NewRequest("POST", "/api/users/signup", body)
				w := httptest.NewRecorder()
				if test.status == http.StatusOK {
					uCaseMock.EXPECT().Signup(&test.inUser).Times(1).Return(test.outUser, nil)
				} else if test.status != http.StatusBadRequest {
					uCaseMock.EXPECT().Signup(&test.inUser).Times(1).Return(domain.User{}, user.InvalidCredentials)
				}
				handler.Signup(w, r)
				current := userTestCase{
					inJSON:  test.inJSON,
					inUser:  test.inUser,
					outUser: test.outUser,
					outJSON: w.Body.String(),
					status:  w.Code,
				}
				require.Equal(t, current, test)
			})
	}
}

var loginTests = []userTestCase{
	{
		inJSON:  "",
		outJSON: `{"message":"json decode"}`,
		status:  http.StatusBadRequest,
	},
	{
		inJSON:  `{"bad_field":"bad_value"}`,
		inUser:  domain.User{},
		outJSON: `{"message":"login"}`,
		status:  http.StatusNotAcceptable,
	},
	{
		inJSON:  "{}",
		outJSON: `{"message":"login"}`,
		status:  http.StatusNotAcceptable,
	},
}

func TestUserHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uCaseMock := mock2.NewMockUserUsecase(ctrl)
	handler := &UserHandler{
		UUsecase: uCaseMock,
	}
	for _, test := range loginTests {
		t.Run(fmt.Sprintf("IN: %v, OUT: %v, CODE: %v", test.inJSON, test.outJSON, test.status),
			func(t *testing.T) {
				test.outJSON += "\n"
				body := bytes.NewReader([]byte(test.inJSON))
				r := httptest.NewRequest("POST", "/api/users/login", body)
				w := httptest.NewRecorder()
				if test.status == http.StatusOK {
					uCaseMock.EXPECT().Login(&test.inUser).Times(1).Return(test.outUser, nil)
				} else if test.status != http.StatusBadRequest {
					uCaseMock.EXPECT().Login(&test.inUser).Times(1).Return(domain.User{}, user.InvalidCredentials)
				}
				handler.Login(w, r)
				current := userTestCase{
					inJSON:  test.inJSON,
					inUser:  test.inUser,
					outUser: test.outUser,
					outJSON: w.Body.String(),
					status:  w.Code,
				}
				require.Equal(t, current, test)
			})
	}
}

type avatarTestCase struct {
	inURL       string
	inURLParams map[string]string
	outJSON     string
	status      int
}

var avatarTests = []avatarTestCase{
	{
		inURL:       "/api/users//avatar",
		inURLParams: map[string]string{},
		outJSON:     `{"message":"url var"}`,
		status:      http.StatusBadRequest,
	},
	{
		inURL:       "/api/users/2/avatar",
		inURLParams: map[string]string{"id": "2"},
		outJSON:     `{"message":"wrong id"}`,
		status:      http.StatusForbidden,
	},
}

func TestUserHandler_Avatar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uCaseMock := mock2.NewMockUserUsecase(ctrl)
	handler := &UserHandler{
		UUsecase: uCaseMock,
	}
	for _, test := range avatarTests {
		t.Run(fmt.Sprintf("IN: %v, OUT: %v, CODE: %v", test.inURL, test.outJSON, test.status),
			func(t *testing.T) {
				test.outJSON += "\n"
				body := bytes.NewReader([]byte(""))
				r := httptest.NewRequest("POST", test.inURL, body)
				r = mux.SetURLVars(r, test.inURLParams)
				w := httptest.NewRecorder()

				err := session.Create(w, r, 1)
				require.NoError(t, err)
				defer session.Delete(w, r, 1)

				handler.Avatar(w, r)
				current := avatarTestCase{
					inURL:       test.inURL,
					inURLParams: test.inURLParams,
					outJSON:     w.Body.String(),
					status:      w.Code,
				}
				require.Equal(t, current, test)
			})
	}
}
