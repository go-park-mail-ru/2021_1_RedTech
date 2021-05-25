package http

import (
	"Redioteka/internal/pkg/actor"
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/utils/jsonerrors"
	"Redioteka/internal/pkg/utils/log"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type ActorHandler struct {
	AUCase domain.ActorUsecase
}

func NewActorHandlers(router *mux.Router, us domain.ActorUsecase) {
	handler := &ActorHandler{
		AUCase: us,
	}
	router.HandleFunc("/actors/{id:[0-9]+}", handler.Get).Methods("GET", "OPTIONS")
}

func (handler *ActorHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id64, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		log.Log.Warn(fmt.Sprintf("error while getting actor id: {#err}"))
		http.Error(w, jsonerrors.URLParams, http.StatusBadRequest)
		return
	}
	id := uint(id64)

	foundActor, err := handler.AUCase.GetById(id)
	if err != nil {
		http.Error(w, jsonerrors.JSONMessage("not found"), actor.CodeFromError(err))
		return
	}
	err = json.NewEncoder(w).Encode(foundActor)
	if err != nil {
		log.Log.Error(err)
		http.Error(w, jsonerrors.JSONEncode, http.StatusInternalServerError)
		return
	}
}
