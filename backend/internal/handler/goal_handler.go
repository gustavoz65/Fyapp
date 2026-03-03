package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/gustavoz65/Fyapp/internal/errs"
	"github.com/gustavoz65/Fyapp/internal/middleware"
	"github.com/gustavoz65/Fyapp/internal/model"
	"github.com/gustavoz65/Fyapp/internal/service"
	"github.com/gustavoz65/Fyapp/internal/validation"
)

type GoalHandler struct {
	goalService *service.GoalService
}

func NewGoalHandler(goalService *service.GoalService) *GoalHandler {
	return &GoalHandler{goalService: goalService}
}

func (h *GoalHandler) GetAll(c echo.Context) error {
	userID := middleware.GetUserID(c)

	if c.QueryParam("active") == "true" {
		goals, err := h.goalService.GetActive(c.Request().Context(), userID)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, goals)
	}

	goals, err := h.goalService.GetAll(c.Request().Context(), userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, goals)
}

func (h *GoalHandler) GetByID(c echo.Context) error {
	userID := middleware.GetUserID(c)

	goalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errs.NewBadRequestError("ID de meta invalido", false, nil, nil, nil)
	}

	goal, err := h.goalService.GetByID(c.Request().Context(), userID, goalID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, goal)
}

func (h *GoalHandler) Create(c echo.Context) error {
	var req model.CreateGoalRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		return err
	}

	userID := middleware.GetUserID(c)

	goal, err := h.goalService.Create(c.Request().Context(), userID, &req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, goal)
}

func (h *GoalHandler) Update(c echo.Context) error {
	var req model.UpdateGoalRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		return err
	}

	userID := middleware.GetUserID(c)

	goalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errs.NewBadRequestError("ID de meta invalido", false, nil, nil, nil)
	}

	goal, err := h.goalService.Update(c.Request().Context(), userID, goalID, &req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, goal)
}

func (h *GoalHandler) Delete(c echo.Context) error {
	userID := middleware.GetUserID(c)

	goalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errs.NewBadRequestError("ID de meta invalido", false, nil, nil, nil)
	}

	if err := h.goalService.Delete(c.Request().Context(), userID, goalID); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *GoalHandler) AddContribution(c echo.Context) error {
	var req model.CreateGoalContributionRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		return err
	}

	userID := middleware.GetUserID(c)

	goalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errs.NewBadRequestError("ID de meta invalido", false, nil, nil, nil)
	}

	contribution, err := h.goalService.AddContribution(c.Request().Context(), userID, goalID, &req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, contribution)
}

func (h *GoalHandler) GetContributions(c echo.Context) error {
	userID := middleware.GetUserID(c)

	goalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errs.NewBadRequestError("ID de meta invalido", false, nil, nil, nil)
	}

	contributions, err := h.goalService.GetContributions(c.Request().Context(), userID, goalID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, contributions)
}

func (h *GoalHandler) GetSummary(c echo.Context) error {
	userID := middleware.GetUserID(c)

	summary, err := h.goalService.GetSummary(c.Request().Context(), userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, summary)
}
