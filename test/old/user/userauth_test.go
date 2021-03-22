package user

import (
	"Redioteka/internal/app/session"
	"Redioteka/internal/app/user"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestCase struct {
	inJSON  string
	outJSON string
	status  int
}

func TestHandlerSignup(t *testing.T) {
	tests := []TestCase{
		{"", `{"error":"bad form"}`, http.StatusBadRequest},
		{`{"bad_field":bad_value}`, `{"error":"bad form"}`, http.StatusBadRequest},
		{`{"username":"a","email":"ya@mail.ru","password":"pass","confirm_password":"not_pass"}`,
			`{"error":"Passwords do not match"}`, http.StatusBadRequest},
		{"{}", `{"error":"Empty fields in form"}`, http.StatusBadRequest},
		{`{"username":"good_user","email":"gmail@mail.ru","password":"pass","confirm_password":"pass"}`,
			`{"id":1,"email":"gmail@mail.ru","username":"good_user"}`, http.StatusOK},
		{`{"username":"good_user","email":"gmail@mail.ru","password":"pass","confirm_password":"pass"}`,
			`{"error":"Wrong username or password"}`, http.StatusBadRequest},
	}

	api := &user.Handler{}
	for i, test := range tests {
		test.outJSON += "\n"
		fmt.Println("TestSignup", i)
		body := bytes.NewReader([]byte(test.inJSON))
		r := httptest.NewRequest("POST", "/users/signup", body)
		w := httptest.NewRecorder()
		api.Signup(w, r)
		actual := TestCase{
			inJSON:  test.inJSON,
			outJSON: w.Body.String(),
			status:  w.Code,
		}
		require.Equal(t, test, actual)
	}
}

func TestHandlerLogin(t *testing.T) {
	testForm := []TestCase{
		{"", `{"error":"bad form"}`, http.StatusBadRequest},
		{`{"bad_field":bad_value}`, `{"error":"bad form"}`, http.StatusBadRequest},
		{"{}", `{"error":"Empty login or password"}`, http.StatusBadRequest},
	}

	api := &user.Handler{}
	for i, test := range testForm {
		test.outJSON += "\n"
		fmt.Println("TestLogin (form parsing)", i)
		body := bytes.NewReader([]byte(test.inJSON))
		r := httptest.NewRequest("POST", "/users/login", body)
		w := httptest.NewRecorder()
		api.Login(w, r)
		actual := TestCase{
			inJSON:  test.inJSON,
			outJSON: w.Body.String(),
			status:  w.Code,
		}
		require.Equal(t, test, actual)
	}

	testData := []TestCase{
		{`{"email":"not_exist@mail.ru","password":"pass"}`,
			`{"error":"Wrong login or password"}`, http.StatusBadRequest},
		{`{"email":"gmail@mail.ru","password":"not_pass"}`,
			`{"error":"Wrong login or password"}`, http.StatusBadRequest},
		{`{"email":"gmail@mail.ru","password":"pass"}`,
			`{"id":1,"email":"gmail@mail.ru","username":"good_user"}`, http.StatusOK},
	}
	testUser := `{"username":"good_user","email":"gmail@mail.ru","password":"pass","confirm_password":"pass"}`
	api.Signup(httptest.NewRecorder(),
		httptest.NewRequest("POST", "/users/signup", bytes.NewReader([]byte(testUser))))

	for i, test := range testData {
		test.outJSON += "\n"
		fmt.Println("TestLogin (data checking)", i)
		body := bytes.NewReader([]byte(test.inJSON))
		r := httptest.NewRequest("POST", "/users/login", body)
		w := httptest.NewRecorder()
		api.Login(w, r)
		actual := TestCase{
			inJSON:  test.inJSON,
			outJSON: w.Body.String(),
			status:  w.Code,
		}
		require.Equal(t, test, actual)
	}
}

func TestHandlerLogout(t *testing.T) {
	userForm := &user.userSignupForm{
		Login:         "user",
		Email:         "gmail.mail.ru",
		Password:      "pass",
		PasswordCheck: "pass",
	}
	api := &user.Handler{}

	r := httptest.NewRequest("GET", "/users/logout", nil)
	w := httptest.NewRecorder()
	api.Logout(w, r)
	require.Equal(t, `{"error":"user not found"}`+"\n", w.Body.String())
	require.Equal(t, http.StatusBadRequest, w.Code)

	form, _ := json.Marshal(userForm)
	body := bytes.NewBuffer(form)
	r = httptest.NewRequest("POST", "/users/signup", body)
	w = httptest.NewRecorder()
	api.Signup(w, r)
	require.Equal(t, http.StatusOK, w.Code)

	cookies := w.Result().Cookies()
	r = httptest.NewRequest("GET", "/users/logout", body)
	for _, cookie := range cookies {
		r.AddCookie(cookie)
	}
	w = httptest.NewRecorder()
	api.Logout(w, r)
	require.Equal(t, `{"status":"OK"}`, w.Body.String())
	require.Equal(t, http.StatusOK, w.Code)

	userID, err := session.Check(r)
	require.Equal(t, uint(0), userID)
	require.Equal(t, nil, err)
}
