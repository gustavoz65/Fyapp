package model

import (
	"time"

	"github.com/google/uuid"
)

type CategoryPattern struct {
	ID         uuid.UUID `json:"id" db:"id"`
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	Keyword    string    `json:"keyword" db:"keyword"`
	CategoryID uuid.UUID `json:"category_id" db:"category_id"`
	Confidence int       `json:"confidence" db:"confidence"`
	Source     string    `json:"source" db:"source"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`

	Category *Category `json:"category,omitempty" db:"-"`
}

type CategorySuggestion struct {
	CategoryID   uuid.UUID `json:"category_id"`
	CategoryName string    `json:"category_name"`
	Confidence   int       `json:"confidence"`
	MatchCount   int       `json:"match_count"`
}
