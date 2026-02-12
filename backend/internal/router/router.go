package router

import (
	"github.com/gustavoz65/Cashing-go/internal/config"
	"github.com/gustavoz65/Cashing-go/internal/database"
	"github.com/gustavoz65/Cashing-go/internal/handler"
	"github.com/gustavoz65/Cashing-go/internal/lib/utils/validator"
	"github.com/gustavoz65/Cashing-go/internal/middleware"
	"github.com/gustavoz65/Cashing-go/internal/repository"
	"github.com/gustavoz65/Cashing-go/internal/server"
	"github.com/gustavoz65/Cashing-go/internal/service"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

// New cria e configura o router Echo com todas as rotas
func New(cfg *config.Config, db *database.Database, logger *zerolog.Logger, srv *server.Server) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Validator = validator.New()

	// Error handler global
	e.HTTPErrorHandler = middleware.ErrorHandler(logger)

	// Middlewares globais
	e.Use(echomiddleware.RequestID())
	e.Use(middleware.RecoveryMiddleware(logger))
	e.Use(middleware.LoggerMiddleware(logger))
	e.Use(middleware.CORSMiddleware(cfg.Server.CORSAllowedOrigins))

	// Repositories
	userRepo := repository.NewUserRepository(db, logger)
	transactionRepo := repository.NewTransactionRepository(db, logger)
	accountRepo := repository.NewBankAccountRepository(db, logger)
	categoryRepo := repository.NewCategoryRepository(db, logger)
	budgetRepo := repository.NewBudgetRepository(db, logger)
	goalRepo := repository.NewGoalRepository(db, logger)
	notificationRepo := repository.NewNotificationRepository(db, logger)
	auditLogRepo := repository.NewAuditLogRepository(db, logger)
	recurringRepo := repository.NewRecurringTransactionRepository(db, logger)

	// Services
	authService := service.NewAuthService(userRepo, cfg, logger)
	userService := service.NewUserService(userRepo, logger)
	transactionService := service.NewTransactionService(transactionRepo, accountRepo, budgetRepo, userRepo, logger)
	accountService := service.NewBankAccountService(accountRepo, logger)
	categoryService := service.NewCategoryService(categoryRepo, logger)
	budgetService := service.NewBudgetService(budgetRepo, notificationRepo, logger)
	goalService := service.NewGoalService(goalRepo, notificationRepo, logger)
	dashboardService := service.NewDashboardService(accountRepo, transactionRepo, budgetRepo, goalRepo, logger)
	notificationService := service.NewNotificationService(notificationRepo, userRepo, logger)
	recurringService := service.NewRecurringTransactionService(recurringRepo, transactionRepo, accountRepo, logger)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	accountHandler := handler.NewBankAccountHandler(accountService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	budgetHandler := handler.NewBudgetHandler(budgetService)
	goalHandler := handler.NewGoalHandler(goalService)
	dashboardHandler := handler.NewDashboardHandler(dashboardService)
	notificationHandler := handler.NewNotificationHandler(notificationService)
	recurringHandler := handler.NewRecurringTransactionHandler(recurringService)

	// Health check
	healthHandler := handler.NewHealthHandler(srv)
	e.GET("/health", healthHandler.CheckHandler)

	// API v1
	api := e.Group("/api/v1")

	//  Rotas publicas (sem autenticacao)
	authRateLimiter := middleware.AuthRateLimit(srv.Redis)

	auth := api.Group("/auth")
	authWithRL := api.Group("/auth", authRateLimiter)

	authWithRL.POST("/register", authHandler.Register)
	authWithRL.POST("/login", authHandler.Login)
	auth.POST("/refresh", authHandler.RefreshToken)
	auth.POST("/logout", authHandler.Logout)

	// --- Rotas autenticadas ---
	authMiddleware := middleware.AuthMiddleware(authService)

	// Audit middleware (deve vir DEPOIS do AuthMiddleware para ter acesso ao user_id)
	auditMiddleware := middleware.NewAuditMiddleware(auditLogRepo, logger)

	// Auth (requer autenticacao)
	authProtected := api.Group("/auth", authMiddleware, auditMiddleware.Handler())
	authProtected.POST("/change-password", authHandler.ChangePassword)

	// Users
	users := api.Group("/users", authMiddleware, auditMiddleware.Handler())
	users.GET("/me", userHandler.GetMe)
	users.PUT("/me", userHandler.UpdateMe)
	users.DELETE("/me", userHandler.DeactivateMe)
	users.GET("/settings", userHandler.GetSettings)
	users.PUT("/settings", userHandler.UpdateSettings)

	// Categories
	categories := api.Group("/categories", authMiddleware, auditMiddleware.Handler())
	categories.GET("", categoryHandler.GetAll)
	categories.GET("/:id", categoryHandler.GetByID)
	categories.POST("", categoryHandler.Create)
	categories.PUT("/:id", categoryHandler.Update)
	categories.DELETE("/:id", categoryHandler.Delete)

	// Bank Accounts
	accounts := api.Group("/accounts", authMiddleware, auditMiddleware.Handler())
	accounts.GET("", accountHandler.GetAll)
	accounts.GET("/balance", accountHandler.GetTotalBalance)
	accounts.GET("/:id", accountHandler.GetByID)
	accounts.POST("", accountHandler.Create)
	accounts.PUT("/:id", accountHandler.Update)
	accounts.DELETE("/:id", accountHandler.Delete)

	// Transactions
	transactions := api.Group("/transactions", authMiddleware, auditMiddleware.Handler())
	transactions.GET("", transactionHandler.GetAll)
	transactions.GET("/upcoming", transactionHandler.GetUpcoming)
	transactions.GET("/:id", transactionHandler.GetByID)
	transactions.POST("", transactionHandler.Create)
	transactions.PUT("/:id", transactionHandler.Update)
	transactions.DELETE("/:id", transactionHandler.Delete)
	transactions.PATCH("/:id/pay", transactionHandler.MarkAsPaid)

	// Recurring Transactions
	recurring := api.Group("/recurring-transactions", authMiddleware, auditMiddleware.Handler())
	recurring.GET("", recurringHandler.GetAll)
	recurring.GET("/stats", recurringHandler.GetStats)
	recurring.GET("/upcoming", recurringHandler.GetUpcoming)
	recurring.GET("/:id", recurringHandler.GetByID)
	recurring.POST("", recurringHandler.Create)
	recurring.PUT("/:id", recurringHandler.Update)
	recurring.DELETE("/:id", recurringHandler.Delete)
	recurring.PATCH("/:id/toggle", recurringHandler.ToggleActive)

	// Budgets
	budgets := api.Group("/budgets", authMiddleware, auditMiddleware.Handler())
	budgets.GET("", budgetHandler.GetAll)
	budgets.GET("/summary", budgetHandler.GetSummary)
	budgets.GET("/:id", budgetHandler.GetByID)
	budgets.POST("", budgetHandler.Create)
	budgets.PUT("/:id", budgetHandler.Update)
	budgets.DELETE("/:id", budgetHandler.Delete)

	// Goals
	goals := api.Group("/goals", authMiddleware, auditMiddleware.Handler())
	goals.GET("", goalHandler.GetAll)
	goals.GET("/summary", goalHandler.GetSummary)
	goals.GET("/:id", goalHandler.GetByID)
	goals.POST("", goalHandler.Create)
	goals.PUT("/:id", goalHandler.Update)
	goals.DELETE("/:id", goalHandler.Delete)
	goals.POST("/:id/contributions", goalHandler.AddContribution)
	goals.GET("/:id/contributions", goalHandler.GetContributions)

	// Dashboard
	dashboard := api.Group("/dashboard", authMiddleware, auditMiddleware.Handler())
	dashboard.GET("", dashboardHandler.GetSummary)
	dashboard.GET("/cash-flow", dashboardHandler.GetCashFlow)
	dashboard.GET("/income-expense", dashboardHandler.GetIncomeVsExpense)
	dashboard.GET("/monthly-comparison", dashboardHandler.GetMonthlyComparison)
	dashboard.GET("/account-balances", dashboardHandler.GetAccountBalances)

	// Notifications
	notifications := api.Group("/notifications", authMiddleware, auditMiddleware.Handler())
	notifications.GET("", notificationHandler.GetAll)
	notifications.GET("/unread", notificationHandler.GetUnread)
	notifications.GET("/unread/count", notificationHandler.GetUnreadCount)
	notifications.PATCH("/:id/read", notificationHandler.MarkAsRead)
	notifications.PATCH("/read-all", notificationHandler.MarkAllAsRead)
	notifications.DELETE("/:id", notificationHandler.Delete)

	srv.Job.SetRecurringService(recurringService)

	return e
}
