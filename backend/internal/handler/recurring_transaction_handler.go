package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gustavoz65/Cashing-go/internal/errs"
	"github.com/gustavoz65/Cashing-go/internal/middleware"
	"github.com/gustavoz65/Cashing-go/internal/model"
	"github.com/gustavoz65/Cashing-go/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
)

type RecurringTransactionHandler struct {
	service *service.RecurringTransactionService
}

func NewRecurringTransactionHandler(service *service.RecurringTransactionService) *RecurringTransactionHandler {
	return &RecurringTransactionHandler{service: service}
}

type CreateRecurringRequest struct {
	BankAccountID uuid.UUID  `json:"bank_account_id" validate:"required"`
	CategoryID    *uuid.UUID `json:"category_id"`
	Type          string     `json:"type" validate:"required,oneof=income expense"`
	Amount        float64    `json:"amount" validate:"required,gt=0"`
	Description   string     `json:"description" validate:"required"`
	Frequency     string     `json:"frequency" validate:"required,oneof=daily weekly biweekly monthly quarterly yearly"`
	DayOfMonth    *int       `json:"day_of_month"`
	DayOfWeek     *int       `json:"day_of_week"`
	StartDate     *time.Time `json:"start_date"`
	EndDate       *time.Time `json:"end_date"`
	AutoConfirm   bool       `json:"auto_confirm"`
}

type UpdateRecurringRequest struct {
	BankAccountID uuid.UUID  `json:"bank_account_id" validate:"required"`
	CategoryID    *uuid.UUID `json:"category_id"`
	Type          string     `json:"type" validate:"required,oneof=income expense"`
	Amount        float64    `json:"amount" validate:"required,gt=0"`
	Description   string     `json:"description" validate:"required"`
	Frequency     string     `json:"frequency" validate:"required,oneof=daily weekly biweekly monthly quarterly yearly"`
	DayOfMonth    *int       `json:"day_of_month"`
	DayOfWeek     *int       `json:"day_of_week"`
	StartDate     time.Time  `json:"start_date" validate:"required"`
	EndDate       *time.Time `json:"end_date"`
	AutoConfirm   bool       `json:"auto_confirm"`
}

type ToggleActiveRequest struct {
	IsActive bool `json:"is_active"`
}

func (h *RecurringTransactionHandler) Create(c echo.Context) error {
	userID := middleware.GetUserID(c)

	var req CreateRecurringRequest
	if err := c.Bind(&req); err != nil {
		return errs.NewBadRequestError("Invalid request body", false, nil, nil, nil)
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	startDate := time.Now()
	if req.StartDate != nil {
		startDate = *req.StartDate
	}

	nextOccurrence := startDate
	if req.Frequency == string(model.FrequencyMonthly) && req.DayOfMonth != nil {
		year, month, _ := startDate.Date()
		nextOccurrence = time.Date(year, month, *req.DayOfMonth, 0, 0, 0, 0, startDate.Location())
		if nextOccurrence.Before(time.Now()) {
			nextOccurrence = nextOccurrence.AddDate(0, 1, 0)
		}
	}

	rt := &model.RecurringTransaction{
		UserID:         userID,
		BankAccountID:  req.BankAccountID,
		CategoryID:     req.CategoryID,
		Type:           model.TransactionType(req.Type),
		Amount:         decimal.NewFromFloat(req.Amount),
		Description:    req.Description,
		Frequency:      model.RecurringFrequency(req.Frequency),
		DayOfMonth:     req.DayOfMonth,
		DayOfWeek:      req.DayOfWeek,
		StartDate:      startDate,
		EndDate:        req.EndDate,
		NextOccurrence: nextOccurrence,
		IsActive:       true,
		AutoConfirm:    req.AutoConfirm,
	}

	if err := h.service.Create(c.Request().Context(), rt); err != nil {
		return errs.NewBadRequestError(err.Error(), false, nil, nil, nil)
	}

	return c.JSON(http.StatusCreated, rt)
}

func (h *RecurringTransactionHandler) GetAll(c echo.Context) error {
	userID := middleware.GetUserID(c)

	var isActive *bool
	if activeParam := c.QueryParam("is_active"); activeParam != "" {
		if val, err := strconv.ParseBool(activeParam); err == nil {
			isActive = &val
		}
	}

	recurrings, err := h.service.GetAll(c.Request().Context(), userID, isActive)
	if err != nil {
		return errs.NewInternalServerError()
	}

	return c.JSON(http.StatusOK, recurrings)
}

func (h *RecurringTransactionHandler) GetByID(c echo.Context) error {
	userID := middleware.GetUserID(c)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errs.NewBadRequestError("Invalid ID", false, nil, nil, nil)
	}

	rt, err := h.service.GetByID(c.Request().Context(), id, userID)
	if err != nil {
		return errs.NewNotFoundError("Recurring transaction not found", false, nil)
	}

	return c.JSON(http.StatusOK, rt)
}

