package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gustavoz65/finext/internal/errs"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	registerDecimalValidations(validate)
	registerCustomValidations(validate)
}

// registerDecimalValidations registra validações customizadas para o tipo decimal.Decimal
// O validator padrão não suporta decimal.Decimal para tags como gt, gte, lt, lte
func registerDecimalValidations(v *validator.Validate) {
	v.RegisterCustomTypeFunc(func(field reflect.Value) interface{} {
		if valuer, ok := field.Interface().(decimal.Decimal); ok {
			return valuer.InexactFloat64()
		}
		return nil
	}, decimal.Decimal{})
}

type Validatable interface {
	any
}

// BindAndValidate faz o bind do request body e valida os campos
func BindAndValidate(c echo.Context, req interface{}) error {
	if err := c.Bind(req); err != nil {
		return errs.NewBadRequestError("Formato de requisicao invalido", false, nil, nil, nil)
	}

	if err := validate.Struct(req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			fieldErrors := make([]errs.FieldError, 0, len(validationErrors))
			for _, ve := range validationErrors {
				fieldErrors = append(fieldErrors, errs.FieldError{
					Field: toSnakeCase(ve.Field()),
					Error: formatValidationError(ve),
				})
			}
			return errs.NewBadRequestError("Erro de validacao", false, nil, fieldErrors, nil)
		}
		return errs.NewBadRequestError("Erro de validacao: "+err.Error(), false, nil, nil, nil)
	}

	return nil
}

// GetValidator retorna a instancia do validator
func GetValidator() *validator.Validate {
	return validate
}

func formatValidationError(fe validator.FieldError) string {
	field := toSnakeCase(fe.Field())

	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("O campo '%s' e obrigatorio", field)
	case "email":
		return fmt.Sprintf("O campo '%s' deve ser um email valido", field)
	case "min":
		return fmt.Sprintf("O campo '%s' deve ter no minimo %s caracteres", field, fe.Param())
	case "max":
		return fmt.Sprintf("O campo '%s' deve ter no maximo %s caracteres", field, fe.Param())
	case "gt":
		return fmt.Sprintf("O campo '%s' deve ser maior que %s", field, fe.Param())
	case "gte":
		return fmt.Sprintf("O campo '%s' deve ser maior ou igual a %s", field, fe.Param())
	case "lte":
		return fmt.Sprintf("O campo '%s' deve ser menor ou igual a %s", field, fe.Param())
	case "oneof":
		return fmt.Sprintf("O campo '%s' deve ser um dos valores: %s", field, fe.Param())
	case "len":
		return fmt.Sprintf("O campo '%s' deve ter exatamente %s caracteres", field, fe.Param())
	case "hexcolor":
		return fmt.Sprintf("O campo '%s' deve ser uma cor hexadecimal valida", field)
	case "gtfield":
		return fmt.Sprintf("O campo '%s' deve ser maior que '%s'", field, toSnakeCase(fe.Param()))
	case "nefield":
		return fmt.Sprintf("O campo '%s' deve ser diferente de '%s'", field, toSnakeCase(fe.Param()))
	case "maxmoney":
		return ErrMaxMoneyValueExceeded
	case "maxcredit":
		return ErrMaxCreditLimitExceeded
	case "maxbudget":
		return ErrMaxBudgetAmountExceeded
	case "maxgoal":
		return ErrMaxGoalAmountExceeded
	case "minpositive":
		return ErrMinPositiveValueRequired
	default:
		return fmt.Sprintf("O campo '%s' e invalido", field)
	}
}

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				result.WriteByte('_')
			}
			result.WriteRune(r + 32)
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// registerCustomValidations registra validações customizadas para limites do sistema
func registerCustomValidations(v *validator.Validate) {
	// Validação para valores monetários máximos
	v.RegisterValidation("maxmoney", validateMaxMoney)
	v.RegisterValidation("maxcredit", validateMaxCredit)
	v.RegisterValidation("maxbudget", validateMaxBudget)
	v.RegisterValidation("maxgoal", validateMaxGoal)
	v.RegisterValidation("minpositive", validateMinPositive)
}

func validateMaxMoney(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(decimal.Decimal)
	if !ok {
		return false
	}
	return value.LessThanOrEqual(MaxMoneyValue)
}

func validateMaxCredit(fl validator.FieldLevel) bool {
	if fl.Field().IsNil() {
		return true
	}
	value, ok := fl.Field().Interface().(*decimal.Decimal)
	if !ok {
		return false
	}
	return value.LessThanOrEqual(MaxCreditLimit)
}

func validateMaxBudget(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(decimal.Decimal)
	if !ok {
		return false
	}
	return value.LessThanOrEqual(MaxBudgetAmount)
}

func validateMaxGoal(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(decimal.Decimal)
	if !ok {
		return false
	}
	return value.LessThanOrEqual(MaxGoalAmount)
}

func validateMinPositive(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(decimal.Decimal)
	if !ok {
		return false
	}
	return value.GreaterThanOrEqual(MinPositiveValue)
}

// ValidateMoneyValue valida se um valor decimal está dentro dos limites permitidos
func ValidateMoneyValue(value decimal.Decimal) error {
	if value.GreaterThan(MaxMoneyValue) {
		return errs.NewBadRequestError(ErrMaxMoneyValueExceeded, false, nil, nil, nil)
	}
	if value.IsPositive() && value.LessThan(MinPositiveValue) {
		return errs.NewBadRequestError(ErrMinPositiveValueRequired, false, nil, nil, nil)
	}
	return nil
}

// ValidateCreditLimit valida se um limite de crédito está dentro dos limites permitidos
func ValidateCreditLimit(value *decimal.Decimal) error {
	if value != nil && value.GreaterThan(MaxCreditLimit) {
		return errs.NewBadRequestError(ErrMaxCreditLimitExceeded, false, nil, nil, nil)
	}
	return nil
}

// ValidateBudgetAmount valida se um valor de orçamento está dentro dos limites permitidos
func ValidateBudgetAmount(value decimal.Decimal) error {
	if value.GreaterThan(MaxBudgetAmount) {
		return errs.NewBadRequestError(ErrMaxBudgetAmountExceeded, false, nil, nil, nil)
	}
	return nil
}

// ValidateGoalAmount valida se um valor de meta está dentro dos limites permitidos
func ValidateGoalAmount(value decimal.Decimal) error {
	if value.GreaterThan(MaxGoalAmount) {
		return errs.NewBadRequestError(ErrMaxGoalAmountExceeded, false, nil, nil, nil)
	}
	return nil
}
