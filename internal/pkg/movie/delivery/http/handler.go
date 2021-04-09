package http

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/movie"
	"Redioteka/internal/pkg/utils/jsonerrors"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"log"
	"net/http"
	"strconv"
)

type MovieHandler struct {
	MUCase domain.MovieUsecase
}

func NewMovieHandlers(router *mux.Router, us domain.MovieUsecase) {
	handler := &MovieHandler{
		MUCase: us,
	}
	router.HandleFunc("/media/movie/{id:[0-9]+}", handler.Get).Methods("GET", "OPTIONS")
	router.HandleFunc("/media/genres", handler.Genres).Methods("GET", "OPTIONS")
	router.HandleFunc("/media/category/{category}", handler.Category).Methods("GET", "OPTIONS")
}

func (handler *MovieHandler) Genres(w http.ResponseWriter, r *http.Request) {
	genres, err := handler.MUCase.GetGenres()
	if err != nil {
		log.Println("Can't get genres")
		http.Error(w, jsonerrors.JSONMessage("can't get"), movie.CodeFromError(movie.NotFoundError))
		return
	}

	err = json.NewEncoder(w).Encode(map[string][]string{"genres": genres})
	if err != nil {
		log.Printf("Error while encoding JSON: %s", err)
		http.Error(w, jsonerrors.JSONEncode, http.StatusInternalServerError)
		return
	}
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
		log.Printf("Error while encoding JSON: %s", err)
		http.Error(w, jsonerrors.JSONEncode, http.StatusInternalServerError)
		return
	}
}

func (handler *MovieHandler) Category(w http.ResponseWriter, r *http.Request) {
	catName, found := mux.Vars(r)["category"]
	if !found {
		log.Printf("Can't parse category")
		http.Error(w, jsonerrors.URLParams, http.StatusBadRequest)
		return
	}

	decoder := schema.NewDecoder()
	filter := domain.MovieFilter{}
	err := decoder.Decode(&filter, r.URL.Query())
	if err != nil {
		log.Printf("Error while parsing querystring %s", err)
		http.Error(w, jsonerrors.JSONMessage("getting movie array"), http.StatusBadRequest)
		return
	}

	foundMovies, err := handler.MUCase.GetByFilter(filter)
	if err != nil {
		log.Printf("Error while getting movie array: %s", err)
		http.Error(w, jsonerrors.JSONMessage("getting movie array"), movie.CodeFromError(err))
		return
	}

	err = json.NewEncoder(w).Encode(map[string][]domain.Movie{catName: foundMovies})
	if err != nil {
		log.Printf("Error while encoding JSON: %s", err)
		http.Error(w, jsonerrors.JSONEncode, http.StatusInternalServerError)
		return
	}
}
