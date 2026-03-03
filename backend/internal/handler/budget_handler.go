package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gustavoz65/Fyapp/internal/errs"
	"github.com/gustavoz65/Fyapp/internal/middleware"
	"github.com/gustavoz65/Fyapp/internal/model"
	"github.com/gustavoz65/Fyapp/internal/service"
	"github.com/gustavoz65/Fyapp/internal/validation"
	"github.com/labstack/echo/v4"
)

type BudgetHandler struct {
	budgetService *service.BudgetService
}

func NewBudgetHandler(budgetService *service.BudgetService) *BudgetHandler {
	return &BudgetHandler{budgetService: budgetService}
}

func (h *BudgetHandler) GetAll(c echo.Context) error {
	userID := middleware.GetUserID(c)

	// Se o query param "active" estiver presente, retorna apenas ativos
	if c.QueryParam("active") == "true" {
		budgets, err := h.budgetService.GetActive(c.Request().Context(), userID)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, budgets)
	}

	budgets, err := h.budgetService.GetAll(c.Request().Context(), userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, budgets)
}

func (h *BudgetHandler) GetByID(c echo.Context) error {
	userID := middleware.GetUserID(c)

	budgetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errs.NewBadRequestError("ID de orcamento invalido", false, nil, nil, nil)
	}

	budget, err := h.budgetService.GetByID(c.Request().Context(), userID, budgetID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, budget)
}

func (h *BudgetHandler) Create(c echo.Context) error {
	var req model.CreateBudgetRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		return err
	}

	userID := middleware.GetUserID(c)

	budget, err := h.budgetService.Create(c.Request().Context(), userID, &req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, budget)
}

func (h *BudgetHandler) Update(c echo.Context) error {
	var req model.UpdateBudgetRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		return err
	}

	userID := middleware.GetUserID(c)

	budgetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errs.NewBadRequestError("ID de orcamento invalido", false, nil, nil, nil)
	}

	budget, err := h.budgetService.Update(c.Request().Context(), userID, budgetID, &req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, budget)
}

func (h *BudgetHandler) Delete(c echo.Context) error {
	userID := middleware.GetUserID(c)

	budgetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errs.NewBadRequestError("ID de orcamento invalido", false, nil, nil, nil)
	}

	if err := h.budgetService.Delete(c.Request().Context(), userID, budgetID); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *BudgetHandler) GetSummary(c echo.Context) error {
	userID := middleware.GetUserID(c)

	summary, err := h.budgetService.GetBudgetSummary(c.Request().Context(), userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, summary)
}
