package http

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/movie"
	"Redioteka/internal/pkg/utils/jsonerrors"
	"Redioteka/internal/pkg/utils/log"
	"Redioteka/internal/pkg/utils/session"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

const (
	addFavourite    = "like"
	removeFavourite = "dislike"
)

type MovieHandler struct {
	MUCase domain.MovieUsecase
}

func NewMovieHandlers(router *mux.Router, us domain.MovieUsecase) {
	handler := &MovieHandler{
		MUCase: us,
	}
	router.HandleFunc("/media/movie/{id:[0-9]+}", handler.Get).Methods("GET", "OPTIONS")

	router.HandleFunc("/media/movie/{id:[0-9]+}/like", handler.Like).Methods("POST", "OPTIONS")

	router.HandleFunc("/media/movie/{id:[0-9]+}/dislike", handler.Dislike).Methods("POST", "OPTIONS")

	router.HandleFunc("/media/movie/{id:[0-9]+}/favourites", handler.SetFavourite).Methods("POST", "OPTIONS")

	router.HandleFunc("/media/genres", handler.Genres).Methods("GET", "OPTIONS")

	router.HandleFunc("/media/category/{category}", handler.Category).Methods("GET", "OPTIONS")

}

func (handler *MovieHandler) Genres(w http.ResponseWriter, r *http.Request) {
	genres, err := handler.MUCase.GetGenres()
	if err != nil {
		http.Error(w, jsonerrors.JSONMessage("can't get"), movie.CodeFromError(movie.NotFoundError))
		return
	}

	err = json.NewEncoder(w).Encode(genres)
	if err != nil {
		log.Log.Error(err)
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
		log.Log.Warn(fmt.Sprintf("Error while getting user id: %s", err))
		http.Error(w, jsonerrors.URLParams, http.StatusBadRequest)
		return
	}
	id := uint(id64)

	sess, _ := session.GetSession(r)
	foundMovie, err := handler.MUCase.GetByID(id, sess)
	if err != nil {
		http.Error(w, jsonerrors.JSONMessage("not found"), movie.CodeFromError(err))
		return
	}

	err = json.NewEncoder(w).Encode(foundMovie)
	if err != nil {
		log.Log.Error(err)
		http.Error(w, jsonerrors.JSONEncode, http.StatusInternalServerError)
		return
	}
}

func parseOrder(catName string) int {
	switch catName {
	case "top":
		return domain.RatingOrder
	case "newest":
		return domain.DateOrder
	default:
		return domain.NoneOrder
	}
}

func (handler *MovieHandler) Category(w http.ResponseWriter, r *http.Request) {
	catName, found := mux.Vars(r)["category"]
	if !found {
		log.Log.Warn("Can't parse category")
		http.Error(w, jsonerrors.URLParams, http.StatusBadRequest)
		return
	}

	decoder := schema.NewDecoder()
	filter := domain.MovieFilter{}
	err := decoder.Decode(&filter, r.URL.Query())
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Error while parsing querystring %s", err))
		http.Error(w, jsonerrors.JSONMessage("getting movie array"), http.StatusBadRequest)
		return
	}
	filter.Order = parseOrder(catName)

	foundMovies, err := handler.MUCase.GetByFilter(filter)
	if err != nil {
		http.Error(w, jsonerrors.JSONMessage("getting movie array"), movie.CodeFromError(err))
		return
	}

	err = json.NewEncoder(w).Encode(map[string][]domain.Movie{catName: foundMovies})
	if err != nil {
		log.Log.Error(err)
		http.Error(w, jsonerrors.JSONEncode, http.StatusInternalServerError)
		return
	}
}
