package http

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/movie"
	"Redioteka/internal/pkg/utils/jsonerrors"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const (
	addFavourite    = "add"
	removeFavourite = "remove"
)

type MovieHandler struct {
	MUCase domain.MovieUsecase
}

func NewMovieHandlers(router *mux.Router, us domain.MovieUsecase) {
	handler := &MovieHandler{
		MUCase: us,
	}
	router.HandleFunc("/media/movie/{id:[0-9]+}", handler.Get).Methods("GET", "OPTIONS")

	router.HandleFunc("/media/movie/{id:[0-9]+}/like", handler.SetFavourite).Methods("POST", "OPTIONS").Name(addFavourite)

	router.HandleFunc("/media/movie/{id:[0-9]+}/dislike", handler.SetFavourite).Methods("POST", "OPTIONS").Name(removeFavourite)
}

func (handler *MovieHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Первый аргумент в парсинге беззнаковых чисел - база системы счисления, второй -
	// количество бит, которые он занимает. Четырех миллиардов пользователей нам хватит
	id64, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		log.Printf("Error while getting user: %s", err)
		http.Error(w, jsonerrors.URLParams, http.StatusBadRequest)
		return
	}
	id := uint(id64)

	foundMovie, err := handler.MUCase.GetById(id)
	if err != nil {
		log.Printf("This movie does not exist")
		http.Error(w, jsonerrors.JSONMessage("not found"), movie.CodeFromError(err))
		return
	}

	err = json.NewEncoder(w).Encode(foundMovie)
	if err != nil {
		log.Printf("error while encoding JSON: %s", err)
		http.Error(w, jsonerrors.JSONEncode, http.StatusInternalServerError)
		return
	}
}
