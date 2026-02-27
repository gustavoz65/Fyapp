package handler

import (
	"net/http"

	"github.com/gustavoz65/finext/internal/config"
	"github.com/gustavoz65/finext/internal/middleware"
	"github.com/gustavoz65/finext/internal/model"
	"github.com/gustavoz65/finext/internal/service"
	"github.com/gustavoz65/finext/internal/validation"
	"github.com/gustavoz65/finext/internal/lib/utils/job"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type AuthHandler struct {
	authService *service.AuthService
	config      *config.Config
	enqueuer    job.Enqueuer
	logger      *zerolog.Logger
}

func NewAuthHandler(authService *service.AuthService, cfg *config.Config, enqueuer job.Enqueuer, logger *zerolog.Logger) *AuthHandler {
	return &AuthHandler{authService: authService, config: cfg, enqueuer: enqueuer, logger: logger}
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

	// Enqueue welcome email (non-blocking)
	if h.enqueuer != nil {
		if task, err := job.NewWelcomeEmailTask(response.User.Email, response.User.FirstName); err == nil {
			if _, err := h.enqueuer.Enqueue(task); err != nil {
				h.logger.Error().Err(err).Str("email", response.User.Email).Msg("failed to enqueue welcome email task from handler")
			}
		}
	}

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

func (h *AuthHandler) SetPassword(c echo.Context) error {
	var req model.SetPasswordRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		return err
	}

	userID := middleware.GetUserID(c)
	if err := h.authService.SetPassword(c.Request().Context(), userID, &req); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// ========================================
// Social Login Handlers
// ========================================

func (h *AuthHandler) SocialLogin(c echo.Context) error {
	var req model.SocialLoginRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		return err
	}

	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()

	response, err := h.authService.SocialLogin(c.Request().Context(), &req, ipAddress, userAgent)
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

func (h *AuthHandler) LinkProvider(c echo.Context) error {
	var req model.LinkProviderRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		return err
	}

	userID := middleware.GetUserID(c)
	if err := h.authService.LinkProvider(c.Request().Context(), userID, &req); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *AuthHandler) UnlinkProvider(c echo.Context) error {
	provider := c.Param("provider")
	if provider == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "provider parameter is required")
	}

	userID := middleware.GetUserID(c)
	if err := h.authService.UnlinkProvider(c.Request().Context(), userID, provider); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *AuthHandler) GetLinkedProviders(c echo.Context) error {
	userID := middleware.GetUserID(c)

	response, err := h.authService.GetLinkedProviders(c.Request().Context(), userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}
