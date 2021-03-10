package cors_rest

import (
	"net/http"
	"strings"

	"github.com/gorilla/handlers"
)

// GetCORS returns a handler with configured CORS
func GetCORS() func(http.Handler) http.Handler {
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions})
	originsAllowed := handlers.AllowedOrigins([]string{"http://localhost:4200", "https://analysis.dystopia.systems"})

	return handlers.CORS(headers, methods, originsAllowed)
}

func PassCors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

func preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{"Content-Type", "Accept"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodDelete}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
	return
}
