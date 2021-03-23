package http

import (
	"Redioteka/internal/pkg/domain"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type ResponseError struct {
	Message string `json:"message"`
}

type UserHandler struct {
	MHanlder domain.MovieUsecase
}

func NewMovieHandlers(router *mux.Router, us domain.MovieUsecase) {
	handler := &UserHandler{
		MHanlder: us,
	}
	router.HandleFunc("/media/movie/{id:[0-9]+}", handler.Get).Methods("GET", "OPTIONS")
}

func (handler *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Первый аргумент в парсинге беззнаковых чисел - база системы счисления, второй -
	// количество бит, которые он занимает. Четырех миллиардов пользователей нам хватит
	id64, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		log.Printf("Error while getting user: %s", err)
		http.Error(w, `{"error":"bad id"}`, http.StatusBadRequest)
		return
	}
	id := uint(id64)

	movie, err := handler.MHanlder.GetById(id)
	if err != nil {
		log.Printf("This movie does not exist")
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(movie)
	if err != nil {
		log.Printf("error while encoding JSON: %s", err)
		http.Error(w, `{"error":"server"}`, http.StatusInternalServerError)
		return
	}
}