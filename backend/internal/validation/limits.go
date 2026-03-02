package validation

import "github.com/shopspring/decimal"

// Limites de valores monetários (em formato decimal)
var (
	// Valor máximo permitido para saldos e transações individuais: 10 milhões
	MaxMoneyValue = decimal.NewFromInt(10_000_000)

	// Valor máximo para limite de crédito: 1 milhão
	MaxCreditLimit = decimal.NewFromInt(1_000_000)

	// Valor máximo para orçamentos mensais: 500 mil
	MaxBudgetAmount = decimal.NewFromInt(500_000)

	// Valor máximo para metas financeiras: 50 milhões
	MaxGoalAmount = decimal.NewFromInt(50_000_000)

	// Valor mínimo positivo (para evitar valores muito pequenos/inválidos)
	MinPositiveValue = decimal.NewFromFloat(0.01)
)

// Limites de quantidade de entidades por usuário
const (
	// Número máximo de contas bancárias por usuário
	MaxBankAccountsPerUser = 20

	// Número máximo de transações por dia
	MaxTransactionsPerDay = 100

	// Número máximo de transações recorrentes ativas
	MaxRecurringTransactions = 50

	// Número máximo de orçamentos ativos simultaneamente
	MaxActiveBudgets = 20

	// Número máximo de metas ativas simultaneamente
	MaxActiveGoals = 15

	// Número máximo de categorias customizadas por usuário
	MaxCustomCategories = 50

	// Número máximo de transferências por dia
	MaxTransfersPerDay = 50
)

// Limites de strings e arrays
const (
	// Tamanho máximo de descrição de transação
	MaxTransactionDescriptionLength = 255

	// Tamanho máximo de notas
	MaxNotesLength = 1000

	// Número máximo de tags por transação
	MaxTagsPerTransaction = 10

	// Tamanho máximo de uma tag
	MaxTagLength = 50

	// Número máximo de parcelas permitidas
	MaxInstallments = 120 // 10 anos
)

// Mensagens de erro para limites
const (
	ErrMaxMoneyValueExceeded      = "O valor máximo permitido é R$ 10.000.000,00"
	ErrMaxCreditLimitExceeded     = "O limite de crédito máximo permitido é R$ 1.000.000,00"
	ErrMaxBudgetAmountExceeded    = "O valor máximo para orçamento é R$ 500.000,00"
	ErrMaxGoalAmountExceeded      = "O valor máximo para meta é R$ 50.000.000,00"
	ErrMaxBankAccountsExceeded    = "Você atingiu o limite máximo de 20 contas bancárias"
	ErrMaxTransactionsPerDay      = "Você atingiu o limite de 100 transações por dia"
	ErrMaxRecurringTransactions   = "Você atingiu o limite de 50 transações recorrentes ativas"
	ErrMaxActiveBudgetsExceeded   = "Você atingiu o limite de 20 orçamentos ativos"
	ErrMaxActiveGoalsExceeded     = "Você atingiu o limite de 15 metas ativas"
	ErrMaxTransfersPerDay         = "Você atingiu o limite de 50 transferências por dia"
	ErrMinPositiveValueRequired   = "O valor mínimo permitido é R$ 0,01"
)
