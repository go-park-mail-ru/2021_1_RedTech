package movie

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestCase struct {
	inURL   string
	outJSON string
	status  int
}

func TestHandlerGet(t *testing.T) {
	tests := []TestCase{
		{"/api/media/movie/2", `{"error":"not found"}`, http.StatusNotFound},
		{"/api/media/movie/1",
			`{"id":1,"title":"Film","description":"Test data","rating":9,"countries":["Japan","South Korea"],"is_free":false,"genres":["Comedy"],"actors":["Sana","Momo","Mina"]}`,
			http.StatusOK},
	}

	api := &Handler{}
	for i, test := range tests {
		test.outJSON += "\n"
		fmt.Println("TestMovieGet", i)
		body := bytes.NewBuffer(nil)
		r := httptest.NewRequest("GET", test.inURL, body)
		w := httptest.NewRecorder()
		api.Get(w, r)
		actual := TestCase{
			inURL:   test.inURL,
			outJSON: w.Body.String(),
			status:  w.Code,
		}
		require.Equal(t, test, actual)
	}
}
