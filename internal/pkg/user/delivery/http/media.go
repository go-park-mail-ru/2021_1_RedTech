package http

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/user"
	"Redioteka/internal/pkg/utils/jsonerrors"
	"Redioteka/internal/pkg/utils/log"
	"Redioteka/internal/pkg/utils/session"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (handler *UserHandler) GetMedia(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	urlID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Log.Warn("Error while getting user id: " + vars["id"])
		http.Error(w, jsonerrors.URLParams, http.StatusBadRequest)
		return
	}
	id := uint(urlID)

	sess, err := session.GetSession(r)
	if err != nil {
		http.Error(w, jsonerrors.Session, http.StatusUnauthorized)
		return
	}

	favouriteMovies, err := handler.UUsecase.GetFavourites(id, sess)
	if err != nil {
		http.Error(w, jsonerrors.JSONMessage("not found"), user.CodeFromError(err))
		return
	}

	err = json.NewEncoder(w).Encode(domain.UserFavourites{Favourites: favouriteMovies})
	if err != nil {
		log.Log.Error(err)
		http.Error(w, jsonerrors.JSONEncode, http.StatusInternalServerError)
		return
	}
}
