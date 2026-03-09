# CSV Parser Bug Fixes Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Fix critical bugs in CSV parser causing data corruption (encoding, duplicates, invalid dates, wrong categorization)

**Architecture:** Add encoding detection layer, enhance line filtering before parse, improve categorization with exclusion rules

**Tech Stack:** Go 1.26, golang.org/x/text/encoding, shopspring/decimal, existing parser infrastructure

---

## Task 1: Create Encoding Detector Module

**Files:**
- Create: `backend/internal/lib/encoding/detector.go`
- Create: `backend/internal/lib/encoding/detector_test.go`

**Step 1: Write the failing test**

Create `backend/internal/lib/encoding/detector_test.go`:

```go
package encoding

import (
	"bytes"
	"testing"
)

func TestDetectAndConvert_UTF8(t *testing.T) {
	input := []byte("Data,Valor,Descrição\n01/01/2026,100.00,Teste")
	result, err := DetectAndConvert(bytes.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := result.String()
	if !bytes.Contains([]byte(output), []byte("Descrição")) {
		t.Errorf("expected UTF-8 preserved, got: %s", output)
	}
}

func TestDetectAndConvert_ISO88591(t *testing.T) {
	// ISO-8859-1 encoded "Lançamento"
	input := []byte{0x4C, 0x61, 0x6E, 0xE7, 0x61, 0x6D, 0x65, 0x6E, 0x74, 0x6F}
	result, err := DetectAndConvert(bytes.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := result.String()
	if !bytes.Contains([]byte(output), []byte("Lançamento")) {
		t.Errorf("expected conversion to UTF-8, got: %s", output)
	}
}

func TestDetectAndConvert_Windows1252(t *testing.T) {
	// Windows-1252 encoded "Cartão"
	input := []byte{0x43, 0x61, 0x72, 0x74, 0xE3, 0x6F}
	result, err := DetectAndConvert(bytes.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := result.String()
	if !bytes.Contains([]byte(output), []byte("Cartão")) {
		t.Errorf("expected conversion to UTF-8, got: %s", output)
	}
}
```

**Step 2: Run test to verify it fails**

```bash
cd backend
go test ./internal/lib/encoding/... -v
```

Expected: FAIL with "no such file or directory" or "undefined: DetectAndConvert"

**Step 3: Write minimal implementation**

Create `backend/internal/lib/encoding/detector.go`:

```go
package encoding

import (
	"bytes"
	"io"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// DetectAndConvert attempts to detect the encoding of the input and convert to UTF-8
func DetectAndConvert(r io.Reader) (*bytes.Buffer, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// Try UTF-8 first (most common, no conversion needed)
	if isUTF8(data) {
		return bytes.NewBuffer(data), nil
	}

	// Try common encodings
	encodings := []encoding.Encoding{
		charmap.ISO8859_1,  // Latin-1
		charmap.Windows1252, // Windows Portuguese
	}

	for _, enc := range encodings {
		decoder := enc.NewDecoder()
		converted, _, err := transform.Bytes(decoder, data)
		if err == nil && isUTF8(converted) {
			return bytes.NewBuffer(converted), nil
		}
	}

	// Fallback: return as-is
	return bytes.NewBuffer(data), nil
}

func isUTF8(data []byte) bool {
	// Check if valid UTF-8
	decoder := unicode.UTF8.NewDecoder()
	_, _, err := transform.Bytes(decoder, data)
	return err == nil
}
```

**Step 4: Run test to verify it passes**

```bash
cd backend
go test ./internal/lib/encoding/... -v
```

Expected: PASS (3 tests)

**Step 5: Commit**

```bash
git add backend/internal/lib/encoding/
git commit -m "feat: add CSV encoding detector for UTF-8/ISO-8859-1/Windows-1252

Detects and converts CSV encodings to UTF-8 automatically.
Fixes BB CSV encoding corruption (Lan�amento → Lançamento).

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Task 2: Add Control Line Filter to Parser

**Files:**
- Modify: `backend/internal/lib/parser/transaction_parser.go`
- Create: `backend/internal/lib/parser/filters_test.go`

**Step 1: Write the failing test**

Create `backend/internal/lib/parser/filters_test.go`:

```go
package parser

