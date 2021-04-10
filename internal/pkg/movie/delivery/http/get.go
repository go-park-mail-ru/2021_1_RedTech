package http

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/movie"
	"Redioteka/internal/pkg/utils/jsonerrors"
	"Redioteka/internal/pkg/utils/log"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"net/http"
	"strconv"
)

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

	foundMovies, err := handler.MUCase.GetByFilter(filter)
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Error while getting movie array: %s", err))
		http.Error(w, jsonerrors.JSONMessage("getting movie array"), movie.CodeFromError(err))
		return
	}

	err = json.NewEncoder(w).Encode(map[string][]domain.Movie{catName: foundMovies})
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Error while encoding JSON: %s", err))
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
		log.Log.Warn(fmt.Sprintf("Error while getting user: %s", err))
		http.Error(w, jsonerrors.URLParams, http.StatusBadRequest)
		return
	}
	id := uint(id64)

	foundMovie, err := handler.MUCase.GetById(id)
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Movie with id %v doesn't exist", id))
		http.Error(w, jsonerrors.JSONMessage("not found"), movie.CodeFromError(err))
		return
	}

	err = json.NewEncoder(w).Encode(foundMovie)
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Error while encoding JSON: %s", err))
		http.Error(w, jsonerrors.JSONEncode, http.StatusInternalServerError)
		return
	}
}
