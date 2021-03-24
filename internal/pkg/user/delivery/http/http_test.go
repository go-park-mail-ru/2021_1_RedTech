package http

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/user"
	"Redioteka/internal/pkg/user/mock"
	"Redioteka/internal/pkg/utils/jsonerrors"
	"Redioteka/internal/pkg/utils/session"
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
	inJSON  string
	outJSON string
	status  int
}

var testCaseUpdate = []UpdateTestCase{
	{
		inJSON:  `{"email":"newmail1@mail.ru"}` + "\n",
		outJSON: `{"id":123,"email":"emaaail@mail.ru","username":"good_user"}`,
		status:  http.StatusOK,
	},
	{
		inJSON:  `{"username":"very_good_user"}` + "\n",
		outJSON: `{"id":123,"email":"emaaail@mail.ru","username":"very_good_user"}`,
		status:  http.StatusOK,
	},
	{
		inJSON:  `{"email":"gates@gmail.com","username":"very_very_good_user"}` + "\n",
		outJSON: `{"id":123,"email":"gates@gmail.com","username":"very_very_good_user"}`,
		status:  http.StatusOK,
	},
	{
		inJSON:  `{"email":"","username":""}` + "\n",
		outJSON: `{"error":"error while updating user"}`,
		status:  http.StatusBadRequest,
	},
	{
		inJSON:  `{}` + "\n",
		outJSON: `{"error":"error while updating user"}`,
		status:  http.StatusBadRequest,
	},
}

func TestUserHandler_Update(t *testing.T) {
}

func TestUserHandler_Login(t *testing.T) {

}

func TestUserHandler_Signup(t *testing.T) {

}

func TestUserHandler_Avatar(t *testing.T) {

}
