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
	"github.com/gustavoz65/Cashing-go/internal/validation"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
)

type TransactionHandler struct {
	transactionService *service.TransactionService
}

func NewTransactionHandler(transactionService *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService}
}

func (h *TransactionHandler) GetAll(c echo.Context) error {
	userID := middleware.GetUserID(c)

	filter := &model.TransactionFilter{
		UserID: userID,
	}

	// Parse query params
	if accountID := c.QueryParam("account_id"); accountID != "" {
		id, err := uuid.Parse(accountID)
		if err == nil {
			filter.AccountID = &id
		}
	}

	if categoryID := c.QueryParam("category_id"); categoryID != "" {
		id, err := uuid.Parse(categoryID)
		if err == nil {
			filter.CategoryID = &id
		}
	}

	if txType := c.QueryParam("type"); txType != "" {
		t := model.TransactionType(txType)
		filter.Type = &t
	}

	if startDate := c.QueryParam("start_date"); startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			filter.StartDate = &t
		}
	}

	if endDate := c.QueryParam("end_date"); endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			filter.EndDate = &t
		}
	}

	if isPaid := c.QueryParam("is_paid"); isPaid != "" {
		paid := isPaid == "true"
		filter.IsPaid = &paid
	}

	if isRecurring := c.QueryParam("is_recurring"); isRecurring != "" {
		recurring := isRecurring == "true"
		filter.IsRecurring = &recurring
	}

	if minAmount := c.QueryParam("min_amount"); minAmount != "" {
		if amount, err := decimal.NewFromString(minAmount); err == nil {
			filter.MinAmount = &amount
		}
	}

	if maxAmount := c.QueryParam("max_amount"); maxAmount != "" {
		if amount, err := decimal.NewFromString(maxAmount); err == nil {
			filter.MaxAmount = &amount
		}
	}

	filter.SearchTerm = c.QueryParam("search")

	// Paginacao
	filter.Page, _ = strconv.Atoi(c.QueryParam("page"))
	filter.PageSize, _ = strconv.Atoi(c.QueryParam("page_size"))
	filter.SortBy = c.QueryParam("sort_by")
	filter.SortDirection = c.QueryParam("sort_dir")

	result, err := h.transactionService.GetByFilter(c.Request().Context(), filter)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

func (h *TransactionHandler) GetByID(c echo.Context) error {
	userID := middleware.GetUserID(c)

	txID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errs.NewBadRequestError("ID de transacao invalido", false, nil, nil, nil)
	}

	tx, err := h.transactionService.GetByID(c.Request().Context(), userID, txID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tx)
}

func (h *TransactionHandler) Create(c echo.Context) error {
	var req model.CreateTransactionRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		return err
	}

	userID := middleware.GetUserID(c)

	tx, err := h.transactionService.Create(c.Request().Context(), userID, &req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, tx)
}

func (h *TransactionHandler) Update(c echo.Context) error {
	var req model.UpdateTransactionRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		return err
	}

	userID := middleware.GetUserID(c)

	txID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errs.NewBadRequestError("ID de transacao invalido", false, nil, nil, nil)
	}

	tx, err := h.transactionService.Update(c.Request().Context(), userID, txID, &req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tx)
}

func (h *TransactionHandler) Delete(c echo.Context) error {
	userID := middleware.GetUserID(c)

	txID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errs.NewBadRequestError("ID de transacao invalido", false, nil, nil, nil)
	}

	if err := h.transactionService.Delete(c.Request().Context(), userID, txID); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *TransactionHandler) MarkAsPaid(c echo.Context) error {
	userID := middleware.GetUserID(c)

	txID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errs.NewBadRequestError("ID de transacao invalido", false, nil, nil, nil)
	}

	tx, err := h.transactionService.MarkAsPaid(c.Request().Context(), userID, txID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tx)
}

func (h *TransactionHandler) GetUpcoming(c echo.Context) error {
	userID := middleware.GetUserID(c)

	days := 7
	if d := c.QueryParam("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 {
			days = parsed
		}
	}

	transactions, err := h.transactionService.GetUpcomingBills(c.Request().Context(), userID, days)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, transactions)
}
