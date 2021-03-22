package http

import (
	"Redioteka/internal/pkg/domain"
	"github.com/gorilla/mux"
)

type ResponseError struct {
	Message string `json:"message"`
}

type UserHandler struct {
	UHandler domain.UserUsecase
}

func NewUserHandlers(router *mux.Router, us domain.UserUsecase) {
	handler := &UserHandler{
		UHandler: us,
	}
	router.HandleFunc("/users/signup", handler.Signup).Methods("POST", "OPTIONS")

	router.HandleFunc("/users/login", handler.Login).Methods("POST", "OPTIONS")

	router.HandleFunc("/users/logout", handler.Logout).Methods("GET", "OPTIONS")

	router.HandleFunc("/users/{id:[0-9]+}", handler.Get).Methods("GET", "OPTIONS")

	router.HandleFunc("/me", handler.Me).Methods("GET", "OPTIONS")

	router.HandleFunc("/users/{id:[0-9]+}", handler.Update).Methods("PATCH", "OPTIONS")

	router.HandleFunc("/users/{id:[0-9]+}/avatar", handler.Avatar).Methods("POST", "PUT", "OPTIONS")
}
