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
		{"{}", `{"id":1,"email":"","username":""}`, http.StatusOK},
		{`{"username":"good_user","email":"gmail@mail.ru","password":"pass","repeated_password":"pass"}`,
			`{"id":2,"email":"gmail@mail.ru","username":"good_user"}`, http.StatusOK},
		{`{"username":"good_user","email":"gmail@mail.ru","password":"pass","repeated_password":"pass"}`,
			`{"error":"Wrong username or password"}`, http.StatusBadRequest},
	}

	api := &Handler{}
	for i, test := range tests {
		fmt.Println("Test", i)
		body := bytes.NewReader([]byte(test.inJSON))
		r := httptest.NewRequest("POST", "/users/signup", body)
		w := httptest.NewRecorder()
		api.Signup(w, r)
		require.Equal(t, test.status, w.Code)
		require.Equal(t, test.outJSON+"\n", w.Body.String())
	}
}
