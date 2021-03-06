package http

import (
	"Redioteka/internal/pkg/middlewares"
	"Redioteka/internal/pkg/subscription/delivery/grpc/proto"
	"Redioteka/internal/pkg/user"
	"Redioteka/internal/pkg/utils/jsonerrors"
	"Redioteka/internal/pkg/utils/log"
	"Redioteka/internal/pkg/utils/session"
	"context"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type SubscriptionHandler struct {
	grpcHandler    proto.SubscriptionClient
	SessionManager session.SessionManager
}

func NewSubscriptionHandlers(router *mux.Router, handlers proto.SubscriptionClient, sm session.SessionManager) {
	handler := &SubscriptionHandler{
		grpcHandler:    handlers,
		SessionManager: sm,
	}

	middL := middlewares.InitMiddleware()
	subrouter := router.NewRoute().Subrouter()
	subrouter.Use(middL.LoggingMiddleware)
	s := subrouter.NewRoute().Subrouter()
	s.Use(middL.CORSMiddleware)
	s.Use(middL.CSRFMiddleware)

	subrouter.HandleFunc("/api/subscriptions", handler.Create).Methods("POST", "OPTIONS")

	s.HandleFunc("/api/subscriptions", handler.Delete).Methods("DELETE", "OPTIONS")
}

func (sh *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	codepro, _ := strconv.ParseBool(r.FormValue("codepro"))
	unaccepted, _ := strconv.ParseBool(r.FormValue("unaccepted"))
	_, err := sh.grpcHandler.Create(context.Background(), &proto.Payment{
		Type:        r.FormValue("notification_type"),
		OperationID: r.FormValue("operation_id"),
		Amount:      r.FormValue("amount"),
		Currency:    r.FormValue("currency"),
		DateTime:    r.FormValue("datetime"),
		Sender:      r.FormValue("sender"),
		CodePro:     codepro,
		Label:       r.FormValue("label"),
		Unaccepted:  unaccepted,
		Hash:        r.FormValue("sha1_hash"),
	})
	if err != nil {
		log.Log.Error(err)
		return
	}
	log.Log.Info("Payment was accepted")
}

func (sh *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	sess, err := session.GetSession(r)
	if err != nil || sh.SessionManager.Check(sess) != nil {
		log.Log.Info("User is unauthorized")
		http.Error(w, jsonerrors.JSONMessage("unauthorized"), user.CodeFromError(user.UnauthorizedError))
		return
	}

	_, err = sh.grpcHandler.Delete(context.Background(), &proto.UserId{
		ID: uint64(sess.UserID),
	})
	if err != nil {
		log.Log.Error(err)
		http.Error(w, jsonerrors.JSONMessage("no subscription"), http.StatusBadRequest)
		return
	}
}
