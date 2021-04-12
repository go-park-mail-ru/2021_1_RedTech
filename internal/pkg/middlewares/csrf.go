package middlewares

import (
	"Redioteka/internal/pkg/utils/jsonerrors"
	"Redioteka/internal/pkg/utils/log"
	"net/http"
)

func (m *GoMiddleware) CSRFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			_, err := r.Cookie("csrf_token")
			if err == http.ErrNoCookie {
				log.Log.Warn("No csrf-token for POST query")
				http.Error(w, jsonerrors.CSRF, http.StatusBadRequest)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
