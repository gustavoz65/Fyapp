package categorization

import (
	"testing"
)

func TestLoadRules(t *testing.T) {
	rules, err := LoadRules("rules.json")
	if err != nil {
		t.Fatalf("LoadRules() error = %v", err)
	}

	if len(rules.ExpensePatterns) == 0 {
		t.Error("Expected expense patterns, got none")
	}

	if len(rules.RegexPatterns) == 0 {
		t.Error("Expected regex patterns, got none")
	}

	// Check first expense pattern
	first := rules.ExpensePatterns[0]
	if len(first.Keywords) == 0 {
		t.Error("Expected keywords in first pattern")
	}
	if first.Category == "" {
		t.Error("Expected category in first pattern")
	}
}
