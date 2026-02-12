package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

// CORSMiddleware configura CORS baseado nas origens permitidas
func CORSMiddleware(allowedOrigins []string) echo.MiddlewareFunc {
	if len(allowedOrigins) == 0 {
		// Default para desenvolvimento - NÃO usar "*" em produção com credentials!
		allowedOrigins = []string{"http://localhost:3000"}
	}

	// IMPORTANTE: Quando AllowCredentials=true, não pode usar "*"
	// Precisa ser lista específica de origens
	allowCredentials := true
	for _, origin := range allowedOrigins {
		if origin == "*" {
			allowCredentials = false
			break
		}
	}

	return echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins: allowedOrigins,
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			"X-Requested-With",
		},
		AllowCredentials: allowCredentials,
		MaxAge:           86400,
	})
}
