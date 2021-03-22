package user

import (
	"Redioteka/internal/app/session"
	"Redioteka/internal/app/user"
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testCaseUpdate = []TestCase{
	{
		inJSON:  `{"email":"emaaail@mail.ru"}` + "\n",
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

func TestUpdate(t *testing.T) {
	api := &user.Handler{}
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
				require.Equal(t, test, current)
			})
	}
}
