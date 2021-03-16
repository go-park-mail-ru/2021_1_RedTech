package user

import (
	"Redioteka/internal/pkg/session"
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testCaseUpdate = []TestCase{
	{
		inJSON:  `{"email":"emaaail@mail.ru","password":"pass"}` + "\n",
		outJSON: `{"username":"good_user","email":"emaaail@mail.ru"}`,
		status:  http.StatusOK,
	},
	{
		inJSON:  `{"username":"very_good_user","password":"pass"}` + "\n",
		outJSON: `{"username":"very_good_user","email":"emaaail@mail.ru"}`,
		status:  http.StatusOK,
	},
	{
		inJSON:  `{"email":"gates@gmail.com","username":"very_good_user","password":"pass"}` + "\n",
		outJSON: `{"username":"very_good_user","email":"gates@gmail.com"}`,
		status:  http.StatusOK,
	},
	{
		inJSON:  `{"email":"","username":"","password":"","new_password":"","confirm_new_password":""}` + "\n",
		outJSON: `{"error":"error while updating user"}`,
		status:  http.StatusBadRequest,
	},
	{
		inJSON:  `{}` + "\n",
		outJSON: `{"error":"error while updating user"}`,
		status:  http.StatusBadRequest,
	},
}

func TestUpdate(t *testing.T) {
	api := &Handler{}
	fillTestData()
	defer clearTestData()

	for _, test := range testCaseUpdate {
		t.Run(fmt.Sprintf("IN: %v, OUT: %v, CODE: %v", test.inJSON, test.outJSON, test.status),
			func(t *testing.T) {
				test.outJSON += "\n"
				body := bytes.NewReader([]byte(test.inJSON))
				r := httptest.NewRequest("GET", "/api/me", body)
				w := httptest.NewRecorder()

				err := session.Create(w, r, 123)
				defer session.Delete(w, r, 123)
				require.Equal(t, nil, err)
				api.Update(w, r)

				current := TestCase{
					inJSON:  test.inJSON,
					outJSON: w.Body.String(),
					status:  w.Code,
				}
				// todo delete this
				fmt.Println("CURRENT   ", current.outJSON)
				require.Equal(t, test, current)
			})
	}
}