import "testing"

func TestIsControlLine(t *testing.T) {
	tests := []struct {
		name     string
		record   []string
		expected bool
	}{
		{
			name:     "saldo anterior",
			record:   []string{"24/02/2026", "Saldo Anterior", "", "", "-439,69", ""},
			expected: true,
		},
		{
			name:     "saldo do dia",
			record:   []string{"00/00/0000", "Saldo do dia", "", "", "-490,87", ""},
			expected: true,
		},
		{
			name:     "S A L D O",
			record:   []string{"05/03/2026", "S A L D O", "", "", "-490,87", ""},
			expected: true,
		},
		{
			name:     "empty tipo lancamento (BB)",
			record:   []string{"05/03/2026", "Pix - Enviado", "Details", "123", "10.00", ""},
			expected: true,
		},
		{
			name:     "valid transaction",
			record:   []string{"02/03/2026", "Pix - Enviado", "Details", "123", "10.00", "Saída"},
			expected: false,
		},
	}

	p := NewTransactionParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := p.isControlLine(tt.record, "bb")
			if result != tt.expected {
				t.Errorf("expected %v, got %v for record: %v", tt.expected, result, tt.record)
			}
		})
	}
}

func TestIsBBDuplicatePayment(t *testing.T) {
	tests := []struct {
		name     string
		record   []string
		expected bool
	}{
		{
			name:     "pagamento pix cartao credito",
			record:   []string{"02/03/2026", "Pagamento Pix Cartão Crédito", "Details", "123", "10.00", "Entrada"},
			expected: true,
		},
		{
			name:     "pagamento pix cartao cred (partial)",
			record:   []string{"02/03/2026", "Pagamento Pix Cartão Créd", "Details", "123", "10.00", "Entrada"},
			expected: true,
		},
		{
			name:     "pix enviado (not duplicate)",
			record:   []string{"02/03/2026", "Pix - Enviado", "Details", "123", "10.00", "Saída"},
			expected: false,
		},
	}

	p := NewTransactionParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := p.isBBDuplicatePayment(tt.record)
			if result != tt.expected {
				t.Errorf("expected %v, got %v for record: %v", tt.expected, result, tt.record)
			}
		})
	}
}
```

**Step 2: Run test to verify it fails**

```bash
cd backend
go test ./internal/lib/parser/... -v -run "TestIsControlLine|TestIsBBDuplicatePayment"
```

Expected: FAIL with "undefined: TransactionParser.isControlLine"

**Step 3: Add filter methods to parser**

In `backend/internal/lib/parser/transaction_parser.go`, add after `isHeaderRow` method:

```go
// isControlLine checks if a CSV record is a control/summary line (not a transaction)
func (p *TransactionParser) isControlLine(record []string, bankType string) bool {
	if len(record) < 2 {
		return false
	}

	description := strings.ToLower(strings.TrimSpace(record[1]))

	// Check for common control line keywords
	controlKeywords := []string{
		"saldo anterior",
		"saldo do dia",
		"s a l d o",
		"saldo devedor",
		"total",
	}

	for _, keyword := range controlKeywords {
		if strings.Contains(description, keyword) {
			return true
		}
	}

	// BB specific: check if "Tipo Lançamento" column (index 5) is empty
	if strings.ToLower(bankType) == "bb" {
		if len(record) > 5 && strings.TrimSpace(record[5]) == "" {
			return true
		}
	}

	return false
}

