package http

import (
	"Redioteka/internal/pkg/user"
	"Redioteka/internal/pkg/utils/jsonerrors"
	"Redioteka/internal/pkg/utils/log"
	"Redioteka/internal/pkg/utils/session"
	"net/http"
)

func (handler *MovieHandler) Like(w http.ResponseWriter, r *http.Request) {
	sess, err := session.GetSession(r)
	if err != nil || session.Manager.Check(sess) != nil {
		log.Log.Warn("Trying to like while unauthorized")
		http.Error(w, jsonerrors.JSONMessage("unauthorized"), user.CodeFromError(user.UnauthorizedError))
		return
	}
}

func (handler *MovieHandler) Dislike(w http.ResponseWriter, r *http.Request) {
	sess, err := session.GetSession(r)
	if err != nil || session.Manager.Check(sess) != nil {
		log.Log.Warn("Trying to like while unauthorized")
		http.Error(w, jsonerrors.JSONMessage("unauthorized"), user.CodeFromError(user.UnauthorizedError))
		return
	}
}
