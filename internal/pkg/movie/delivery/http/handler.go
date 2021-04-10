package http

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/movie"
	"Redioteka/internal/pkg/user"
	"Redioteka/internal/pkg/utils/jsonerrors"
	"Redioteka/internal/pkg/utils/session"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

	router.HandleFunc("/media/movie/{id:[0-9]+}/stream", handler.Stream).Methods("GET", "OPTIONS")
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
func (handler *MovieHandler) Stream(w http.ResponseWriter, r *http.Request) {
	sess, err := session.GetSession(r)
	if err != nil || session.Manager.Check(sess) != nil {
		log.Printf("Trying to get stream while unauthorized")
		http.Error(w, jsonerrors.JSONMessage("unauthorized"), user.CodeFromError(user.UnauthorizedError))
		return
	}

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

	foundStream, err := handler.MUCase.GetStream(id)
	if err != nil {
		log.Printf("This movie does not exist")
		http.Error(w, jsonerrors.JSONMessage("not found"), movie.CodeFromError(err))
		return
	}

	err = json.NewEncoder(w).Encode(foundStream)
	if err != nil {
		log.Printf("Error while encoding JSON: %s", err)
		http.Error(w, jsonerrors.JSONEncode, http.StatusInternalServerError)
		return
	}
}
