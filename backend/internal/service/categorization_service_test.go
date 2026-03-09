package service

import (
	"testing"

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