func (h *RecurringTransactionHandler) Update(c echo.Context) error {
	userID := middleware.GetUserID(c)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errs.NewBadRequestError("Invalid ID", false, nil, nil, nil)
	}

	var req UpdateRecurringRequest
	if err := c.Bind(&req); err != nil {
		return errs.NewBadRequestError("Invalid request body", false, nil, nil, nil)
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	existing, err := h.service.GetByID(c.Request().Context(), id, userID)
	if err != nil {
		return errs.NewNotFoundError("Recurring transaction not found", false, nil)
	}

	existing.BankAccountID = req.BankAccountID
	existing.CategoryID = req.CategoryID
	existing.Type = model.TransactionType(req.Type)
	existing.Amount = decimal.NewFromFloat(req.Amount)
	existing.Description = req.Description
	existing.Frequency = model.RecurringFrequency(req.Frequency)
	existing.DayOfMonth = req.DayOfMonth
	existing.DayOfWeek = req.DayOfWeek
	existing.StartDate = req.StartDate
	existing.EndDate = req.EndDate
	existing.AutoConfirm = req.AutoConfirm

	if err := h.service.Update(c.Request().Context(), existing); err != nil {
		return errs.NewBadRequestError(err.Error(), false, nil, nil, nil)
	}

	return c.JSON(http.StatusOK, existing)
}

func (h *RecurringTransactionHandler) Delete(c echo.Context) error {
	userID := middleware.GetUserID(c)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errs.NewBadRequestError("Invalid ID", false, nil, nil, nil)
	}

	if err := h.service.Delete(c.Request().Context(), id, userID); err != nil {
		return errs.NewNotFoundError("Recurring transaction not found", false, nil)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *RecurringTransactionHandler) ToggleActive(c echo.Context) error {
	userID := middleware.GetUserID(c)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errs.NewBadRequestError("Invalid ID", false, nil, nil, nil)
	}

	var req ToggleActiveRequest
	if err := c.Bind(&req); err != nil {
		return errs.NewBadRequestError("Invalid request body", false, nil, nil, nil)
	}

	if err := h.service.ToggleActive(c.Request().Context(), id, userID, req.IsActive); err != nil {
		return errs.NewNotFoundError("Recurring transaction not found", false, nil)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *RecurringTransactionHandler) GetStats(c echo.Context) error {
	userID := middleware.GetUserID(c)

	stats, err := h.service.GetStats(c.Request().Context(), userID)
	if err != nil {
		return errs.NewInternalServerError()
	}

	return c.JSON(http.StatusOK, stats)
}

func (h *RecurringTransactionHandler) GetUpcoming(c echo.Context) error {
	userID := middleware.GetUserID(c)

	days := 30
	if daysParam := c.QueryParam("days"); daysParam != "" {
		if val, err := strconv.Atoi(daysParam); err == nil && val > 0 {
			days = val
		}
	}

	upcoming, err := h.service.GetUpcoming(c.Request().Context(), userID, days)
	if err != nil {
		return errs.NewInternalServerError()
	}

	return c.JSON(http.StatusOK, upcoming)
}