// isBBDuplicatePayment checks if record is a BB "Pagamento Pix Cartão Crédito" duplicate
func (p *TransactionParser) isBBDuplicatePayment(record []string) bool {
	if len(record) < 2 {
		return false
	}

	lancamento := strings.ToLower(strings.TrimSpace(record[1]))

	// These are accounting entries, not real transactions
	// The actual debit appears separately as "Pix - Enviado"
	if strings.Contains(lancamento, "pagamento pix cart") &&
	   strings.Contains(lancamento, "cr") {
		return true
	}

	return false
}
```

**Step 4: Integrate filters into ParseCSV method**

In `transaction_parser.go`, modify the `ParseCSV` loop (around line 100-115):

```go
for {
	record, err := csvReader.Read()
	if err == io.EOF {
		break
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao ler a linha %d do CSV: %w", lineNumber, err)
	}

	lineNumber++

	// Skip header row
	if lineNumber == 1 {
		if p.isHeaderRow(record) {
			continue
		}
	}

	// Skip empty lines
	if len(record) == 0 || (len(record) == 1 && record[0] == "") {
		continue
	}

	// NEW: Skip control lines
	if p.isControlLine(record, bankType) {
		continue
	}

	// NEW: Skip BB duplicate payments
	if strings.ToLower(bankType) == "bb" && p.isBBDuplicatePayment(record) {
		continue
	}

	transaction, err := p.parseCSVRecord(record, mapping, bankType)
	if err != nil {
		fmt.Printf("Aviso: pulando a linha %d: %v\n", lineNumber, err)
		continue
	}

	if !seenTransactions[transaction.ExternalID] {
		transactions = append(transactions, transaction)
		seenTransactions[transaction.ExternalID] = true
	}
}
```

**Step 5: Run tests to verify they pass**

```bash
cd backend
go test ./internal/lib/parser/... -v
```

Expected: PASS (all tests including new filter tests)

**Step 6: Commit**

```bash
git add backend/internal/lib/parser/
git commit -m "feat: add control line and duplicate payment filters

Filters out:
- Control lines (Saldo Anterior, Saldo do dia, S A L D O)
- BB duplicate 'Pagamento Pix Cartão Crédito' entries
- Empty 'Tipo Lançamento' rows

Fixes dashboard showing invalid transactions.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Task 3: Integrate Encoding Detector into Parser

**Files:**
- Modify: `backend/internal/lib/parser/transaction_parser.go`
- Create: `backend/internal/lib/parser/integration_test.go`

**Step 1: Write integration test**

Create `backend/internal/lib/parser/integration_test.go`:

```go
package parser

import (
	"bytes"
	"strings"
	"testing"
)

func TestParseCSV_WithEncoding_BB(t *testing.T) {
	// Simulated BB CSV with ISO-8859-1 encoding
	csvData := `"Data","Lançamento","Detalhes","Nº documento","Valor","Tipo Lançamento"
"02/03/2026","Pix - Enviado","28/02 23:24 Kauê Rodrigues","30202","-10,00","Saída"
"02/03/2026","Pagamento Pix Cartão Crédito","28/02 23:24 Kauê Rodrigues","100611","10,00","Entrada"
"00/00/0000","Saldo do dia","","","-490,87",""
"05/03/2026","S A L D O","","","-490,87",""`

	parser := NewTransactionParser()
	transactions, err := parser.ParseCSV(strings.NewReader(csvData), "bb")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should only have 1 valid transaction (Pix - Enviado)
	// Filters out: Pagamento Pix Cartão, Saldo do dia, S A L D O
	if len(transactions) != 1 {
		t.Errorf("expected 1 transaction, got %d", len(transactions))
	}

	if len(transactions) > 0 {
		if !strings.Contains(transactions[0].Description, "Pix - Enviado") {
			t.Errorf("expected Pix - Enviado, got: %s", transactions[0].Description)
		}
	}
}
```

**Step 2: Run test to verify current behavior**

```bash
cd backend
go test ./internal/lib/parser/... -v -run TestParseCSV_WithEncoding
```

Expected: May PASS or FAIL depending on current state (establishes baseline)

**Step 3: Add encoding detection to ParseCSV**

In `transaction_parser.go`, modify `ParseCSV` method (around line 75-82):

