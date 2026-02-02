package model

import (
	"time"

	"github.com/google/uuid"
)

type CategoryType string

const (
	CategoryTypeIncome  CategoryType = "income"
	CategoryTypeExpense CategoryType = "expense"
)

type Category struct {
	ID          uuid.UUID    `json:"id" db:"id"`
	UserID      *uuid.UUID   `json:"user_id,omitempty" db:"user_id"`
	Name        string       `json:"name" db:"name"`
	Description *string      `json:"description,omitempty" db:"description"`
	Type        CategoryType `json:"type" db:"type"`
	Color       string       `json:"color" db:"color"`
	Icon        string       `json:"icon" db:"icon"`
	IsSystem    bool         `json:"is_system" db:"is_system"`
	IsActive    bool         `json:"is_active" db:"is_active"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
}

func (c *Category) IsUserCategory() bool {
	return c.UserID != nil
}

func (c *Category) CanBeModified() bool {
	return !c.IsSystem
}

func (c *Category) CanBeDeleted() bool {
	return !c.IsSystem && c.UserID != nil
}
