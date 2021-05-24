package middlewares

import (
	"Redioteka/internal/pkg/utils/jsonerrors"
	"Redioteka/internal/pkg/utils/log"
	"net/http"
)

var whiteListOrigin = map[string]struct{}{
	"http://localhost":            {},
	"http://localhost:3000":       {},
	"http://127.0.0.1":            {},
	"http://127.0.0.1:3000":       {},
	"https://redioteka.com":       {},
	"https://redioteka.com:3000":  {},
	"https://89.208.198.192":      {},
	"https://89.208.198.192:3000": {},
}

func (m *GoMiddleware) CORSMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if r.Method != http.MethodGet {
			if _, found := whiteListOrigin[origin]; found {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else {
				log.Log.Warn("Request from unknown host: " + origin)
				http.Error(w, jsonerrors.JSONMessage("unknown origin"), http.StatusMethodNotAllowed)
				return
			}
		}

		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Accept-Language, "+
			"Content-Language, Content-Type, Content-Encoding, X-CSRF-Token")
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}
