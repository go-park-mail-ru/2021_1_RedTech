package http

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/movie"
	mock2 "Redioteka/internal/pkg/movie/usecase/mock"
	"bytes"
	"encoding/json"
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
		Title:       "Some japanese comedy",
		Description: "Test data",
		Rating:      9,
		Countries:   []string{"Japan", "South Korea"},
		IsFree:      false,
		Genres:      []string{"Comedy"},
		Actors:      []string{"Sono", "Chi", "No", "Sadame", "Mina"},
		Avatar:      "/static/movies/default.jpg",
		Type:        domain.MovieT,
		Year:        "2011",
		Director:    []string{"Director Directorovich"},
	},
	2: {
		ID:          2,
		Title:       "KOOOOOOOREAAAAA",
		Description: "Test data",
		Rating:      2,
		Countries:   []string{"USA", "South Korea"},
		IsFree:      true,
		Genres:      []string{"Comedy"},
		Actors:      []string{"John", "Wick"},
		Avatar:      "/static/movies/default.jpg",
		Type:        domain.MovieT,
		Year:        "1999",
		Director:    []string{"Test Cameron"},
	},
	3: {
		ID:          3,
		Title:       "Samurai vs Samurai",
		Description: "Test data",
		Rating:      2,
		Countries:   []string{"Japan", "USA"},
		IsFree:      false,
		Genres:      []string{"Action", "Thriller"},
		Actors:      []string{"Actor1 ", "Actor2", "Actor3"},
		Avatar:      "/static/movies/default.jpg",
		Type:        domain.MovieT,
		Year:        "3000",
		Director:    []string{"Stephen Spielberg"},
	},
	4: {
		ID:          4,
		Title:       "Series 1-1",
		Description: "Test data",
		Rating:      2,
		Countries:   []string{"Russia"},
		IsFree:      true,
		Genres:      []string{"Horror", "Comedy"},
		Actors:      []string{"Ivan", "Nikolay", "Berezka"},
		Avatar:      "/static/movies/default.jpg",
		Type:        domain.SeriesT,
		Year:        "3000",
		Director:    []string{"Mikhalkov"},
	},
	5: {
		ID:          5,
		Title:       "Russische Leute in Deutschland",
		Description: "Test data",
		Rating:      7,
		Countries:   []string{"Germany", "Russia"},
		IsFree:      false,
		Genres:      []string{"Cartoon"},
		Actors:      []string{"Fluegegeheimen", "Alexandr", "Leshiy Ivanov"},
		Avatar:      "/static/movies/default.jpg",
		Type:        domain.SeriesT,
		Year:        "1222",
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
		outJSON:  `{"message":"url params"}`,
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
		outJSON:  `{"id":1,"title":"Some japanese comedy","description":"Test data","rating":9,"countries":["Japan","South Korea"],"is_free":false,"genres":["Comedy"],"actors":["Sono","Chi","No","Sadame","Mina"],"movie_avatar":"/static/movies/default.jpg","type":"movie","year":"2011","director":["Director Directorovich"]}`,
		outMovie: movieTestData[1],
		status:   http.StatusOK,
	},
}

func TestMovieHandler_Get(t *testing.T) {
	res, _ := json.Marshal(movieTestData[1])
	fmt.Println(string(res))
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mCaseMock := mock2.NewMockMovieUsecase(ctrl)
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

