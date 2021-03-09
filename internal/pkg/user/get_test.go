package user

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

type getTestCase struct {
	ID      string
	outJSON string
	status  int
}

var testCaseGet = []getTestCase{
	{
		ID:      "1",
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
	testUser := `{"username":"good_user","email":"gmail@mail.ru","password":"pass","confirm_password":"pass"}` + "\n"
	api.Signup(httptest.NewRecorder(),
		httptest.NewRequest("POST", "/api/users/signup", bytes.NewReader([]byte(testUser))))
	for _, test := range testCaseGet {
		t.Run(fmt.Sprintf("IN: %v, OUT: %v, CODE: %v", test.ID, test.outJSON, test.status),
			func(t *testing.T) {
				test.outJSON += "\n"
				r := httptest.NewRequest("GET", "/api/users/"+test.ID, nil)
				r = mux.SetURLVars(r, map[string]string{"id": test.ID})
				w := httptest.NewRecorder()
				api.Get(w, r)
				current := getTestCase{
					ID:      test.ID,
					outJSON: w.Body.String(),
					status:  w.Code,
				}
				require.Equal(t, test, current)
			})
	}
}

func TestMe(t *testing.T) {
	require.Equal(t, true, true)
}
