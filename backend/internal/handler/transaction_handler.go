package handler

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"

	"github.com/gustavoz65/Fyapp/internal/errs"
	"github.com/gustavoz65/Fyapp/internal/lib/utils/job"
	"github.com/gustavoz65/Fyapp/internal/middleware"
	"github.com/gustavoz65/Fyapp/internal/model"
	"github.com/gustavoz65/Fyapp/internal/service"
	"github.com/gustavoz65/Fyapp/internal/utils"
	"github.com/gustavoz65/Fyapp/internal/validation"
)

type TransactionHandler struct {
	transactionService *service.TransactionService
	accountService     *service.BankAccountService
	jobService         *job.JobService
}

func NewTransactionHandler(transactionService *service.TransactionService, accountService *service.BankAccountService, jobService *job.JobService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
		accountService:     accountService,
		jobService:         jobService,
	}
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

	sourceVal, err := qv.GetEnum("source", false, []string{"manual", "bank_sync", "recurring"})
	if err != nil {
		return err
	}
	if sourceVal != nil {
		s := model.TransactionSource(*sourceVal)
		filter.Source = &s
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

func (h *TransactionHandler) Import(c echo.Context) error {
	userID := middleware.GetUserID(c)

	bankType := c.FormValue("bank_type")
	if bankType == "" {
		bankType = "generic"
	}

	file, err := c.FormFile("file")
	if err != nil {
		return errs.NewBadRequestError("file is required", false, nil, nil, nil)
	}

	if file.Header.Get("Content-Type") != "text/csv" &&
		file.Header.Get("Content-Type") != "application/vnd.ms-excel" &&
		file.Header.Get("Content-Type") != "application/csv" {
		if len(file.Filename) < 4 || file.Filename[len(file.Filename)-4:] != ".csv" {
			return errs.NewBadRequestError("apenas arquivos CSV são suportados", false, nil, nil, nil)
		}
	}

	const maxFileSize = 5 * 1024 * 1024
	if file.Size > maxFileSize {
		return errs.NewBadRequestError("tamanho do arquivo excede o limite de 5MB", false, nil, nil, nil)
	}

	if file.Size == 0 {
		return errs.NewBadRequestError("arquivo está vazio", false, nil, nil, nil)
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	csvData, err := io.ReadAll(io.LimitReader(src, maxFileSize))
	if err != nil {
		return errs.NewBadRequestError("falha ao ler arquivo", false, nil, nil, nil)
	}

	if err := h.validateCSVContent(string(csvData)); err != nil {
		return err
	}

	// Handle account creation or use existing account
	var bankAccountID uuid.UUID
	createAccount := c.FormValue("create_account") == "true"

	if createAccount {
		// Create new account
		accountName := c.FormValue("account_name")
		if accountName == "" {
			return errs.NewBadRequestError("account_name é obrigatório ao criar nova conta", false, nil, nil, nil)
		}

		// Extract initial balance from CSV if it's Banco do Brasil
		initialBalance, err := h.extractInitialBalance(string(csvData), bankType)
		if err != nil {
			// If we can't extract balance, default to zero
			initialBalance = "0"
		}

		// Parse initial balance
		initialBalanceDecimal, err := decimal.NewFromString(initialBalance)
		if err != nil {
			initialBalanceDecimal = decimal.Zero
		}

		// Create bank account
		req := &model.CreateBankAccountRequest{
			Name:           accountName,
			AccountType:    model.AccountTypeChecking, // Default to checking
			InitialBalance: initialBalanceDecimal,
			Currency:       "BRL",
			Color:          "#10B981",
			Icon:           "bank",
		}

		// Set bank name based on bank type
		switch bankType {
		case "nubank":
			req.BankName = "Nubank"
			req.BankCode = "260"
		case "bb":
			req.BankName = "Banco do Brasil"
			req.BankCode = "001"
		case "inter":
			req.BankName = "Inter"
			req.BankCode = "077"
		case "itau":
			req.BankName = "Itaú"
			req.BankCode = "341"
		}

		account, err := h.accountService.Create(c.Request().Context(), userID, req)
		if err != nil {
			return errs.NewBadRequestError(fmt.Sprintf("Erro ao criar conta: %v", err), false, nil, nil, nil)
		}

		bankAccountID = account.ID
	} else {
		// Use existing account
		bankAccountIDStr := c.FormValue("bank_account_id")
		if bankAccountIDStr == "" {
			return errs.NewBadRequestError("bank_account_id é obrigatório", false, nil, nil, nil)
		}

		bankAccountID, err = uuid.Parse(bankAccountIDStr)
		if err != nil {
			return errs.NewBadRequestError("formato de bank_account_id inválido", false, nil, nil, nil)
		}
	}

	// Generate job ID
	jobID := uuid.New().String()

	// Check if force reimport
	forceReimport := c.FormValue("force_reimport") == "true"

	// Create async task
	task, err := job.NewImportTransactionsTask(job.ImportTransactionsPayload{
		JobID:         jobID,
		UserID:        userID.String(),
		BankAccountID: bankAccountID.String(),
		BankType:      bankType,
		ForceReimport: forceReimport,
		CSVData:       string(csvData),
	})
	if err != nil {
		return err
	}

	// Enqueue task
	if _, err := h.jobService.Client.Enqueue(task); err != nil {
		return err
	}

	// Return job ID immediately
	return c.JSON(http.StatusAccepted, map[string]interface{}{
		"job_id":          jobID,
		"status":          "processing",
		"message":         "Import job criado. Use o job_id para consultar o progresso.",
		"bank_account_id": bankAccountID.String(),
		"account_created": createAccount,
	})
}

// GetImportStatus returns the status of an import job
func (h *TransactionHandler) GetImportStatus(c echo.Context) error {
	jobID := c.Param("job_id")
	if jobID == "" {
		return errs.NewBadRequestError("job_id is required", false, nil, nil, nil)
	}

	status, err := h.jobService.GetImportStatus(c.Request().Context(), jobID)
	if err != nil {
		return errs.NewNotFoundError("Job não encontrado ou expirado", false, nil)
	}

	return c.JSON(http.StatusOK, status)
}

// DeleteAllByAccount deleta todas as transações de uma conta
func (h *TransactionHandler) DeleteAllByAccount(c echo.Context) error {
	userID := middleware.GetUserID(c)

	accountID, err := uuid.Parse(c.Param("account_id"))
	if err != nil {
		return errs.NewBadRequestError("ID de conta inválido", false, nil, nil, nil)
	}

	if err := h.transactionService.DeleteAllByAccount(c.Request().Context(), userID, accountID); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":    "Todas as transações foram deletadas",
		"account_id": accountID.String(),
	})
}

