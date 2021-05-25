package http

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/movie"
	"Redioteka/internal/pkg/user"
	"Redioteka/internal/pkg/utils/jsonerrors"
	"Redioteka/internal/pkg/utils/log"
	"Redioteka/internal/pkg/utils/session"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type StreamHandler struct {
	SUCase domain.StreamUsecase
	sessionManager session.SessionManager
}

func NewStreamHandlers(router *mux.Router, us domain.StreamUsecase, sm session.SessionManager) {
	handler := &StreamHandler{
		SUCase: us,
		sessionManager: sm,
	}
	router.HandleFunc("/media/movie/{id:[0-9]+}/stream", handler.Stream).Methods("GET", "OPTIONS")
}

func (handler *StreamHandler) Stream(w http.ResponseWriter, r *http.Request) {
	sess, err := session.GetSession(r)
	if err != nil || session.Manager.Check(sess) != nil {
		log.Log.Warn("Trying to get stream while unauthorized")
		http.Error(w, jsonerrors.JSONMessage("unauthorized"), user.CodeFromError(user.UnauthorizedError))
		return
	}

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

	foundStream, err := handler.SUCase.GetStream(id)
	if err != nil {
		http.Error(w, jsonerrors.JSONMessage("not found"), movie.CodeFromError(err))
		return
	}

	err = json.NewEncoder(w).Encode(foundStream)
	if err != nil {
		log.Log.Error(err)
		http.Error(w, jsonerrors.JSONEncode, http.StatusInternalServerError)
		return
	}
}
