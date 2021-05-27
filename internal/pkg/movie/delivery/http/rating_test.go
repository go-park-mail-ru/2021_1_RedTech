package http

import (
	"Redioteka/internal/pkg/movie/usecase/mock"
	"Redioteka/internal/pkg/utils/session"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

type movieRateTestCase struct {
	inURL      string
	inParams   map[string]string
	outJSON    string
	authorized bool
	status     int
}

/*
var movieLikeTests = []movieRateTestCase{
	{
		inURL:    "/api/media/movie/1/like",
		inParams: map[string]string{},
		outJSON:  `{"message":"url params"}`,
		status:   http.StatusBadRequest,
	},
	{
		inURL:      "/api/media/movie/1/like",
		inParams:   map[string]string{"id": "1"},
		outJSON:    `{"message":"unauthorized"}`,
		status:     http.StatusUnauthorized,
		authorized: false,
	},
	{
		inURL:      "/api/media/movie/1/like",
		inParams:   map[string]string{"id": "1"},
		outJSON:    `{"message":"OK"}`,
		status:     http.StatusOK,
		authorized: true,
	},
}

func TestMovieHandler_Like(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mCaseMock := mock.NewMockMovieUsecase(ctrl)
	handler := &MovieHandler{
		MUCase: mCaseMock,
	}
	for _, test := range movieLikeTests {
		t.Run(fmt.Sprintf("IN: %v, OUT: %v, CODE: %v", test.inURL, test.outJSON, test.status),
			func(t *testing.T) {
				test.outJSON += "\n"
				r := httptest.NewRequest("POST", test.inURL, nil)
				r = mux.SetURLVars(r, test.inParams)
				w := httptest.NewRecorder()

				if test.authorized {
					sess := &session.Session{UserID: 1}
					err := session.Manager.Create(sess)
					require.NoError(t, err)
					defer sessDelete(session.Manager, sess)
					r.AddCookie(&http.Cookie{
						Name:    "session_id",
						Value:   sess.Cookie,
						Expires: sess.CookieExpiration,
					})
				}

				if test.status == http.StatusOK {
					mCaseMock.EXPECT().Like(uint(1), uint(1)).Times(1).Return(nil)
				}

				handler.Like(w, r)
				require.Equal(t, test.status, w.Code)
				require.Equal(t, test.outJSON, w.Body.String())
			})
	}
}
*/
var movieDislikeTests = []movieRateTestCase{
	{
		inURL:    "/api/media/movie/1/dislike",
		inParams: map[string]string{},
		outJSON:  `{"message":"url params"}`,
		status:   http.StatusBadRequest,
	},
	{
		inURL:      "/api/media/movie/1/dislike",
		inParams:   map[string]string{"id": "1"},
		outJSON:    `{"message":"unauthorized"}`,
		status:     http.StatusUnauthorized,
		authorized: false,
	},
	{
		inURL:      "/api/media/movie/1/dislike",
		inParams:   map[string]string{"id": "1"},
		outJSON:    `{"message":"OK"}`,
		status:     http.StatusOK,
		authorized: true,
	},
}

func TestMovieHandler_Dislike(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mCaseMock := mock.NewMockMovieUsecase(ctrl)
	handler := &MovieHandler{
		MUCase: mCaseMock,
	}
	for _, test := range movieDislikeTests {
		t.Run(fmt.Sprintf("IN: %v, OUT: %v, CODE: %v", test.inURL, test.outJSON, test.status),
			func(t *testing.T) {
				test.outJSON += "\n"
				r := httptest.NewRequest("POST", test.inURL, nil)
				r = mux.SetURLVars(r, test.inParams)
				w := httptest.NewRecorder()

				if test.authorized {
					sess := &session.Session{UserID: 1}
					err := session.Manager.Create(sess)
					require.NoError(t, err)
					defer sessDelete(session.Manager, sess)
					r.AddCookie(&http.Cookie{
						Name:    "session_id",
						Value:   sess.Cookie,
						Expires: sess.CookieExpiration,
					})
				}

				if test.status == http.StatusOK {
					mCaseMock.EXPECT().Dislike(uint(1), uint(1)).Times(1).Return(nil)
				}

				handler.Dislike(w, r)
				require.Equal(t, test.status, w.Code)
				require.Equal(t, test.outJSON, w.Body.String())
			})
	}
}
