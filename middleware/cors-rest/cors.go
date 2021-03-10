package cors_rest

import (
	"net/http"

	"github.com/vectorman1/analysis/analysis-api/common"

	"github.com/gorilla/handlers"
)

// GetCORS returns a handler with configured CORS
func GetCORS() func(http.Handler) http.Handler {
	headers := handlers.AllowedHeaders(common.GetAllowedHeaders())
	methods := handlers.AllowedMethods(common.GetAllowedMethods())
	origins := handlers.AllowedOrigins([]string{"*"})

	return handlers.CORS(headers, methods, origins)
}