// BulkDelete deletes multiple transactions within a date range
func (h *TransactionHandler) BulkDelete(c echo.Context) error {
	userID := middleware.GetUserID(c)

	var req service.BulkDeleteRequest
	if err := c.Bind(&req); err != nil {
		return errs.NewBadRequestError("Dados inválidos", false, nil, nil, nil)
	}

	// Validate required fields
	if req.StartDate.IsZero() || req.EndDate.IsZero() {
		return errs.NewBadRequestError("Data inicial e final são obrigatórias", false, nil, nil, nil)
	}

	result, err := h.transactionService.BulkDeleteByDateRange(c.Request().Context(), userID, &req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

func (h *TransactionHandler) validateCSVContent(csvData string) error {
	const maxLines = 50000
	const maxLineLength = 10000

	scanner := bufio.NewScanner(strings.NewReader(csvData))
	lineCount := 0

	for scanner.Scan() {
		lineCount++
		if lineCount > maxLines {
			return errs.NewBadRequestError("CSV excede o limite de 50.000 linhas", false, nil, nil, nil)
		}

		line := scanner.Text()
		if len(line) > maxLineLength {
			return errs.NewBadRequestError("CSV contém linhas muito longas", false, nil, nil, nil)
		}

		if strings.Contains(line, "\x00") {
			return errs.NewBadRequestError("CSV contém caracteres inválidos", false, nil, nil, nil)
		}
	}

	if lineCount == 0 {
		return errs.NewBadRequestError("CSV está vazio", false, nil, nil, nil)
	}

	if err := scanner.Err(); err != nil {
		return errs.NewBadRequestError("erro ao processar CSV", false, nil, nil, nil)
	}

	return nil
}

func (h *TransactionHandler) extractInitialBalance(csvData, bankType string) (string, error) {
	if bankType != "bb" {
		return "0", fmt.Errorf("balance extraction only supported for Banco do Brasil")
	}

	scanner := bufio.NewScanner(strings.NewReader(csvData))
	for scanner.Scan() {
		line := scanner.Text()
		lowerLine := strings.ToLower(line)

		if strings.Contains(lowerLine, "saldo anterior") {
			parts := strings.Split(line, ",")
			if len(parts) >= 2 {
				balanceStr := strings.TrimSpace(parts[len(parts)-1])
				balanceStr = strings.ReplaceAll(balanceStr, "R$", "")
				balanceStr = strings.TrimSpace(balanceStr)
				balanceStr = strings.ReplaceAll(balanceStr, ",", ".")

				if _, err := decimal.NewFromString(balanceStr); err == nil {
					return balanceStr, nil
				}
			}
		}
	}

	return "0", fmt.Errorf("saldo anterior not found in CSV")
}
