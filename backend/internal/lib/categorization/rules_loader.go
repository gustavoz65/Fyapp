package categorization

import (
	"embed"
	"encoding/json"
	"fmt"
)

//go:embed rules.json
var rulesFS embed.FS

type ExpensePattern struct {
	Keywords   []string `json:"keywords"`
	Category   string   `json:"category"`
	Confidence int      `json:"confidence"`
}

type IncomePattern struct {
	Keywords   []string `json:"keywords"`
	Category   string   `json:"category"`
	Confidence int      `json:"confidence"`
}

type RegexPattern struct {
	Pattern  string `json:"pattern"`
	Category string `json:"category"`
	Type     string `json:"type"`
}

type ExclusionPattern struct {
	Keywords          []string `json:"keywords"`
	ExcludeCategories []string `json:"exclude_categories"`
	Reason            string   `json:"reason"`
}

type Rules struct {
	ExpensePatterns   []ExpensePattern   `json:"expense_patterns"`
	IncomePatterns    []IncomePattern    `json:"income_patterns"`
	RegexPatterns     []RegexPattern     `json:"regex_patterns"`
	ExclusionPatterns []ExclusionPattern `json:"exclusion_patterns"`
}

var cachedRules *Rules

func LoadRules(filename string) (*Rules, error) {
	if cachedRules != nil {
		return cachedRules, nil
	}

	data, err := rulesFS.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read rules file: %w", err)
	}

	var rules Rules
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, fmt.Errorf("failed to parse rules JSON: %w", err)
	}

	cachedRules = &rules
	return cachedRules, nil
}
