package handler

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gustavoz65/Fyapp/internal/errs"
	"github.com/gustavoz65/Fyapp/internal/middleware"
	"github.com/gustavoz65/Fyapp/internal/model"
	"github.com/gustavoz65/Fyapp/internal/service"
	"github.com/gustavoz65/Fyapp/internal/utils"
	"github.com/gustavoz65/Fyapp/internal/validation"
	"github.com/labstack/echo/v4"
)

type TransactionHandler struct {
	transactionService *service.TransactionService
}

func NewTransactionHandler(transactionService *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService}
}

func (h *TransactionHandler) GetAll(c echo.Context) error {
	userID := middleware.GetUserID(c)
	qv := utils.NewQueryValidator(c)

	filter := &model.TransactionFilter{
		UserID: userID,
	}

	var err error

	// Validar UUID params
	filter.AccountID, err = qv.GetUUID("account_id", false)
	if err != nil {
		return err
	}

	filter.CategoryID, err = qv.GetUUID("category_id", false)
	if err != nil {
		return err
	}

	// Validar enum params
	typeVal, err := qv.GetEnum("type", false, []string{"income", "expense"})
	if err != nil {
		return err
	}
	if typeVal != nil {
		t := model.TransactionType(*typeVal)
		filter.Type = &t
	}

	// Validar date params
	filter.StartDate, err = qv.GetDate("start_date", false)
	if err != nil {
		return err
	}

	filter.EndDate, err = qv.GetDate("end_date", false)
	if err != nil {
		return err
	}

	// Bool params
	filter.IsPaid = qv.GetBool("is_paid")
	filter.IsRecurring = qv.GetBool("is_recurring")

	// Decimal params
	filter.MinAmount, err = qv.GetDecimal("min_amount", false)
	if err != nil {
		return err
	}

	filter.MaxAmount, err = qv.GetDecimal("max_amount", false)
	if err != nil {
		return err
	}

	// String search (sem validação especial)
	searchTerm, _ := qv.GetString("search", false, 255)
	if searchTerm != nil {
		filter.SearchTerm = *searchTerm
	}

	// Paginação com limites (máximo 100 por página)
	filter.Page, filter.PageSize, err = qv.GetPagination()
	if err != nil {
		return err
	}

	// Sorting
	sortBy, err := qv.GetEnum("sort_by", false, []string{"date", "amount", "description", "created_at"})
	if err != nil {
		return err
	}
	if sortBy != nil {
		filter.SortBy = *sortBy
	}

	sortDir, err := qv.GetEnum("sort_dir", false, []string{"asc", "desc"})
	if err != nil {
		return err
	}
	if sortDir != nil {
		filter.SortDirection = *sortDir
	}

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
