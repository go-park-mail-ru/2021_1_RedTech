package middlewares

import (
	"log"
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
		if _, found := whiteListOrigin[origin]; found {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			log.Printf("Request from unknown host: %s", origin)
		}

		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Accept-Language, "+
			"Content-Language, Content-Type, Content-Encoding")
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}
