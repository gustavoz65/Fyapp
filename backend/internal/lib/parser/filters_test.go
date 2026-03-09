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
