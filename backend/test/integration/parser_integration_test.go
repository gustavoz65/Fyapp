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