```go
func (p *TransactionParser) ParseCSV(reader io.Reader, bankType string) ([]TransactionImport, error) {
	mapping, ok := BankMappings[strings.ToLower(bankType)]
	if !ok {
		mapping = BankMappings["generic"]
	}

	// NEW: Detect and convert encoding to UTF-8
	convertedReader, err := encoding.DetectAndConvert(reader)
	if err != nil {
		return nil, fmt.Errorf("erro ao detectar encoding: %w", err)
	}

	csvReader := csv.NewReader(convertedReader)
	csvReader.FieldsPerRecord = -1
	csvReader.TrimLeadingSpace = true
	// ... rest of method
```

Add import at top of file:

```go
import (
	"crypto/md5"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gustavoz65/Fyapp/internal/lib/encoding"  // NEW
	"github.com/shopspring/decimal"
)
```

**Step 4: Run integration test**

```bash
cd backend
go test ./internal/lib/parser/... -v
```

Expected: PASS (all tests)

**Step 5: Commit**

```bash
git add backend/internal/lib/parser/transaction_parser.go backend/internal/lib/parser/integration_test.go
git commit -m "feat: integrate encoding detection into CSV parser

Automatically detects and converts CSV encoding to UTF-8.
Fixes corrupted characters in BB imports.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Task 4: Improve Categorization - Add Exclusion Rules

**Files:**
- Modify: `backend/internal/lib/categorization/rules.json`
- Modify: `backend/internal/lib/categorization/rules_loader.go`
- Modify: `backend/internal/service/categorization_service.go`
- Create: `backend/internal/service/categorization_service_test.go`

**Step 1: Write categorization test**

Create `backend/internal/service/categorization_service_test.go`:

```go
package service

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func TestCategorize_MercadoPago_NotFood(t *testing.T) {
	logger := zerolog.New(nil)

	// Mock repos (nil for this isolated test)
	svc := NewCategorizationService(nil, nil, &logger)

	descriptions := []string{
		"PIX RECEBIDO MERCADO PAGO",
		"Pagamento Mercado Pago",
		"MERCADOPAGO*TRANSFER",
	}

	for _, desc := range descriptions {
		result := svc.shouldExcludeCategory(desc, "alimentação")
		if !result {
			t.Errorf("expected Mercado Pago to exclude 'alimentação', desc: %s", desc)
		}
	}
}

func TestCategorize_Mercado_IsFood(t *testing.T) {
	logger := zerolog.New(nil)
	svc := NewCategorizationService(nil, nil, &logger)

	descriptions := []string{
		"MERCADO SAO JOAO",
		"SUPERMERCADO ABC",
		"Compra no mercado",
	}

	for _, desc := range descriptions {
		result := svc.shouldExcludeCategory(desc, "alimentação")
		if result {
			t.Errorf("expected regular 'mercado' to NOT exclude 'alimentação', desc: %s", desc)
		}
	}
}
```

**Step 2: Run test to verify it fails**

```bash
cd backend
go test ./internal/service/... -v -run TestCategorize
```

Expected: FAIL with "undefined: shouldExcludeCategory"

**Step 3: Add exclusion patterns to rules.json**

In `backend/internal/lib/categorization/rules.json`, add new section after `regex_patterns`:

```json
{
  "expense_patterns": [ ... ],
  "income_patterns": [ ... ],
  "regex_patterns": [ ... ],
  "exclusion_patterns": [
    {
      "keywords": ["mercado pago", "mercadopago", "pagamento mercado pago"],
      "exclude_categories": ["alimentação", "supermercado"],
      "reason": "payment_platform"
    },
    {
      "keywords": ["pagseguro", "pag seguro"],
      "exclude_categories": ["alimentação", "supermercado"],
      "reason": "payment_platform"
    },
    {
      "keywords": ["picpay"],
      "exclude_categories": ["alimentação", "supermercado"],
      "reason": "payment_platform"
    }
  ]
}
```

**Step 4: Update rules_loader.go to load exclusions**

In `backend/internal/lib/categorization/rules_loader.go`:

```go
type ExclusionPattern struct {
	Keywords         []string `json:"keywords"`
	ExcludeCategories []string `json:"exclude_categories"`
	Reason           string   `json:"reason"`
}

