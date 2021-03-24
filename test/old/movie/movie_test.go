package movie

import (
	"Redioteka/internal/pkg/movie"
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

type TestCase struct {
	in  TestIn
	out TestOut
}

func TestHandlerGet(t *testing.T) {
	tests := []TestCase{
		{
			TestIn{"/api/media/movie/", map[string]string{}},
			TestOut{`{"error":"server"}`, http.StatusInternalServerError},
		},
		{
			TestIn{"/api/media/movie/2", map[string]string{"id": "2"}},
			TestOut{`{"error":"not found"}`, http.StatusNotFound},
		},
		{
			TestIn{"/api/media/movie/1", map[string]string{"id": "1"}},
			TestOut{`{"id":1,"title":"Film","description":"Test data","rating":9,"countries":["Japan","South Korea"],"is_free":false,"genres":["Comedy"],"actors":["Sana","Momo","Mina"],"movie_avatar":"/static/movies/default.jpg","type":"movie","year":"2012","director":["James Cameron"]}`,
				http.StatusOK},
		},
	}

	api := &movie.Handler{}
	for i, test := range tests {
		test.out.JSON += "\n"
		fmt.Println("TestMovieGet", i)
		r := httptest.NewRequest("GET", test.in.URL, nil)
		r = mux.SetURLVars(r, test.in.URLParams)
		w := httptest.NewRecorder()
		api.Get(w, r)
		actual := TestOut{
			JSON:   w.Body.String(),
			status: w.Code,
		}
		require.Equal(t, test.out, actual)
	}
}
