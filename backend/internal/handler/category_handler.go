package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gustavoz65/finext/internal/errs"
	"github.com/gustavoz65/finext/internal/middleware"
	"github.com/gustavoz65/finext/internal/model"
	"github.com/gustavoz65/finext/internal/service"
	"github.com/gustavoz65/finext/internal/validation"
	"github.com/labstack/echo/v4"
)

type CategoryHandler struct {
	categoryService   *service.CategoryService
	categorizationSvc *service.CategorizationService
}

func NewCategoryHandler(categoryService *service.CategoryService, categorizationSvc *service.CategorizationService) *CategoryHandler {
	return &CategoryHandler{
		categoryService:   categoryService,
		categorizationSvc: categorizationSvc,
	}
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

// SuggestCategory sugere categorias baseado na descricao fornecida
func (h *CategoryHandler) SuggestCategory(c echo.Context) error {
	userID := middleware.GetUserID(c)

	description := c.QueryParam("description")
	if description == "" {
		return errs.NewBadRequestError("descricao e obrigatoria", false, nil, nil, nil)
	}

	result, err := h.categorizationSvc.SuggestCategory(c.Request().Context(), userID, description)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}
