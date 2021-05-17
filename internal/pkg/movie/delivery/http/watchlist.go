package http

import (
	"Redioteka/internal/pkg/movie"
	"Redioteka/internal/pkg/utils/jsonerrors"
	"Redioteka/internal/pkg/utils/log"
	"Redioteka/internal/pkg/utils/session"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (handler *MovieHandler) SetWatchlist(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	urlID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Log.Warn("Error while getting movie id: " + vars["id"])
		http.Error(w, jsonerrors.URLParams, http.StatusBadRequest)
		return
	}
	id := uint(urlID)

	sess, err := session.GetSession(r)
	if err != nil {
		http.Error(w, jsonerrors.Session, http.StatusUnauthorized)
		return
	}

	qa := r.URL.Query().Get("action")
	switch qa {
	case "save":
		err = handler.MUCase.AddWatchlist(id, sess)
	case "delete":
		err = handler.MUCase.RemoveWatchlist(id, sess)
	default:
		err = movie.BadParamsError
	}
	if err != nil {
		http.Error(w, jsonerrors.URLParams, movie.CodeFromError(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}
