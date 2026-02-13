package handler

import (
	"net/http"

	"github.com/gustavoz65/Cashing-go/internal/config"
	"github.com/gustavoz65/Cashing-go/internal/middleware"
	"github.com/gustavoz65/Cashing-go/internal/model"
	"github.com/gustavoz65/Cashing-go/internal/service"
	"github.com/gustavoz65/Cashing-go/internal/validation"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService *service.AuthService
	config      *config.Config
}

func NewAuthHandler(authService *service.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{authService: authService, config: cfg}
}

func (h *AuthHandler) setAuthCookies(c echo.Context, accessToken, refreshToken string) {
	isProduction := h.config.Primary.Env == "production"

	c.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   isProduction,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   h.config.Auth.AccessTokenDuration * 60,
	})

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/api/v1/auth",
		HttpOnly: true,
		Secure:   isProduction,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   h.config.Auth.RefreshTokenDuration * 3600,
	})
}

func (h *AuthHandler) clearAuthCookies(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/api/v1/auth",
		HttpOnly: true,
		MaxAge:   -1,
	})
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req model.RegisterRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		return err
	}

	response, err := h.authService.Register(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	h.setAuthCookies(c, response.AccessToken, response.RefreshToken)

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"user":       response.User,
		"expires_at": response.ExpiresAt,
	})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req model.LoginRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		return err
	}

	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()

	response, err := h.authService.Login(c.Request().Context(), &req, ipAddress, userAgent)
	if err != nil {
		return err
	}

	// Setar cookies httpOnly
	h.setAuthCookies(c, response.AccessToken, response.RefreshToken)

	// Retornar response SEM tokens (por segurança)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"user":       response.User,
		"expires_at": response.ExpiresAt,
	})
}

func (h *AuthHandler) RefreshToken(c echo.Context) error {
	// Tentar ler refresh_token do cookie primeiro
	var refreshToken string
	cookie, err := c.Cookie("refresh_token")
	if err == nil && cookie.Value != "" {
		refreshToken = cookie.Value
	} else {
		// Fallback para body (backward compatibility)
		var req model.RefreshTokenRequest
		if err := validation.BindAndValidate(c, &req); err != nil {
			return err
		}
		refreshToken = req.RefreshToken
	}

	response, err := h.authService.RefreshToken(c.Request().Context(), refreshToken)
	if err != nil {
		return err
	}

	// Setar novos cookies httpOnly
	h.setAuthCookies(c, response.AccessToken, response.RefreshToken)

	// Retornar response SEM tokens (por segurança)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"expires_at": response.ExpiresAt,
	})
}

func (h *AuthHandler) Logout(c echo.Context) error {
	// Tentar ler refresh_token do cookie primeiro
	var refreshToken string
	cookie, err := c.Cookie("refresh_token")
	if err == nil && cookie.Value != "" {
		refreshToken = cookie.Value
	} else {
		// Fallback para body (backward compatibility)
		var req model.RefreshTokenRequest
		if err := validation.BindAndValidate(c, &req); err != nil {
			return err
		}
		refreshToken = req.RefreshToken
	}

	if err := h.authService.Logout(c.Request().Context(), refreshToken); err != nil {
		return err
	}

	// Limpar cookies
	h.clearAuthCookies(c)

	return c.NoContent(http.StatusNoContent)
}

func (h *AuthHandler) ChangePassword(c echo.Context) error {
	var req model.ChangePasswordRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		return err
	}

	userID := middleware.GetUserID(c)
	if err := h.authService.ChangePassword(c.Request().Context(), userID, &req); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
