package parser

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
	"github.com/shopspring/decimal"
)

type TransactionImport struct {
	Date        time.Time
	Description string
	Amount      decimal.Decimal
	Type        string
	ExternalID  string
	RawData     map[string]string
}

type CSVMapping struct {
	DateColumn        int
	DescriptionColumn int
	AmountColumn      int
	TypeColumn        *int
	DateFormat        string
}

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
		DescriptionColumn: 2,
		AmountColumn:      4,
		TypeColumn:        &[]int{5}[0],
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

func (p *TransactionParser) ParseCSV(reader io.Reader, bankType string) ([]TransactionImport, error) {
	mapping, ok := BankMappings[strings.ToLower(bankType)]
	if !ok {
		mapping = BankMappings["generic"]
	}

	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = -1
	csvReader.TrimLeadingSpace = true

	var transactions []TransactionImport
	seenTransactions := make(map[string]bool)
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

		if lineNumber == 1 {
			if p.isHeaderRow(record) {
				continue
			}
		}

		if len(record) == 0 || (len(record) == 1 && record[0] == "") {
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

	return transactions, nil
}

func (p *TransactionParser) isHeaderRow(record []string) bool {
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

func (p *TransactionParser) parseCSVRecord(record []string, mapping CSVMapping, bankType string) (TransactionImport, error) {
	if len(record) <= mapping.DateColumn || len(record) <= mapping.DescriptionColumn || len(record) <= mapping.AmountColumn {
		return TransactionImport{}, fmt.Errorf("registro com colunas insuficientes: %d", len(record))
	}

	if mapping.TypeColumn != nil && len(record) > *mapping.TypeColumn {
		tipoLancamento := strings.TrimSpace(record[*mapping.TypeColumn])
		if tipoLancamento == "" {
			return TransactionImport{}, fmt.Errorf("linha sem tipo de lançamento (saldo/subtotal)")
		}
	}

	dateStr := strings.TrimSpace(record[mapping.DateColumn])
	if dateStr == "00/00/0000" || dateStr == "" {
		return TransactionImport{}, fmt.Errorf("data inválida ou vazia")
	}

	baseDate, err := p.parseDate(dateStr, mapping.DateFormat)
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

	description := strings.TrimSpace(record[mapping.DescriptionColumn])
	lancamento := ""
	if mapping.TypeColumn != nil && len(record) > 1 {
		lancamento = strings.TrimSpace(record[1])
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

	finalDate := baseDate
	if strings.ToLower(bankType) == "bb" {
		if extractedDate, ok := p.extractDateFromBBDescription(record[mapping.DescriptionColumn], baseDate.Year()); ok {
			finalDate = extractedDate
		}

		if strings.Contains(strings.ToLower(lancamento), "pagamento pix cart") &&
		   strings.Contains(strings.ToLower(lancamento), "cr") {
			transactionType = "expense"
		}
	}

	externalID := p.generateExternalID(finalDate, description, amount.Abs())

	rawData := make(map[string]string)
	for i, value := range record {
		rawData[fmt.Sprintf("col_%d", i)] = value
	}

	return TransactionImport{
		Date:        finalDate,
		Description: description,
		Amount:      amount,
		Type:        transactionType,
		ExternalID:  externalID,
		RawData:     rawData,
	}, nil
}

func (p *TransactionParser) extractDateFromBBDescription(description string, baseYear int) (time.Time, bool) {
	re := regexp.MustCompile(`(\d{2}/\d{2})\s+(\d{2}:\d{2})`)
	matches := re.FindStringSubmatch(description)

	if len(matches) >= 2 {
		dateTimeStr := fmt.Sprintf("%s/%d %s", matches[1], baseYear, matches[2])

		formats := []string{
			"02/01/2006 15:04",
			"01/02/2006 15:04",
		}

		for _, format := range formats {
			if t, err := time.Parse(format, dateTimeStr); err == nil {
				return t, true
			}
		}
	}

	return time.Time{}, false
}

func (p *TransactionParser) parseDate(dateStr, format string) (time.Time, error) {
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
	cleaned := strings.TrimSpace(amountStr)
	cleaned = strings.ReplaceAll(cleaned, "R$", "")
	cleaned = strings.ReplaceAll(cleaned, "$", "")
	cleaned = strings.ReplaceAll(cleaned, " ", "")

	isNegative := false
	if strings.HasPrefix(cleaned, "-") || strings.HasPrefix(cleaned, "(") {
		isNegative = true
		cleaned = strings.TrimPrefix(cleaned, "-")
		cleaned = strings.Trim(cleaned, "()")
	}

	dotCount := strings.Count(cleaned, ".")
	commaCount := strings.Count(cleaned, ",")

	if commaCount > 0 && dotCount > 0 {
		lastDot := strings.LastIndex(cleaned, ".")
		lastComma := strings.LastIndex(cleaned, ",")

		if lastComma > lastDot {
			cleaned = strings.ReplaceAll(cleaned, ".", "")
			cleaned = strings.ReplaceAll(cleaned, ",", ".")
		} else {
			cleaned = strings.ReplaceAll(cleaned, ",", "")
		}
	} else if commaCount > 0 {
		parts := strings.Split(cleaned, ",")
		if len(parts) == 2 && len(parts[1]) == 2 {
			cleaned = strings.ReplaceAll(cleaned, ",", ".")
		} else {
			cleaned = strings.ReplaceAll(cleaned, ",", "")
		}
	}

	amount, err := decimal.NewFromString(cleaned)
	if err != nil {
		floatVal, err := strconv.ParseFloat(cleaned, 64)
		if err != nil {
			return decimal.Zero, "", fmt.Errorf("invalid amount: %s", amountStr)
		}
		amount = decimal.NewFromFloat(floatVal)
	}

	amount = amount.Abs()

	transactionType := "income"
	if isNegative {
		transactionType = "expense"
	}

	return amount, transactionType, nil
}

func (p *TransactionParser) generateExternalID(date time.Time, description string, amount decimal.Decimal) string {
	data := fmt.Sprintf("%s|%s|%s",
		date.Format("2006-01-02 15:04"),
		strings.ToLower(strings.TrimSpace(description)),
		amount.String(),
	)

	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

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
