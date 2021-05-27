package http

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/utils/randstring"
	"Redioteka/internal/pkg/utils/session"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type ResponseError struct {
	Message string `json:"message"`
}

type UserHandler struct {
	UUsecase       domain.UserUsecase
	SessionManager session.SessionManager
}

func NewUserHandlers(router *mux.Router, uc domain.UserUsecase, sm session.SessionManager) {
	handler := &UserHandler{
		UUsecase:       uc,
		SessionManager: sm,
	}
	router.HandleFunc("/users/signup", handler.Signup).Methods("POST", "OPTIONS")

	router.HandleFunc("/users/login", handler.Login).Methods("POST", "OPTIONS")

	router.HandleFunc("/users/logout", handler.Logout).Methods("GET", "OPTIONS")

	router.HandleFunc("/users/{id:[0-9]+}", handler.Get).Methods("GET", "OPTIONS")

	router.HandleFunc("/me", handler.Me).Methods("GET", "OPTIONS")

	router.HandleFunc("/users/{id:[0-9]+}", handler.Update).Methods("PATCH", "OPTIONS")

	router.HandleFunc("/users/{id:[0-9]+}/avatar", handler.Avatar).Methods("POST", "PUT", "OPTIONS")

	router.HandleFunc("/users/{id:[0-9]+}/media", handler.GetMedia).Methods("GET", "OPTIONS")

	router.HandleFunc("/csrf", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     "csrf_token",
			Value:    randstring.RandString(32),
			Path:     "/",
			Expires:  time.Now().Add(900 * time.Second),
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})
		w.WriteHeader(http.StatusNoContent)
	}).Methods("GET", "OPTIONS")
}
