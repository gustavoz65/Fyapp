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

type BankAccountHandler struct {
	accountService *service.BankAccountService
}

func NewBankAccountHandler(accountService *service.BankAccountService) *BankAccountHandler {
	return &BankAccountHandler{accountService: accountService}
}

func (h *BankAccountHandler) GetAll(c echo.Context) error {
	userID := middleware.GetUserID(c)

	accounts, err := h.accountService.GetAll(c.Request().Context(), userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, accounts)
}

func (h *BankAccountHandler) GetByID(c echo.Context) error {
	userID := middleware.GetUserID(c)

	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errs.NewBadRequestError("ID de conta invalido", false, nil, nil, nil)
	}

	account, err := h.accountService.GetByID(c.Request().Context(), userID, accountID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, account)
}

func (h *BankAccountHandler) Create(c echo.Context) error {
	var req model.CreateBankAccountRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		return err
	}

	userID := middleware.GetUserID(c)

	account, err := h.accountService.Create(c.Request().Context(), userID, &req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, account)
}

func (h *BankAccountHandler) Update(c echo.Context) error {
	var req model.UpdateBankAccountRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		return err
	}

	userID := middleware.GetUserID(c)

	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errs.NewBadRequestError("ID de conta invalido", false, nil, nil, nil)
	}

	account, err := h.accountService.Update(c.Request().Context(), userID, accountID, &req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, account)
}

func (h *BankAccountHandler) Delete(c echo.Context) error {
	userID := middleware.GetUserID(c)

	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return errs.NewBadRequestError("ID de conta invalido", false, nil, nil, nil)
	}

	if err := h.accountService.Delete(c.Request().Context(), userID, accountID); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *BankAccountHandler) GetTotalBalance(c echo.Context) error {
	userID := middleware.GetUserID(c)

	total, err := h.accountService.GetTotalBalance(c.Request().Context(), userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"total_balance": total,
	})
}
