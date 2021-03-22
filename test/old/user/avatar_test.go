package user

import (
	"Redioteka/internal/pkg/user"
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

type TestIn struct {
	URL       string
	URLParams map[string]string
}

type TestOut struct {
	JSON   string
	status int
}

type TestCaseAvatar struct {
	in  TestIn
	out TestOut
}

func TestAvatar(t *testing.T) {
	tests := []TestCaseAvatar{
		{
			TestIn{"/api/users//avatar", map[string]string{}},
			TestOut{`{"error":"server"}`, http.StatusInternalServerError},
		},
		{
			TestIn{"/api/users/2/avatar", map[string]string{"id": "2"}},
			TestOut{`{"error":"can't find user'"}`, http.StatusBadRequest},
		},
		/*{
			TestIn{"/api/media/movie/1", map[string]string{"id": "1"}},
			TestOut{`{"id":1,"title":"Film","description":"Test data","rating":9,"countries":["Japan","South Korea"],"is_free":false,"genres":["Comedy"],"actors":["Sana","Momo","Mina"],"movie_avatar":"/static/movies/default.jpg"}`,
				http.StatusOK},
		},*/
	}

	api := &user.Handler{}
	for i, test := range tests {
		test.out.JSON += "\n"
		fmt.Println("TestAvatar", i)
		body := bytes.NewReader([]byte(""))
		r := httptest.NewRequest("POST", test.in.URL, body)
		r = mux.SetURLVars(r, test.in.URLParams)
		w := httptest.NewRecorder()
		api.Avatar(w, r)
		actual := TestOut{
			JSON:   w.Body.String(),
			status: w.Code,
		}
		require.Equal(t, test.out, actual)
	}
}
