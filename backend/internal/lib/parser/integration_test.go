package parser

import (
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
