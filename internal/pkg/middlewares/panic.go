package middlewares

import (
	"Redioteka/internal/pkg/utils/log"
	"fmt"
	"net/http"
)

func (m *GoMiddleware) PanicRecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Log.Warn(fmt.Sprintf("Recovered from panic with err: '%s' on url: %s", err, r.RequestURI))
				http.Error(w, `{"error":"server"}`, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
