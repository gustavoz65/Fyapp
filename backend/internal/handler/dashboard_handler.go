package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gustavoz65/finext/internal/errs"
	"github.com/gustavoz65/finext/internal/middleware"
	"github.com/gustavoz65/finext/internal/service"
	"github.com/labstack/echo/v4"
)

type DashboardHandler struct {
	dashboardService *service.DashboardService
}

func NewDashboardHandler(dashboardService *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService}
}

func (h *DashboardHandler) GetSummary(c echo.Context) error {
	userID := middleware.GetUserID(c)

	summary, err := h.dashboardService.GetDashboardSummary(c.Request().Context(), userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, summary)
}

func (h *DashboardHandler) GetCashFlow(c echo.Context) error {
	userID := middleware.GetUserID(c)

	startDate, endDate, err := parseDateRange(c)
	if err != nil {
		return err
	}

	report, err := h.dashboardService.GetCashFlowReport(c.Request().Context(), userID, startDate, endDate)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, report)
}

func (h *DashboardHandler) GetIncomeVsExpense(c echo.Context) error {
	userID := middleware.GetUserID(c)

	startDate, endDate, err := parseDateRange(c)
	if err != nil {
		return err
	}

	report, err := h.dashboardService.GetIncomeVsExpenseReport(c.Request().Context(), userID, startDate, endDate)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, report)
}

func (h *DashboardHandler) GetMonthlyComparison(c echo.Context) error {
	userID := middleware.GetUserID(c)

	months := 6
	if m := c.QueryParam("months"); m != "" {
		if parsed, err := strconv.Atoi(m); err == nil && parsed > 0 && parsed <= 24 {
			months = parsed
		}
	}

	comparison, err := h.dashboardService.GetMonthlyComparison(c.Request().Context(), userID, months)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, comparison)
}

func (h *DashboardHandler) GetAccountBalances(c echo.Context) error {
	userID := middleware.GetUserID(c)

	accounts, err := h.dashboardService.GetAccountBalances(c.Request().Context(), userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, accounts)
}

// parseDateRange extrai start_date e end_date dos query params.
// Se nao fornecidos, usa o mes atual como padrao.
func parseDateRange(c echo.Context) (time.Time, time.Time, error) {
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, -1)

	if sd := c.QueryParam("start_date"); sd != "" {
		parsed, err := time.Parse("2006-01-02", sd)
		if err != nil {
			return time.Time{}, time.Time{}, errs.NewBadRequestError("Formato de data de inicio invalido. Use YYYY-MM-DD", false, nil, nil, nil)
		}
		startDate = parsed
	}

	if ed := c.QueryParam("end_date"); ed != "" {
		parsed, err := time.Parse("2006-01-02", ed)
		if err != nil {
			return time.Time{}, time.Time{}, errs.NewBadRequestError("Formato de data de fim invalido. Use YYYY-MM-DD", false, nil, nil, nil)
		}
		endDate = parsed
	}

	if endDate.Before(startDate) {
		return time.Time{}, time.Time{}, errs.NewBadRequestError("A data de fim deve ser posterior a data de inicio", false, nil, nil, nil)
	}

	return startDate, endDate, nil
}
