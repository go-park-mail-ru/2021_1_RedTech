package user

import (
	"Redioteka/internal/pkg/session"
	"crypto/sha256"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

var usersTestData = []User{
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
	for i, _ := range usersTestData {
		data.addUser(&usersTestData[i])
	}
}

func clearTestData() {
	for _, user := range usersTestData {
		data.deleteById(user.ID)
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
		outJSON: `{"username":"good_user","email":"gmail@mail.ru"}`,
		status:  http.StatusOK,
	},
	{
		ID:      "100",
		outJSON: `{"error":"server can't send user'"}`,
		status:  http.StatusBadRequest,
	},
}

func TestGet(t *testing.T) {
	api := &Handler{}
	fillTestData()
	defer clearTestData()

	for _, test := range testCaseGet {
		t.Run(fmt.Sprintf("IN: %v, OUT: %v, CODE: %v", test.ID, test.outJSON, test.status),
			func(t *testing.T) {
				test.outJSON += "\n"
				r := httptest.NewRequest("GET", "/api/users/"+test.ID, nil)
				r = mux.SetURLVars(r, map[string]string{"id": test.ID})
				w := httptest.NewRecorder()
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
	{
		ID:      "123",
		outJSON: `{"username":"good_user","email":"gmail@mail.ru"}`,
		status:  http.StatusOK,
	},
}

func TestMe(t *testing.T) {
	api := &Handler{}
	fillTestData()
	defer clearTestData()

	for _, test := range testCaseMe {
		t.Run(fmt.Sprintf("IN: %v, OUT: %v, CODE: %v", test.ID, test.outJSON, test.status),
			func(t *testing.T) {
				test.outJSON += "\n"
				r := httptest.NewRequest("GET", "/api/me", nil)
				w := httptest.NewRecorder()

				err := session.Create(w, r, 123)
				defer session.Delete(w, r, 123)
				require.Equal(t, nil, err)
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