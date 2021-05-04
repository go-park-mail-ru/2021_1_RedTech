package http

import (
	"Redioteka/internal/pkg/domain"
	"Redioteka/internal/pkg/utils/jsonerrors"
	"Redioteka/internal/pkg/utils/log"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type SearchHandler struct {
	SUCase domain.SearchUsecase
}

func NewSearchHandlers(router *mux.Router, us domain.SearchUsecase) {
	handler := &SearchHandler{
		SUCase: us,
	}
	router.HandleFunc("/search", handler.Get).Methods("GET", "OPTIONS")
}

func (handler *SearchHandler) Get(w http.ResponseWriter, r *http.Request) {
	query, found := r.URL.Query()["query"]
	if !found {
		log.Log.Warn("Can't parse search query")
		http.Error(w, jsonerrors.URLParams, http.StatusBadRequest)
		return
	}
	queryRes, err := handler.SUCase.Get(query[0])
	if err != nil {
		log.Log.Warn(fmt.Sprintf("Error while getting user id: %s", err))
		http.Error(w, jsonerrors.URLParams, http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(queryRes)
	if err != nil {
		log.Log.Error(err)
		http.Error(w, jsonerrors.JSONEncode, http.StatusInternalServerError)
		return
	}
}
