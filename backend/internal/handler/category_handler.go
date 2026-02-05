package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gustavoz65/Cashing-go/internal/errs"
	"github.com/gustavoz65/Cashing-go/internal/middleware"
	"github.com/gustavoz65/Cashing-go/internal/model"
	"github.com/gustavoz65/Cashing-go/internal/service"
	"github.com/gustavoz65/Cashing-go/internal/validation"
	"github.com/labstack/echo/v4"
)

type CategoryHandler struct {
	categoryService *service.CategoryService
}

func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

func (h *CategoryHandler) GetAll(c echo.Context) error {
	userID := middleware.GetUserID(c)

	// Filtro por tipo opcional
	categoryType := c.QueryParam("type")
	if categoryType != "" {
		categories, err := h.categoryService.GetByType(c.Request().Context(), userID, model.CategoryType(categoryType))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, categories)
	}

	categories, err := h.categoryService.GetAll(c.Request().Context(), userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, categories)
}

func (h *CategoryHandler) GetByID(c echo.Context) error {
	userID := middleware.GetUserID(c)

	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errs.NewBadRequestError("ID de categoria invalido", false, nil, nil, nil)
	}

	category, err := h.categoryService.GetByID(c.Request().Context(), userID, categoryID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) Create(c echo.Context) error {
	var req model.CreateCategoryRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		return err
	}

	userID := middleware.GetUserID(c)

	category, err := h.categoryService.Create(c.Request().Context(), userID, &req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, category)
}

func (h *CategoryHandler) Update(c echo.Context) error {
	var req model.UpdateCategoryRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		return err
	}

	userID := middleware.GetUserID(c)

	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errs.NewBadRequestError("ID de categoria invalido", false, nil, nil, nil)
	}

	category, err := h.categoryService.Update(c.Request().Context(), userID, categoryID, &req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) Delete(c echo.Context) error {
	userID := middleware.GetUserID(c)

	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errs.NewBadRequestError("ID de categoria invalido", false, nil, nil, nil)
	}

	if err := h.categoryService.Delete(c.Request().Context(), userID, categoryID); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
