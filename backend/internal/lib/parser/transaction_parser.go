package parser

import (
	"crypto/md5"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransactionImport struct {
	Date        time.Time
	Description string
	Amount      decimal.Decimal
	Type        string // "income" or "expense"
	ExternalID  string // Hash for deduplication
	RawData     map[string]string
}

type CSVMapping struct {
	DateColumn        int
	DescriptionColumn int
	AmountColumn      int
	TypeColumn        *int // Optional, will infer from amount if not present
	DateFormat        string
}

// BankMappings contains default CSV column mappings for common Brazilian banks
var BankMappings = map[string]CSVMapping{
	"nubank": {
		DateColumn:        0,
		DescriptionColumn: 3,
		AmountColumn:      1,
		DateFormat:        "02/01/2006",
	},
	"inter": {
		DateColumn:        0,
		DescriptionColumn: 1,
		AmountColumn:      2,
		DateFormat:        "02/01/2006",
	},
	"itau": {
		DateColumn:        0,
		DescriptionColumn: 1,
		AmountColumn:      3,
		DateFormat:        "02/01/2006",
	},
	"bb": {
		DateColumn:        0,
		DescriptionColumn: 2,            // Detalhes column
		AmountColumn:      4,            // Valor column
		TypeColumn:        &[]int{5}[0], // Tipo Lançamento (Entrada/Saída)
		DateFormat:        "02/01/2006",
	},
	"generic": {
		DateColumn:        0,
		DescriptionColumn: 1,
		AmountColumn:      2,
		DateFormat:        "02/01/2006",
	},
}

type TransactionParser struct{}

func NewTransactionParser() *TransactionParser {
	return &TransactionParser{}
}

// ParseCSV parses CSV file and returns transactions
func (p *TransactionParser) ParseCSV(reader io.Reader, bankType string) ([]TransactionImport, error) {
	mapping, ok := BankMappings[strings.ToLower(bankType)]
	if !ok {
		mapping = BankMappings["generic"]
	}

	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = -1 // Lida com linhas com número variável de colunas
	csvReader.TrimLeadingSpace = true

	var transactions []TransactionImport
	lineNumber := 0

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("erro ao ler a linha %d do CSV: %w", lineNumber, err)
		}

		lineNumber++

		// Pular linha de cabeçalho
		if lineNumber == 1 {
			if p.isHeaderRow(record) {
				continue
			}
		}

		// Pular linhas vazias
		if len(record) == 0 || (len(record) == 1 && record[0] == "") {
			continue
		}

		transaction, err := p.parseCSVRecord(record, mapping)
		if err != nil {
			// Log error but continue processing
			fmt.Printf("Aviso: pulando a linha %d: %v\n", lineNumber, err)
			continue
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

// isHeaderRow tenta identificar se uma linha é um cabeçalho
func (p *TransactionParser) isHeaderRow(record []string) bool {
	// Check if row contains common header keywords
	headerKeywords := []string{"data", "date", "descri", "valor", "amount", "tipo", "type"}

	for _, field := range record {
		fieldLower := strings.ToLower(field)
		for _, keyword := range headerKeywords {
			if strings.Contains(fieldLower, keyword) {
				return true
			}
		}
	}
	return false
}

// parseCSVRecord faz o parsing de uma linha do CSV usando o mapeamento fornecido e retorna um TransactionImport
func (p *TransactionParser) parseCSVRecord(record []string, mapping CSVMapping) (TransactionImport, error) {
	if len(record) <= mapping.DateColumn || len(record) <= mapping.DescriptionColumn || len(record) <= mapping.AmountColumn {
		return TransactionImport{}, fmt.Errorf("registro com colunas insuficientes: %d", len(record))
	}

	if mapping.TypeColumn != nil && len(record) > *mapping.TypeColumn {
		tipoLancamento := strings.TrimSpace(record[*mapping.TypeColumn])
		if tipoLancamento == "" {
			return TransactionImport{}, fmt.Errorf("linha sem tipo de lançamento (saldo/subtotal)")
		}
	}

	// Parse date
	dateStr := strings.TrimSpace(record[mapping.DateColumn])
	// Skip invalid dates like "00/00/0000" from Banco do Brasil
	if dateStr == "00/00/0000" || dateStr == "" {
		return TransactionImport{}, fmt.Errorf("data inválida ou vazia")
	}
	date, err := p.parseDate(dateStr, mapping.DateFormat)
	if err != nil {
		return TransactionImport{}, fmt.Errorf("data inválida '%s': %w", dateStr, err)
	}

	amountStr := strings.TrimSpace(record[mapping.AmountColumn])
	if amountStr == "" {
		return TransactionImport{}, fmt.Errorf("valor vazio")
	}
	amount, transactionType, err := p.parseAmount(amountStr)
	if err != nil {
		return TransactionImport{}, fmt.Errorf("valor inválido '%s': %w", amountStr, err)
	}

	// Parse description
	description := strings.TrimSpace(record[mapping.DescriptionColumn])
	if mapping.TypeColumn != nil && len(record) > 1 {
		lancamento := strings.TrimSpace(record[1])
		detalhes := strings.TrimSpace(record[mapping.DescriptionColumn])
		if lancamento != "" && detalhes != "" {
			description = lancamento + " - " + detalhes
		} else if lancamento != "" {
			description = lancamento
		}
	}

	if description == "" {
		description = "Transação importada"
	}

	// Generate external ID for deduplication
	externalID := p.generateExternalID(date, description, amount.Abs())

	// Store raw data
	rawData := make(map[string]string)
	for i, value := range record {
		rawData[fmt.Sprintf("col_%d", i)] = value
	}

	return TransactionImport{
		Date:        date,
		Description: description,
		Amount:      amount,
		Type:        transactionType,
		ExternalID:  externalID,
		RawData:     rawData,
	}, nil
}

func (p *TransactionParser) parseDate(dateStr, format string) (time.Time, error) {
	// Try common formats if the specified one fails
	formats := []string{
		format,
		"2006-01-02",
		"02/01/2006",
		"01/02/2006",
		"2006/01/02",
		"02-01-2006",
		"01-02-2006",
	}

	for _, fmt := range formats {
		if t, err := time.Parse(fmt, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

func (p *TransactionParser) parseAmount(amountStr string) (decimal.Decimal, string, error) {
	// Clean amount string
	cleaned := strings.TrimSpace(amountStr)
	cleaned = strings.ReplaceAll(cleaned, "R$", "")
	cleaned = strings.ReplaceAll(cleaned, "$", "")
	cleaned = strings.ReplaceAll(cleaned, " ", "")

	// Determine if it's negative (expense) or positive (income)
	isNegative := false
	if strings.HasPrefix(cleaned, "-") || strings.HasPrefix(cleaned, "(") {
		isNegative = true
		cleaned = strings.TrimPrefix(cleaned, "-")
		cleaned = strings.Trim(cleaned, "()")
	}

	// Handle Brazilian format (1.234,56) vs US format (1,234.56)
	dotCount := strings.Count(cleaned, ".")
	commaCount := strings.Count(cleaned, ",")

	if commaCount > 0 && dotCount > 0 {
		// Has both - determine which is decimal separator
		lastDot := strings.LastIndex(cleaned, ".")
		lastComma := strings.LastIndex(cleaned, ",")

		if lastComma > lastDot {
			// Brazilian format: 1.234,56
			cleaned = strings.ReplaceAll(cleaned, ".", "")
			cleaned = strings.ReplaceAll(cleaned, ",", ".")
		} else {
			// US format: 1,234.56
			cleaned = strings.ReplaceAll(cleaned, ",", "")
		}
	} else if commaCount > 0 {
		// Only commas - check if it's decimal or thousands
		parts := strings.Split(cleaned, ",")
		if len(parts) == 2 && len(parts[1]) == 2 {
			// Likely decimal: 1234,56
			cleaned = strings.ReplaceAll(cleaned, ",", ".")
		} else {
			// Likely thousands: 1,234
			cleaned = strings.ReplaceAll(cleaned, ",", "")
		}
	}

	// Parse to decimal
	amount, err := decimal.NewFromString(cleaned)
	if err != nil {
		// Try as float
		floatVal, err := strconv.ParseFloat(cleaned, 64)
		if err != nil {
			return decimal.Zero, "", fmt.Errorf("invalid amount: %s", amountStr)
		}
		amount = decimal.NewFromFloat(floatVal)
	}

	// Make amount positive (we determine type separately)
	amount = amount.Abs()

	// Determine transaction type
	transactionType := "income"
	if isNegative {
		transactionType = "expense"
	}

	return amount, transactionType, nil
}

func (p *TransactionParser) generateExternalID(date time.Time, description string, amount decimal.Decimal) string {
	// Create hash from date + description + amount for deduplication
	data := fmt.Sprintf("%s|%s|%s",
		date.Format("2006-01-02"),
		strings.ToLower(strings.TrimSpace(description)),
		amount.String(),
	)

	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

// ValidateImport checks for duplicates and returns stats
func (p *TransactionParser) ValidateImport(transactions []TransactionImport, userID uuid.UUID, existingExternalIDs map[string]bool) (new, duplicate int) {
	for _, tx := range transactions {
		if existingExternalIDs[tx.ExternalID] {
			duplicate++
		} else {
			new++
		}
	}
	return new, duplicate
}
