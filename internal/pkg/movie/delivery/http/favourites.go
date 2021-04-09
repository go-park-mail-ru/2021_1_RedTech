package http

import (
	"Redioteka/internal/pkg/movie"
	"Redioteka/internal/pkg/utils/jsonerrors"
	"Redioteka/internal/pkg/utils/log"
	"Redioteka/internal/pkg/utils/session"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (handler *MovieHandler) SetFavourite(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	urlID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Log.Warn(fmt.Sprint("Error while getting movie id: ", vars["id"]))
		http.Error(w, jsonerrors.URLParams, http.StatusBadRequest)
		return
	}
	id := uint(urlID)

	sess, err := session.GetSession(r)
	if err != nil {
		http.Error(w, jsonerrors.Session, http.StatusUnauthorized)
		return
	}

	switch mux.CurrentRoute(r).GetName() {
	case addFavourite:
		err = handler.MUCase.AddFavourite(id, sess)
	case removeFavourite:
		err = handler.MUCase.RemoveFavourite(id, sess)
	default:
		log.Log.Warn("Unknown name of handler")
		err = errors.New("Wrong url")
	}
	if err != nil {
		http.Error(w, jsonerrors.URLParams, movie.CodeFromError(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}
