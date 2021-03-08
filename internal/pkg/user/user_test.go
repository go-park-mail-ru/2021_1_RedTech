package user

import (
	"bytes"
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
		{`{"username":"a","email":"ya@mail.ru","password":"pass","repeated_password":"not_pass"}`,
			`{"error":"Passwords do not match"}`, http.StatusBadRequest},
		{"{}", `{"error":"Empty fields in form"}`, http.StatusBadRequest},
		{`{"username":"good_user","email":"gmail@mail.ru","password":"pass","repeated_password":"pass"}`,
			`{"id":1,"email":"gmail@mail.ru","username":"good_user"}`, http.StatusOK},
		{`{"username":"good_user","email":"gmail@mail.ru","password":"pass","repeated_password":"pass"}`,
			`{"error":"Wrong username or password"}`, http.StatusBadRequest},
	}

	api := &Handler{}
	for i, test := range tests {
		fmt.Println("TestSignup", i)
		body := bytes.NewReader([]byte(test.inJSON))
		r := httptest.NewRequest("POST", "/users/signup", body)
		w := httptest.NewRecorder()
		api.Signup(w, r)
		require.Equal(t, test.status, w.Code)
		require.Equal(t, test.outJSON+"\n", w.Body.String())
	}
}

func TestHandlerLogin(t *testing.T) {
	testForm := []TestCase{
		{"", `{"error":"bad form"}`, http.StatusBadRequest},
		{`{"bad_field":bad_value}`, `{"error":"bad form"}`, http.StatusBadRequest},
		{"{}", `{"error":"Empty login or password"}`, http.StatusBadRequest},
	}

	api := &Handler{}
	for i, test := range testForm {
		fmt.Println("TestLogin (form parsing)", i)
		body := bytes.NewReader([]byte(test.inJSON))
		r := httptest.NewRequest("POST", "/users/login", body)
		w := httptest.NewRecorder()
		api.Login(w, r)
		require.Equal(t, test.status, w.Code)
		require.Equal(t, test.outJSON+"\n", w.Body.String())
	}

	testData := []TestCase{
		{`{"email":"not_exist@mail.ru","password":"pass"}`,
			`{"error":"Wrong login or password"}`, http.StatusBadRequest},
		{`{"email":"gmail@mail.ru","password":"not_pass"}`,
			`{"error":"Wrong login or password"}`, http.StatusBadRequest},
		{`{"email":"gmail@mail.ru","password":"pass"}`,
			`{"id":1,"email":"gmail@mail.ru","username":"good_user"}`, http.StatusOK},
	}
	testUser := `{"username":"good_user","email":"gmail@mail.ru","password":"pass","repeated_password":"pass"}`
	api.Signup(httptest.NewRecorder(),
		httptest.NewRequest("POST", "/users/signup", bytes.NewReader([]byte(testUser))))

	for i, test := range testData {
		fmt.Println("TestLogin (data checking)", i)
		body := bytes.NewReader([]byte(test.inJSON))
		r := httptest.NewRequest("POST", "/users/login", body)
		w := httptest.NewRecorder()
		api.Login(w, r)
		require.Equal(t, test.status, w.Code)
		require.Equal(t, test.outJSON+"\n", w.Body.String())
	}
}