type Rules struct {
	ExpensePatterns   []ExpensePattern   `json:"expense_patterns"`
	IncomePatterns    []IncomePattern    `json:"income_patterns"`
	RegexPatterns     []RegexPattern     `json:"regex_patterns"`
	ExclusionPatterns []ExclusionPattern `json:"exclusion_patterns"` // NEW
}
```

**Step 5: Add shouldExcludeCategory method to categorization_service.go**

In `backend/internal/service/categorization_service.go`, add method:

```go
// shouldExcludeCategory checks if a description matches exclusion patterns for a category
func (s *CategorizationService) shouldExcludeCategory(description, categoryName string) bool {
	if s.rules == nil {
		return false
	}

	descLower := strings.ToLower(description)
	categoryLower := strings.ToLower(categoryName)

	for _, pattern := range s.rules.ExclusionPatterns {
		// Check if any keyword matches
		for _, keyword := range pattern.Keywords {
			if strings.Contains(descLower, strings.ToLower(keyword)) {
				// Check if this category should be excluded
				for _, excludeCat := range pattern.ExcludeCategories {
					if strings.ToLower(excludeCat) == categoryLower {
						return true
					}
				}
			}
		}
	}

	return false
}
```

**Step 6: Integrate exclusion check into SuggestCategory**

Find the `SuggestCategory` method in `categorization_service.go` and add exclusion check:

```go
// After finding a potential category match, check exclusions
if s.shouldExcludeCategory(description, category.Name) {
	continue // Skip this category
}
```

**Step 7: Run tests**

```bash
cd backend
go test ./internal/service/... -v -run TestCategorize
```

Expected: PASS

**Step 8: Commit**

```bash
git add backend/internal/lib/categorization/ backend/internal/service/categorization_service.go backend/internal/service/categorization_service_test.go
git commit -m "feat: add categorization exclusion rules for payment platforms

Prevents false positives:
- Mercado Pago not categorized as food/supermarket
- PagSeguro, PicPay also excluded from food categories

Uses multi-word matching for better context awareness.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Task 5: End-to-End Testing with Real BB CSV

**Files:**
- Create: `backend/test/fixtures/bb_sample.csv`
- Create: `backend/test/integration/parser_integration_test.go`

**Step 1: Create test fixture with real BB data**

Create `backend/test/fixtures/bb_sample.csv` (use the provided BB CSV):

```csv
"Data","Lançamento","Detalhes","Nº documento","Valor","Tipo Lançamento"
"24/02/2026","Saldo Anterior","","","-439,69",""
"02/03/2026","Pagamento Pix Cartão Crédito","28/02 23:24 Kauê Rodrigues da Silva","100611000281685","10,00","Entrada"
"02/03/2026","Pix - Enviado","28/02 23:24 Kauê Rodrigues da Silva","30202","-10,00","Saída"
"00/00/0000","Saldo do dia","","","-490,87",""
"05/03/2026","S A L D O","","","-490,87",""
```

**Step 2: Write end-to-end test**

Create `backend/test/integration/parser_integration_test.go`:

```go
package integration

import (
	"os"
	"testing"

	"github.com/gustavoz65/Fyapp/internal/lib/parser"
)

func TestBBCSVImport_E2E(t *testing.T) {
	file, err := os.Open("../fixtures/bb_sample.csv")
	if err != nil {
		t.Fatalf("failed to open fixture: %v", err)
	}
	defer file.Close()

	p := parser.NewTransactionParser()
	transactions, err := p.ParseCSV(file, "bb")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	// Should only import valid transactions:
	// - Skip "Saldo Anterior" (control line)
	// - Skip "Pagamento Pix Cartão" (duplicate)
	// - Import "Pix - Enviado" ✓
	// - Skip "Saldo do dia" (invalid date)
	// - Skip "S A L D O" (control line)

	if len(transactions) != 1 {
		t.Errorf("expected 1 valid transaction, got %d", len(transactions))
		for i, tx := range transactions {
			t.Logf("Transaction %d: %s - %s", i, tx.Date, tx.Description)
		}
	}

	if len(transactions) > 0 {
		tx := transactions[0]
		if tx.Amount.String() != "10" {
			t.Errorf("expected amount 10, got %s", tx.Amount)
		}
		if tx.Type != "expense" {
			t.Errorf("expected type expense, got %s", tx.Type)
		}
	}
}
```

