package middlewares

import (
	"log"
	"net/http"
)

func (m *GoMiddleware) PanicRecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recovered from panic with err: '%s' on url: %s", err, r.RequestURI)
				http.Error(w, `{"error":"server"}`, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
