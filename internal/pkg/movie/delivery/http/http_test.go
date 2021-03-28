package http

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/movie"
	"Redioteka/internal/pkg/movie/mock"
	"bytes"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

var movieTestData = map[uint]domain.Movie{
	1: {
		ID:          1,
		Title:       "Film",
		Description: "Test data",
		Rating:      9,
		Countries:   []string{"Japan", "South Korea"},
		IsFree:      false,
		Genres:      []string{"Comedy"},
		Actors:      []string{"Sana", "Momo", "Mina"},
		Avatar:      "/static/movies/default.jpg",
		Type:        domain.MovieT,
		Year:        "2012",
		Director:    []string{"James Cameron"},
	},
}

type movieTestCase struct {
	inURL    string
	inParams map[string]string
	outJSON  string
	outMovie domain.Movie
	status   int
}

var movieGetTests = []movieTestCase{
	{
		inURL:    "/api/media/movie/",
		inParams: map[string]string{},
		outJSON:  `{"message":"params"}`,
		status:   http.StatusBadRequest,
	},
	{
		inURL:    "/api/media/movie/2",
		inParams: map[string]string{"id": "2"},
		outJSON:  `{"message":"not found"}`,
		status:   http.StatusNotFound,
	},
	{
		inURL:    "/api/media/movie/1",
		inParams: map[string]string{"id": "1"},
		outJSON:  `{"id":1,"title":"Film","description":"Test data","rating":9,"countries":["Japan","South Korea"],"is_free":false,"genres":["Comedy"],"actors":["Sana","Momo","Mina"],"movie_avatar":"/static/movies/default.jpg","type":"movie","year":"2012","director":["James Cameron"]}`,
		outMovie: movieTestData[1],
		status:   http.StatusOK,
	},
}

func TestUserHandler_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mCaseMock := mock.NewMockMovieUsecase(ctrl)
	handler := &MovieHandler{
		MUCase: mCaseMock,
	}
	for _, test := range movieGetTests {
		t.Run(fmt.Sprintf("IN: %v, OUT: %v, CODE: %v", test.inURL, test.outJSON, test.status),
			func(t *testing.T) {
				test.outJSON += "\n"
				body := bytes.NewReader([]byte(""))
				r := httptest.NewRequest("GET", test.inURL, body)
				r = mux.SetURLVars(r, test.inParams)
				w := httptest.NewRecorder()

				if test.status == http.StatusOK {
					mCaseMock.EXPECT().GetById(uint(1)).Times(1).Return(movieTestData[1], nil)
				} else if test.status == http.StatusNotFound {
					mCaseMock.EXPECT().GetById(uint(2)).Times(1).Return(domain.Movie{}, movie.NotFoundError)
				}

				handler.Get(w, r)
				current := movieTestCase{
					inURL:    test.inURL,
					inParams: test.inParams,
					outMovie: test.outMovie,
					outJSON:  w.Body.String(),
					status:   w.Code,
				}
				require.Equal(t, current, test)
			})
	}
}
