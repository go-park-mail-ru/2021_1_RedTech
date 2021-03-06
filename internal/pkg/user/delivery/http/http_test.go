package http

/*
import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/user"
	mock2 "Redioteka/internal/pkg/user/usecase/mock"
	"Redioteka/internal/pkg/utils/jsonerrors"
	"Redioteka/internal/pkg/utils/log"
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

var usersTestData = map[uint]domain.User{
	1: {
		ID:       1,
		Email:    "mail1@mail.ru",
		Username: "user1",
	},
	2: {
		ID:       2,
		Email:    "mail2@mail.ru",
		Username: "user2",
	},
}

var userTestErrors = map[uint]error{
	3: user.InvalidCredentials,
}

type TestCaseGet struct {
	ID      uint
	outJSON string
	status  int
}

var testCaseGet = []TestCaseGet{
	{
		ID:      1,
		outJSON: `{"id":1,"email":"mail1@mail.ru","username":"user1"}`,
		status:  http.StatusOK,
	},
	{
		ID:      2,
		outJSON: `{"id":2,"username":"user2"}`,
		status:  http.StatusOK,
	},
	{
		ID:      3,
		outJSON: jsonerrors.JSONMessage("get"),
		status:  http.StatusForbidden,
	},
}

func sessDelete(m session.SessionManager, s *session.Session) {
	err := m.Delete(s)
	if err != nil {
		log.Log.Error(err)
	}
}

func TestUserHandler_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uCaseMock := mock2.NewMockUserUsecase(ctrl)
	for uid, value := range usersTestData {
		uCaseMock.EXPECT().GetById(uid).Times(1).Return(value, nil)
	}
	for uid, errValue := range userTestErrors {
		uCaseMock.EXPECT().GetById(uid).Times(1).Return(domain.User{}, errValue)
	}
	handler := &UserHandler{
		UUsecase: uCaseMock,
	}
	for _, test := range testCaseGet {
		t.Run(fmt.Sprintf("IN: %v, OUT: %v, CODE: %v", test.ID, test.outJSON, test.status),
			func(t *testing.T) {
				test.outJSON += "\n"
				r := httptest.NewRequest("GET", "/api/users/"+strconv.Itoa(int(test.ID)), nil)
				r = mux.SetURLVars(r, map[string]string{"id": strconv.Itoa(int(test.ID))})
				w := httptest.NewRecorder()

				sess := &session.Session{UserID: 1}
				err := session.Manager.Create(sess)
				require.NoError(t, err)
				r.AddCookie(&http.Cookie{
					Name:    "session_id",
					Value:   sess.Cookie,
					Expires: sess.CookieExpiration,
				})
				defer sessDelete(session.Manager, sess)

				handler.Get(w, r)
				current := TestCaseGet{
					ID:      test.ID,
					outJSON: w.Body.String(),
					status:  w.Code,
				}
				require.Equal(t, test, current)
			})
	}
}

var meTests = []TestCaseGet{
	{
		// 1 is authorized
		ID:      1,
		outJSON: `{"id":1}`,
		status:  http.StatusOK,
	},
	{
		// 2 is unauthorized
		ID:      2,
		outJSON: jsonerrors.JSONMessage("unauthorized"),
		status:  http.StatusUnauthorized,
	},
}

func TestUserHandler_Me(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uCaseMock := mock2.NewMockUserUsecase(ctrl)
	handler := &UserHandler{
		UUsecase: uCaseMock,
	}
	for _, test := range meTests {
		t.Run(fmt.Sprintf("IN: %v, OUT: %v, CODE: %v", test.ID, test.outJSON, test.status),
			func(t *testing.T) {
				test.outJSON += "\n"
				r := httptest.NewRequest("GET", "/api/users/"+strconv.Itoa(int(test.ID)), nil)
				r = mux.SetURLVars(r, map[string]string{"id": strconv.Itoa(int(test.ID))})
				w := httptest.NewRecorder()

				if test.ID == 1 {
					sess := &session.Session{UserID: test.ID}
					err := session.Manager.Create(sess)
					require.NoError(t, err)
					r.AddCookie(&http.Cookie{
						Name:    "session_id",
						Value:   sess.Cookie,
						Expires: sess.CookieExpiration,
					})
					defer sessDelete(session.Manager, sess)
				}

				handler.Me(w, r)
				current := TestCaseGet{
					ID:      test.ID,
					outJSON: w.Body.String(),
					status:  w.Code,
				}
				require.Equal(t, test, current)
			})
	}
}

type userTestCase struct {
	inJSON  string
	inUser  domain.User
	outUser domain.User
	outJSON string
	status  int
}

var updateTests = []userTestCase{
	{
		inJSON: `{"email":"newmail1@mail.ru"}` + "\n",
		inUser: domain.User{
			ID:    1,
			Email: "newmail1@mail.ru",
		},
		outUser: domain.User{
			ID:       1,
			Email:    "newmail1@mail.ru",
			Username: "user1",
		},
		outJSON: `{"id":1,"email":"newmail1@mail.ru","username":"user1"}`,
		status:  http.StatusOK,
	},
	{
		inJSON: `{"username":"new_user1"}` + "\n",
		inUser: domain.User{
			ID:       1,
			Username: "new_user1",
		},
		outUser: domain.User{
			ID:       1,
			Username: "new_user1",
			Email:    "mail1@mail.ru",
		},
		outJSON: `{"id":1,"email":"mail1@mail.ru","username":"new_user1"}`,
		status:  http.StatusOK,
	},
	{
		inJSON: `{"email":"newmail1@mail.ru","username":"new_user1"}` + "\n",
		inUser: domain.User{
			ID:       1,
			Email:    "newmail1@mail.ru",
			Username: "new_user1",
		},
		outUser: domain.User{
			ID:       1,
			Email:    "newmail1@mail.ru",
			Username: "new_user1",
		},
		outJSON: `{"id":1,"email":"newmail1@mail.ru","username":"new_user1"}`,
		status:  http.StatusOK,
	},
	{
		inJSON: `{"email":"","username":""}` + "\n",
		inUser: domain.User{
			ID:       1,
			Email:    "",
			Username: "",
		},
		outUser: domain.User{},
		outJSON: `{"message":"invalid update"}`,
		status:  http.StatusNotAcceptable,
	},
	{
		inJSON: `{}` + "\n",
		inUser: domain.User{
			ID: 1,
		},
		outUser: domain.User{},
		outJSON: `{"message":"invalid update"}`,
		status:  http.StatusNotAcceptable,
	},
}

func TestUserHandler_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uCaseMock := mock2.NewMockUserUsecase(ctrl)
	handler := &UserHandler{
		UUsecase: uCaseMock,
	}
	for _, test := range updateTests {
		t.Run(fmt.Sprintf("IN: %v, OUT: %v, CODE: %v", test.inJSON, test.outJSON, test.status),
			func(t *testing.T) {
				test.outJSON += "\n"
				body := bytes.NewReader([]byte(test.inJSON))
				r := httptest.NewRequest("PATCH", "/api/users/1", body)
				r = mux.SetURLVars(r, map[string]string{"id": "1"})
				w := httptest.NewRecorder()

				sess := &session.Session{UserID: test.inUser.ID}
				err := session.Manager.Create(sess)
				require.NoError(t, err)
				r.AddCookie(&http.Cookie{
					Name:    "session_id",
					Value:   sess.Cookie,
					Expires: sess.CookieExpiration,
				})
				defer sessDelete(session.Manager, sess)

				if test.status == http.StatusOK {
					uCaseMock.EXPECT().Update(&test.inUser).Times(1).Return(nil)
					uCaseMock.EXPECT().GetById(uint(1)).Times(1).Return(test.outUser, nil)
				} else {
					uCaseMock.EXPECT().Update(&test.inUser).Times(1).Return(user.InvalidUpdateError)
				}
				handler.Update(w, r)
				current := userTestCase{
					inJSON:  test.inJSON,
					inUser:  test.inUser,
					outUser: test.outUser,
					outJSON: w.Body.String(),
					status:  w.Code,
				}
				require.Equal(t, test, current)
			})
	}
}

var signupTests = []userTestCase{
	{
		inJSON:  "",
		inUser:  domain.User{},
		outUser: domain.User{},
		outJSON: `{"message":"json decode"}`,
		status:  http.StatusBadRequest,
	},
	{
		inJSON:  `{"bad_field":"bad_value"}`,
		inUser:  domain.User{},
		outUser: domain.User{},
		outJSON: `{"message":"signup"}`,
		status:  http.StatusForbidden,
	},
	{
		inJSON: `{"username":"a","email":"ya@mail.ru","password":"pass","confirm_password":"not_pass"}`,
		inUser: domain.User{
			Username:             "a",
			Email:                "ya@mail.ru",
			InputPassword:        "pass",
			ConfirmInputPassword: "not_pass",
		},
		outUser: domain.User{},
		outJSON: `{"message":"signup"}`,
		status:  http.StatusForbidden,
	},
	{
		inJSON:  "{}",
		inUser:  domain.User{},
		outUser: domain.User{},
		outJSON: `{"message":"signup"}`,
		status:  http.StatusForbidden,
	},
	{
		inJSON: `{"username":"good_user","email":"gmail@mail.ru","password":"pass","confirm_password":"pass"}`,
		inUser: domain.User{
			Email:                "gmail@mail.ru",
			Username:             "good_user",
			InputPassword:        "pass",
			ConfirmInputPassword: "pass",
		},
		outUser: domain.User{
			ID:       1,
			Email:    "gmail@mail.ru",
			Username: "good_user",
		},
		outJSON: `{"id":1,"email":"gmail@mail.ru","username":"good_user"}`,
		status:  http.StatusOK,
	},
}

func TestUserHandler_Signup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uCaseMock := mock2.NewMockUserUsecase(ctrl)
	handler := &UserHandler{
		UUsecase: uCaseMock,
	}
	for _, test := range signupTests {
		t.Run(fmt.Sprintf("IN: %v, OUT: %v, CODE: %v", test.inJSON, test.outJSON, test.status),
			func(t *testing.T) {
				test.outJSON += "\n"
				body := bytes.NewReader([]byte(test.inJSON))
				r := httptest.NewRequest("POST", "/api/users/signup", body)
				w := httptest.NewRecorder()
				if test.status == http.StatusOK {
					uCaseMock.EXPECT().Signup(&test.inUser).Times(1).Return(test.outUser, &session.Session{}, nil)
				} else if test.status != http.StatusBadRequest {
					uCaseMock.EXPECT().Signup(&test.inUser).Times(1).Return(domain.User{}, nil, user.InvalidCredentials)
				}
				handler.Signup(w, r)
				current := userTestCase{
					inJSON:  test.inJSON,
					inUser:  test.inUser,
					outUser: test.outUser,
					outJSON: w.Body.String(),
					status:  w.Code,
				}
				require.Equal(t, current, test)
			})
	}
}

var loginTests = []userTestCase{
	{
		inJSON:  "",
		outJSON: `{"message":"json decode"}`,
		status:  http.StatusBadRequest,
	},
	{
		inJSON:  `{"bad_field":"bad_value"}`,
		inUser:  domain.User{},
		outJSON: `{"message":"login"}`,
		status:  http.StatusForbidden,
	},
	{
		inJSON:  "{}",
		outJSON: `{"message":"login"}`,
		status:  http.StatusForbidden,
	},
}

func TestUserHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uCaseMock := mock2.NewMockUserUsecase(ctrl)
	handler := &UserHandler{
		UUsecase: uCaseMock,
	}
	for _, test := range loginTests {
		t.Run(fmt.Sprintf("IN: %v, OUT: %v, CODE: %v", test.inJSON, test.outJSON, test.status),
			func(t *testing.T) {
				test.outJSON += "\n"
				body := bytes.NewReader([]byte(test.inJSON))
				r := httptest.NewRequest("POST", "/api/users/login", body)
				w := httptest.NewRecorder()
				if test.status == http.StatusOK {
					uCaseMock.EXPECT().Login(&test.inUser).Times(1).Return(test.outUser, &session.Session{}, nil)
				} else if test.status != http.StatusBadRequest {
					uCaseMock.EXPECT().Login(&test.inUser).Times(1).Return(domain.User{}, nil, user.InvalidCredentials)
				}
				handler.Login(w, r)
				current := userTestCase{
					inJSON:  test.inJSON,
					inUser:  test.inUser,
					outUser: test.outUser,
					outJSON: w.Body.String(),
					status:  w.Code,
				}
				require.Equal(t, current, test)
			})
	}
}

type logoutTestCase struct {
	sess   *session.Session
	status int
}

var logoutTests = []logoutTestCase{
	{
		sess:   &session.Session{},
		status: http.StatusBadRequest,
	},
	{
		sess:   &session.Session{UserID: 1},
		status: http.StatusInternalServerError,
	},
	{
		sess:   &session.Session{UserID: 2},
		status: http.StatusOK,
	},
}

func TestUserHandler_Logout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uCaseMock := mock2.NewMockUserUsecase(ctrl)
	handler := &UserHandler{
		UUsecase: uCaseMock,
	}
	for _, test := range logoutTests {
		t.Run(fmt.Sprintf("CODE: %v", test.status),
			func(t *testing.T) {
				r := httptest.NewRequest("GET", "/api/users/logout", nil)
				w := httptest.NewRecorder()
				if test.status != http.StatusBadRequest {
					err := session.Manager.Create(test.sess)
					require.NoError(t, err)
					r.AddCookie(&http.Cookie{
						Name:    "session_id",
						Value:   test.sess.Cookie,
						Expires: test.sess.CookieExpiration,
					})
					test.sess = &session.Session{Cookie: test.sess.Cookie}
					defer sessDelete(session.Manager, test.sess)
				}
				if test.status == http.StatusOK {
					uCaseMock.EXPECT().Logout(test.sess).Times(1).Return(test.sess, nil)
				} else if test.status != http.StatusBadRequest {
					uCaseMock.EXPECT().Logout(test.sess).Times(1).Return(test.sess, user.UnauthorizedError)
				}
				handler.Logout(w, r)
				require.Equal(t, test.status, w.Code)
			})
	}
}

type avatarTestCase struct {
	inURL       string
	inURLParams map[string]string
	outJSON     string
	status      int
}

var avatarTests = []avatarTestCase{
	{
		inURL:       "/api/users//avatar",
		inURLParams: map[string]string{},
		outJSON:     `{"message":"url params"}`,
		status:      http.StatusBadRequest,
	},
	{
		inURL:       "/api/users/2/avatar",
		inURLParams: map[string]string{"id": "2"},
		outJSON:     `{"message":"wrong id"}`,
		status:      http.StatusForbidden,
	},
}

func TestUserHandler_Avatar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uCaseMock := mock2.NewMockUserUsecase(ctrl)
	handler := &UserHandler{
		UUsecase: uCaseMock,
	}
	for _, test := range avatarTests {
		t.Run(fmt.Sprintf("IN: %v, OUT: %v, CODE: %v", test.inURL, test.outJSON, test.status),
			func(t *testing.T) {
				test.outJSON += "\n"
				body := bytes.NewReader([]byte(""))
				r := httptest.NewRequest("POST", test.inURL, body)
				r = mux.SetURLVars(r, test.inURLParams)
				w := httptest.NewRecorder()

				sess := &session.Session{UserID: 1}
				err := session.Manager.Create(sess)
				require.NoError(t, err)
				r.AddCookie(&http.Cookie{
					Name:    "session_id",
					Value:   sess.Cookie,
					Expires: sess.CookieExpiration,
				})
				defer sessDelete(session.Manager, sess)

				handler.Avatar(w, r)
				current := avatarTestCase{
					inURL:       test.inURL,
					inURLParams: test.inURLParams,
					outJSON:     w.Body.String(),
					status:      w.Code,
				}
				require.Equal(t, current, test)
			})
	}
}

type getMediaTestCase struct {
	inURL       string
	inURLParams map[string]string
	outJSON     string
	status      int
}

var getMediaTests = []getMediaTestCase{
	{
		inURL:       "/api/users//media",
		inURLParams: map[string]string{},
		outJSON:     jsonerrors.URLParams,
		status:      http.StatusBadRequest,
	},
	{
		inURL:       "/api/users/2/media",
		inURLParams: map[string]string{"id": "2"},
		outJSON:     jsonerrors.Session,
		status:      http.StatusUnauthorized,
	},
	{
		inURL:       "/api/users/3/media",
		inURLParams: map[string]string{"id": "3"},
		outJSON:     `{"favourites":[{"id":1,"rating":9,"title":"Film","is_free":false,"movie_avatar":"/static/movies/default.jpg"}]}`,
		status:      http.StatusOK,
	},
}

func TestUserHandler_GetMedia(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	uCaseMock := mock2.NewMockUserUsecase(ctrl)
	handler := &UserHandler{
		UUsecase: uCaseMock,
	}
	for _, test := range getMediaTests {
		t.Run(fmt.Sprintf("IN: %v, OUT: %v, CODE: %v", test.inURL, test.outJSON, test.status),
			func(t *testing.T) {
				test.outJSON += "\n"
				r := httptest.NewRequest("GET", test.inURL, nil)
				r = mux.SetURLVars(r, test.inURLParams)
				w := httptest.NewRecorder()

				if test.status == http.StatusOK {
					id, err := strconv.Atoi(test.inURLParams["id"])
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
					defer sessDelete(session.Manager, sess)

					testMovies := []domain.Movie{
						{
							ID:     1,
							Title:  "Film",
							Rating: 9,
							IsFree: false,
							Avatar: "/static/movies/default.jpg",
						},
					}
					uCaseMock.EXPECT().GetFavourites(uint(id), sess).Times(1).Return(testMovies, nil)
				}

				handler.GetMedia(w, r)
				current := getMediaTestCase{
					inURL:       test.inURL,
					inURLParams: test.inURLParams,
					outJSON:     w.Body.String(),
					status:      w.Code,
				}
				require.Equal(t, test, current)
			})
	}
}
*/
