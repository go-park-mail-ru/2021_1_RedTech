package user

import (
	"Redioteka/internal/app/session"
	"Redioteka/internal/app/user"
	"crypto/sha256"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

var usersTestData = []user.User{
	{
		ID:       123,
		Email:    "gmail@mail.ru",
		Username: "good_user",
		Password: sha256.Sum256([]byte("pass")),
	},
	{
		ID:       124,
		Email:    "mail@mail.ru",
		Username: "user_user",
		Password: sha256.Sum256([]byte("pass")),
	},
}

func fillTestData() {
	for i := range usersTestData {
		user.data.addUser(&usersTestData[i])
	}
}

func clearTestData() {
	for _, user := range usersTestData {
		user.data.deleteById(user.ID)
	}
}

type TestCaseGet struct {
	ID      string
	outJSON string
	status  int
}

var testCaseGet = []TestCaseGet{
	{
		ID:      "123",
		outJSON: `{"id":123,"email":"gmail@mail.ru","username":"good_user"}`,
		status:  http.StatusOK,
	},
	{
		ID:      "124",
		outJSON: `{"id":124,"username":"user_user"}`,
		status:  http.StatusOK,
	},
	{
		ID:      "100",
		outJSON: `{"error":"server can't send user'"}`,
		status:  http.StatusBadRequest,
	},
}

func TestGet(t *testing.T) {
	api := &user.Handler{}
	fillTestData()
	defer clearTestData()

	for _, test := range testCaseGet {
		t.Run(fmt.Sprintf("IN: %v, OUT: %v, CODE: %v", test.ID, test.outJSON, test.status),
			func(t *testing.T) {
				test.outJSON += "\n"
				r := httptest.NewRequest("GET", "/api/users/"+test.ID, nil)
				r = mux.SetURLVars(r, map[string]string{"id": test.ID})
				w := httptest.NewRecorder()

				err := session.Create(w, r, 123)
				defer session.Delete(w, r, 123)
				require.Equal(t, nil, err)

				api.Get(w, r)
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
	// authorize only for id 123
	{
		ID:      "123",
		outJSON: `{"id":123}`,
		status:  http.StatusOK,
	},
	{
		ID:      "124",
		outJSON: `{"error":"can't find user"}`,
		status:  http.StatusBadRequest,
	},
}

func TestMe(t *testing.T) {
	api := &user.Handler{}
	fillTestData()
	defer clearTestData()

	for _, test := range testCaseMe {
		t.Run(fmt.Sprintf("IN: %v, OUT: %v, CODE: %v", test.ID, test.outJSON, test.status),
			func(t *testing.T) {
				test.outJSON += "\n"
				r := httptest.NewRequest("GET", "/api/me", nil)
				w := httptest.NewRecorder()

				if test.ID == "123" {
					err := session.Create(w, r, 123)
					defer session.Delete(w, r, 123)
					require.Equal(t, nil, err)
				}

				api.Me(w, r)

				current := TestCaseGet{
					ID:      test.ID,
					outJSON: w.Body.String(),
					status:  w.Code,
				}
				require.Equal(t, test, current)
			})
	}
}
