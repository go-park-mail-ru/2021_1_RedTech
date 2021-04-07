package http

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/user"
	"Redioteka/internal/pkg/utils/jsonerrors"
	"Redioteka/internal/pkg/utils/log"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (handler *UserHandler) GetMedia(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	urlID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Log.Warn(fmt.Sprint("Error while getting user id: ", vars["id"]))
		http.Error(w, jsonerrors.URLParams, http.StatusBadRequest)
		return
	}
	id := uint(urlID)

	sess, err := getSession(r)
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
