package http

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/user"
	"Redioteka/internal/pkg/user/mock"
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
		status:  http.StatusBadRequest,
	},
}

func TestUserHandler_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uCaseMock := mock.NewMockUserUsecase(ctrl)
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

var testCaseMe = []TestCaseGet{
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
	uCaseMock := mock.NewMockUserUsecase(ctrl)
	handler := &UserHandler{
		UUsecase: uCaseMock,
	}
	for _, test := range testCaseMe {
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

type UpdateTestCase struct {
	inJSON    string
	inUpdate  domain.User
	outUpdate domain.User
	outJSON   string
	status    int
}

var testCaseUpdate = []UpdateTestCase{
	{
		inJSON: `{"email":"newmail1@mail.ru"}` + "\n",
		inUpdate: domain.User{
			ID:    1,
			Email: "newmail1@mail.ru",
		},
		outUpdate: domain.User{
			ID:       1,
			Email:    "newmail1@mail.ru",
			Username: "user1",
		},
		outJSON: `{"id":1,"email":"newmail1@mail.ru","username":"user1"}`,
		status:  http.StatusOK,
	},
	{
		inJSON: `{"username":"new_user1"}` + "\n",
		inUpdate: domain.User{
			ID:       1,
			Username: "new_user1",
		},
		outUpdate: domain.User{
			ID:       1,
			Username: "new_user1",
			Email:    "mail1@mail.ru",
		},
		outJSON: `{"id":1,"email":"mail1@mail.ru","username":"new_user1"}`,
		status:  http.StatusOK,
	},
	{
		inJSON: `{"email":"newmail1@mail.ru","username":"new_user1"}` + "\n",
		inUpdate: domain.User{
			ID:       1,
			Email:    "newmail1@mail.ru",
			Username: "new_user1",
		},
		outUpdate: domain.User{
			ID:       1,
			Email:    "newmail1@mail.ru",
			Username: "new_user1",
		},
		outJSON: `{"id":1,"email":"newmail1@mail.ru","username":"new_user1"}`,
		status:  http.StatusOK,
	},
	{
		inJSON: `{"email":"","username":""}` + "\n",
		inUpdate: domain.User{
			ID:       1,
			Email:    "",
			Username: "",
		},
		outUpdate: domain.User{},
		outJSON:   `{"message":"invalid update"}`,
		status:    http.StatusBadRequest,
	},
	{
		inJSON: `{}` + "\n",
		inUpdate: domain.User{
			ID: 1,
		},
		outUpdate: domain.User{},
		outJSON:   `{"message":"invalid update"}`,
		status:    http.StatusBadRequest,
	},
}

func TestUserHandler_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uCaseMock := mock.NewMockUserUsecase(ctrl)
	handler := &UserHandler{
		UUsecase: uCaseMock,
	}
	for _, test := range testCaseUpdate {
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
					uCaseMock.EXPECT().Update(&test.inUpdate).Times(1).Return(nil)
					uCaseMock.EXPECT().GetById(uint(1)).Times(1).Return(test.outUpdate, nil)
				} else {
					uCaseMock.EXPECT().Update(&test.inUpdate).Times(1).Return(user.InvalidUpdateError)
				}
				handler.Update(w, r)
				current := UpdateTestCase{
					inJSON:   test.inJSON,
					inUpdate: test.inUpdate,
					outUpdate: test.outUpdate,
					outJSON:  w.Body.String(),
					status:   w.Code,
				}
				require.Equal(t, test, current)
			})
	}
}

func TestUserHandler_Login(t *testing.T) {
}

func TestUserHandler_Signup(t *testing.T) {

}

func TestUserHandler_Avatar(t *testing.T) {

}
