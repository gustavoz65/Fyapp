package handler

import (
	"net/http"

	"github.com/gustavoz65/Cashing-go/internal/middleware"
	"github.com/gustavoz65/Cashing-go/internal/model"
	"github.com/gustavoz65/Cashing-go/internal/service"
	"github.com/gustavoz65/Cashing-go/internal/validation"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req model.RegisterRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		return err
	}

	user, err := h.authService.Register(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, user)
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

	return c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var req model.RefreshTokenRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		return err
	}

	response, err := h.authService.RefreshToken(c.Request().Context(), req.RefreshToken)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) Logout(c echo.Context) error {
	var req model.RefreshTokenRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		return err
	}

	if err := h.authService.Logout(c.Request().Context(), req.RefreshToken); err != nil {
		return err
	}

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
