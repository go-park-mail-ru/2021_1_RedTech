package middlewares

import (
	"Redioteka/internal/pkg/utils/log"
	"fmt"
	"net/http"
	"time"
)

func (m *GoMiddleware) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Log.Info(fmt.Sprintf("%s %s %s %s", r.Method, r.RequestURI, r.RemoteAddr, time.Since(start)))
	})
}
