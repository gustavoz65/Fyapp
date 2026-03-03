package router

import (
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"

	"github.com/gustavoz65/Fyapp/internal/config"
	"github.com/gustavoz65/Fyapp/internal/database"
	"github.com/gustavoz65/Fyapp/internal/handler"
	"github.com/gustavoz65/Fyapp/internal/lib/utils/validator"
	"github.com/gustavoz65/Fyapp/internal/middleware"
	"github.com/gustavoz65/Fyapp/internal/repository"
	"github.com/gustavoz65/Fyapp/internal/server"
	"github.com/gustavoz65/Fyapp/internal/service"
)

// New cria e configura o router Echo com todas as rotas
func New(cfg *config.Config, db *database.Database, logger *zerolog.Logger, srv *server.Server) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Validator = validator.New()

	// Error handler global
	e.HTTPErrorHandler = middleware.ErrorHandler(logger)

	e.Use(echomiddleware.RequestID())
	e.Use(middleware.RecoveryMiddleware(logger))
	e.Use(middleware.LoggerMiddleware(logger))
	e.Use(middleware.SecurityHeadersMiddleware())
	e.Use(middleware.CORSMiddleware(cfg.Server.CORSAllowedOrigins))
	e.Use(middleware.CSRFTokenGenerator())

	// Repositories
	userRepo := repository.NewUserRepository(db, logger)
	providerRepo := repository.NewOAuthProviderRepository(db, logger)
	transactionRepo := repository.NewTransactionRepository(db, logger)
	accountRepo := repository.NewBankAccountRepository(db, logger)
	categoryRepo := repository.NewCategoryRepository(db, logger)
	categoryPatternRepo := repository.NewCategoryPatternRepository(db, logger)
	budgetRepo := repository.NewBudgetRepository(db, logger)
	goalRepo := repository.NewGoalRepository(db, logger)
	notificationRepo := repository.NewNotificationRepository(db, logger)
	auditLogRepo := repository.NewAuditLogRepository(db, logger)
	recurringRepo := repository.NewRecurringTransactionRepository(db, logger)

	// Services
	authService := service.NewAuthService(userRepo, providerRepo, srv.FirebaseClient, cfg, logger)
	userService := service.NewUserService(userRepo, logger)
	categorizationService := service.NewCategorizationService(categoryPatternRepo, logger)
	transactionService := service.NewTransactionService(transactionRepo, accountRepo, budgetRepo, userRepo, categorizationService, logger)
	accountService := service.NewBankAccountService(accountRepo, transactionRepo, logger)
	categoryService := service.NewCategoryService(categoryRepo, logger)
	notificationService := service.NewNotificationService(notificationRepo, userRepo, logger, srv.Job.Client)
	budgetService := service.NewBudgetService(budgetRepo, notificationService, logger)
	goalService := service.NewGoalService(goalRepo, notificationService, logger)
	dashboardService := service.NewDashboardService(accountRepo, transactionRepo, budgetRepo, goalRepo, logger)
	recurringService := service.NewRecurringTransactionService(recurringRepo, transactionRepo, accountRepo, logger)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	transactionHandler := handler.NewTransactionHandler(transactionService, accountService, srv.Job)
	accountHandler := handler.NewBankAccountHandler(accountService)
	categoryHandler := handler.NewCategoryHandler(categoryService, categorizationService)
	budgetHandler := handler.NewBudgetHandler(budgetService)
	goalHandler := handler.NewGoalHandler(goalService)
	dashboardHandler := handler.NewDashboardHandler(dashboardService)
	notificationHandler := handler.NewNotificationHandler(notificationService)
	recurringHandler := handler.NewRecurringTransactionHandler(recurringService)
	wsHandler := handler.NewWebSocketHandler(srv.Job, logger)

	// Health check
	healthHandler := handler.NewHealthHandler(srv)
	e.GET("/health", healthHandler.CheckHandler)

	api := e.Group("/api/v1")

	authRateLimiter := middleware.AuthRateLimit(srv.Redis)
	refreshRateLimiter := middleware.RefreshRateLimit(srv.Redis)
	csrfTokenGen := middleware.CSRFTokenGenerator()

	authWithRL := api.Group("/auth", authRateLimiter, csrfTokenGen)
	authRefresh := api.Group("/auth", refreshRateLimiter, csrfTokenGen)

	authWithRL.POST("/register", authHandler.Register)
	authWithRL.POST("/login", authHandler.Login)
	authWithRL.POST("/social/login", authHandler.SocialLogin)
	authRefresh.POST("/refresh", authHandler.RefreshToken)
	authRefresh.POST("/logout", authHandler.Logout)

	authMiddleware := middleware.AuthMiddleware(authService)
	auditMiddleware := middleware.NewAuditMiddleware(auditLogRepo, logger)
	csrfMiddleware := middleware.CSRFMiddleware()
	mutationRL := middleware.MutationRateLimit(srv.Redis)
	readRL := middleware.ReadRateLimit(srv.Redis)
	uploadRL := middleware.UploadRateLimit(srv.Redis)

	// WebSocket routes (with auth via query parameter)
	wsAuthMiddleware := middleware.WebSocketAuthMiddleware(authService)
	ws := e.Group("/ws", wsAuthMiddleware)
	ws.GET("/import-progress", wsHandler.ImportProgress)

	authProtected := api.Group("/auth", authMiddleware, auditMiddleware.Handler(), csrfMiddleware)
	authProtected.POST("/change-password", authHandler.ChangePassword, mutationRL)
	authProtected.POST("/set-password", authHandler.SetPassword, mutationRL)
	authProtected.POST("/social/link", authHandler.LinkProvider, mutationRL)
	authProtected.DELETE("/social/:provider", authHandler.UnlinkProvider, mutationRL)
	authProtected.GET("/social/providers", authHandler.GetLinkedProviders, readRL)

	users := api.Group("/users", authMiddleware, auditMiddleware.Handler(), csrfMiddleware)
	users.GET("/me", userHandler.GetMe, readRL)
	users.PUT("/me", userHandler.UpdateMe, mutationRL)
	users.DELETE("/me", userHandler.DeactivateMe, mutationRL)
	users.GET("/settings", userHandler.GetSettings, readRL)
	users.PUT("/settings", userHandler.UpdateSettings, mutationRL)

	categories := api.Group("/categories", authMiddleware, auditMiddleware.Handler(), csrfMiddleware)
	categories.GET("", categoryHandler.GetAll, readRL)
	categories.GET("/suggest", categoryHandler.SuggestCategory, readRL)
	categories.GET("/:id", categoryHandler.GetByID, readRL)
	categories.POST("", categoryHandler.Create, mutationRL)
	categories.PUT("/:id", categoryHandler.Update, mutationRL)
	categories.DELETE("/:id", categoryHandler.Delete, mutationRL)

	accounts := api.Group("/accounts", authMiddleware, auditMiddleware.Handler(), csrfMiddleware)
	accounts.GET("", accountHandler.GetAll, readRL)
	accounts.GET("/balance", accountHandler.GetTotalBalance, readRL)
	accounts.GET("/:id", accountHandler.GetByID, readRL)
	accounts.POST("", accountHandler.Create, mutationRL)
	accounts.POST("/:id/recalculate", accountHandler.RecalculateBalance, mutationRL)
	accounts.PUT("/:id", accountHandler.Update, mutationRL)
	accounts.DELETE("/:id", accountHandler.Delete, mutationRL)

	transactions := api.Group("/transactions", authMiddleware, auditMiddleware.Handler(), csrfMiddleware)
	transactions.GET("", transactionHandler.GetAll, readRL)
	transactions.GET("/upcoming", transactionHandler.GetUpcoming, readRL)
	transactions.GET("/:id", transactionHandler.GetByID, readRL)
	transactions.POST("", transactionHandler.Create, mutationRL)
	transactions.POST("/import", transactionHandler.Import, mutationRL, uploadRL)
	transactions.GET("/import/:job_id", transactionHandler.GetImportStatus, readRL)
	transactions.PUT("/:id", transactionHandler.Update, mutationRL)
	transactions.DELETE("/:id", transactionHandler.Delete, mutationRL)
	transactions.DELETE("/account/:account_id", transactionHandler.DeleteAllByAccount, mutationRL)
	transactions.PATCH("/:id/pay", transactionHandler.MarkAsPaid, mutationRL)

	recurring := api.Group("/recurring-transactions", authMiddleware, auditMiddleware.Handler(), csrfMiddleware)
	recurring.GET("", recurringHandler.GetAll, readRL)
	recurring.GET("/stats", recurringHandler.GetStats, readRL)
	recurring.GET("/upcoming", recurringHandler.GetUpcoming, readRL)
	recurring.GET("/:id", recurringHandler.GetByID, readRL)
	recurring.POST("", recurringHandler.Create, mutationRL)
	recurring.PUT("/:id", recurringHandler.Update, mutationRL)
	recurring.DELETE("/:id", recurringHandler.Delete, mutationRL)
	recurring.PATCH("/:id/toggle", recurringHandler.ToggleActive, mutationRL)

	budgets := api.Group("/budgets", authMiddleware, auditMiddleware.Handler(), csrfMiddleware)
	budgets.GET("", budgetHandler.GetAll, readRL)
	budgets.GET("/summary", budgetHandler.GetSummary, readRL)
	budgets.GET("/:id", budgetHandler.GetByID, readRL)
	budgets.POST("", budgetHandler.Create, mutationRL)
	budgets.PUT("/:id", budgetHandler.Update, mutationRL)
	budgets.DELETE("/:id", budgetHandler.Delete, mutationRL)

	goals := api.Group("/goals", authMiddleware, auditMiddleware.Handler(), csrfMiddleware)
	goals.GET("", goalHandler.GetAll, readRL)
	goals.GET("/summary", goalHandler.GetSummary, readRL)
	goals.GET("/:id", goalHandler.GetByID, readRL)
	goals.POST("", goalHandler.Create, mutationRL)
	goals.PUT("/:id", goalHandler.Update, mutationRL)
	goals.DELETE("/:id", goalHandler.Delete, mutationRL)
	goals.POST("/:id/contributions", goalHandler.AddContribution, mutationRL)
	goals.GET("/:id/contributions", goalHandler.GetContributions, readRL)

	dashboard := api.Group("/dashboard", authMiddleware, auditMiddleware.Handler(), csrfMiddleware)
	dashboard.GET("", dashboardHandler.GetSummary, readRL)
	dashboard.GET("/cash-flow", dashboardHandler.GetCashFlow, readRL)
	dashboard.GET("/income-expense", dashboardHandler.GetIncomeVsExpense, readRL)
	dashboard.GET("/monthly-comparison", dashboardHandler.GetMonthlyComparison, readRL)
	dashboard.GET("/account-balances", dashboardHandler.GetAccountBalances, readRL)

	notifications := api.Group("/notifications", authMiddleware, auditMiddleware.Handler(), csrfMiddleware)
	notifications.GET("", notificationHandler.GetAll, readRL)
	notifications.GET("/unread", notificationHandler.GetUnread, readRL)
	notifications.GET("/unread/count", notificationHandler.GetUnreadCount, readRL)
	notifications.PATCH("/:id/read", notificationHandler.MarkAsRead, mutationRL)
	notifications.PATCH("/read-all", notificationHandler.MarkAllAsRead, mutationRL)
	notifications.DELETE("/:id", notificationHandler.Delete, mutationRL)

	srv.Job.SetRecurringService(recurringService)
	srv.Job.SetTransactionService(transactionService)

	return e
}
