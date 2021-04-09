package http

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/movie"
	mock2 "Redioteka/internal/pkg/movie/usecase/mock"
	"Redioteka/internal/pkg/utils/session"
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
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

type movieGetTestCase struct {
	inURL    string
	inParams map[string]string
	outJSON  string
	outMovie domain.Movie
	status   int
}

var movieGetTests = []movieGetTestCase{
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
		outJSON:  `{"id":1,"title":"Film","description":"Test data","rating":9,"countries":["Japan","South Korea"],"is_free":false,"genres":["Comedy"],"actors":["Sana","Momo","Mina"],"movie_avatar":"/static/movies/default.jpg","type":"movie","year":"2012","director":["James Cameron"]}`,
		outMovie: movieTestData[1],
		status:   http.StatusOK,
	},
}

func TestMovieHandler_Get(t *testing.T) {
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
				r := httptest.NewRequest("GET", test.inURL, nil)
				r = mux.SetURLVars(r, test.inParams)
				w := httptest.NewRecorder()

				if test.status == http.StatusOK {
					mCaseMock.EXPECT().GetById(uint(1)).Times(1).Return(movieTestData[1], nil)
				} else if test.status == http.StatusNotFound {
					mCaseMock.EXPECT().GetById(uint(2)).Times(1).Return(domain.Movie{}, movie.NotFoundError)
				}

				handler.Get(w, r)
				current := movieGetTestCase{
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

type movieSetFavouriteTestCase struct {
	inURL       string
	inRouteName string
	inParams    map[string]string
	err         error
	status      int
}

var movieSetFavouriteTests = []movieSetFavouriteTestCase{
	{
		inURL:       "/api/media/movie//like",
		inRouteName: addFavourite,
		inParams:    map[string]string{},
		err:         nil,
		status:      http.StatusBadRequest,
	},
	{
		inURL:       "/api/media/movie/1/dislike",
		inRouteName: removeFavourite,
		inParams:    map[string]string{"id": "1"},
		err:         nil,
		status:      http.StatusUnauthorized,
	},
	{
		inURL:       "/api/media/movie/2/like",
		inRouteName: addFavourite,
		inParams:    map[string]string{"id": "2"},
		err:         movie.NotFoundError,
		status:      http.StatusNotFound,
	},
	{
		inURL:       "/api/media/movie/3/dislike",
		inRouteName: removeFavourite,
		inParams:    map[string]string{"id": "3"},
		err:         nil,
		status:      http.StatusOK,
	},
}

func TestMovieHandler_SetFavourite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mCaseMock := mock2.NewMockMovieUsecase(ctrl)
	handler := &MovieHandler{
		MUCase: mCaseMock,
	}

	for _, test := range movieSetFavouriteTests {
		t.Run(fmt.Sprintf("IN: %v CODE: %v", test.inURL, test.status),
			func(t *testing.T) {
				body := bytes.NewReader([]byte(""))
				r := httptest.NewRequest("POST", test.inURL, body)
				r = mux.SetURLVars(r, test.inParams)
				w := httptest.NewRecorder()

				if test.status == http.StatusOK || test.status == http.StatusNotFound {
					id, err := strconv.Atoi(test.inParams["id"])
					require.NoError(t, err)
					sess := &session.Session{UserID: uint(id)}
					err = session.Manager.Create(sess)
					require.NoError(t, err)
					sess = &session.Session{Cookie: sess.Cookie}
					r.AddCookie(&http.Cookie{
						Name:    "session_id",
						Value:   sess.Cookie,
						Expires: sess.CookieExpiration,
					})
					defer session.Manager.Delete(sess)

					if test.inRouteName == addFavourite {
						mCaseMock.EXPECT().AddFavourite(uint(id), sess).Times(1).Return(test.err)
					} else if test.inRouteName == removeFavourite {
						mCaseMock.EXPECT().RemoveFavourite(uint(id), sess).Times(1).Return(test.err)
					}
				}

				handler.SetFavourite(w, r)
				require.Equal(t, test.status, w.Code)
			})
	}
}