**Step 3: Create test directories**

```bash
mkdir -p backend/test/fixtures
mkdir -p backend/test/integration
```

**Step 4: Run E2E test**

```bash
cd backend
go test ./test/integration/... -v
```

Expected: PASS

**Step 5: Commit**

```bash
git add backend/test/
git commit -m "test: add end-to-end BB CSV import test

Validates complete import pipeline with real BB CSV data.
Ensures all filters work together correctly.

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

## Task 6: Add Missing Category Rules (Deferred)

**Files:**
- Modify: `backend/internal/lib/categorization/rules.json`

**Step 1: Review existing categories and add missing patterns**

This task is marked as **deferred to end**. When ready:

1. Review `rules.json` for gaps
2. Add patterns for common Brazilian transaction types
3. Test with real transaction descriptions
4. Commit with descriptive message

**Placeholder for now** - will be implemented after core bugs are fixed.

---

## Task 7: Run Full Test Suite

**Step 1: Run all tests**

```bash
cd backend
go test ./... -v -cover
```

Expected: PASS (all tests)

**Step 2: Check test coverage**

```bash
cd backend
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

Review coverage report in browser.

**Step 3: If coverage < 70%, add missing tests**

Focus on:
- Edge cases in encoding detection
- Various BB CSV formats
- Categorization exclusion edge cases

---

## Task 8: Manual Verification with Dashboard

**Step 1: Start backend server**

```bash
cd backend
go run cmd/main.go
```

**Step 2: Import BB CSV through API**

Use the BB CSV file provided by user (5 months of data).

**Step 3: Verify dashboard**

Check:
- ✅ No encoding corruption in transaction descriptions
- ✅ No duplicate "Pagamento Pix Cartão" entries
- ✅ Valid dates only (no "Invalid Date")
- ✅ Correct balance calculations
- ✅ "Mercado Pago" not categorized as food

**Step 4: Take screenshots/notes of fixes**

Document before/after for reference.

---

## Success Criteria Checklist

- [ ] Encoding: BB CSV imports without corruption (ã, ç, ê display correctly)
- [ ] Duplicates: "Pagamento Pix Cartão Crédito" filtered out
- [ ] Control Lines: "Saldo Anterior", "Saldo do dia", "S A L D O" filtered out
- [ ] Invalid Dates: "00/00/0000" filtered out before parse
- [ ] Categorization: "Mercado Pago" not categorized as food
- [ ] Dashboard: Shows valid dates on chart X-axis
- [ ] Dashboard: Receitas > R$ 0,00 (not zero)
- [ ] Dashboard: Saldo matches imported data
- [ ] Tests: 100% passing
- [ ] Coverage: >70% on modified files

---

## Post-Implementation Tasks

After all fixes are complete:

1. Update CHANGELOG.md with bug fixes
2. Create PR with detailed description
3. Tag release (e.g., v1.2.1 - Bug Fixes)
4. Update user documentation if needed
5. Monitor for any new edge cases in production

---

## Notes for Implementation

- **TDD Required**: Write test first, watch it fail, implement, watch it pass
- **Frequent Commits**: Commit after each passing test
- **DRY**: Reuse encoding detector, don't duplicate logic
- **YAGNI**: Don't add features not in this plan
- **Error Handling**: Graceful degradation (return as-is if encoding detection fails)
- **Backwards Compatibility**: Don't break existing CSV imports from other banks
