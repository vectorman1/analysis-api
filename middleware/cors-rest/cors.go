package cors_rest

import (
	"net/http"

	"github.com/gorilla/handlers"
)

// GetCORS returns a handler with configured CORS
func GetCORS() func(http.Handler) http.Handler {
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions})
	originsAllowed := handlers.AllowedOrigins([]string{"*", "localhost", "*.dystopia.systems"})

	return handlers.CORS(headers, methods, originsAllowed)
}
