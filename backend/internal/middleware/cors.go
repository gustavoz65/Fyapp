package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

// CORSMiddleware configura CORS baseado nas origens permitidas
func CORSMiddleware(allowedOrigins []string) echo.MiddlewareFunc {
	if len(allowedOrigins) == 0 {
		allowedOrigins = []string{"*"}
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
		AllowCredentials: true,
		MaxAge:           86400,
	})
}
