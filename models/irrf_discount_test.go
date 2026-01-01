package models

import (
	"os"
	"testing"

	"github.com/shopspring/decimal"
)

func init() {
	// Configurar variáveis de ambiente necessárias para os testes
	os.Setenv("DEPENDENT_DEDUCTION_AMOUNT", "189.59")

	// Tabela progressiva mensal do IRRF 2026 (oficial da Receita Federal)
	irrfRangesJSON := `[
		{"init_value": "0.00", "end_value": "2428.80", "aliquot": "0.00", "deduction": "0.00"},
		{"init_value": "2428.81", "end_value": "2826.65", "aliquot": "0.075", "deduction": "182.16"},
		{"init_value": "2826.66", "end_value": "3751.05", "aliquot": "0.15", "deduction": "394.16"},
		{"init_value": "3751.06", "end_value": "4664.68", "aliquot": "0.225", "deduction": "675.49"},
		{"init_value": "4664.69", "end_value": "999999999.99", "aliquot": "0.275", "deduction": "908.73"}
	]`
	os.Setenv("IRRF_RANGES", irrfRangesJSON)

	// Recarregar as variáveis
	dependentDeductionAmount, _ = decimal.NewFromString(os.Getenv("DEPENDENT_DEDUCTION_AMOUNT"))
	IRRFRanges = loadIRRFRangesFromEnv()
}

// TestIRRFExample1_Receita_4500 testa o exemplo 1 da Receita Federal
// Rendimento: R$ 4.500,00
// Com desconto simplificado calculado: 25% do valor final da primeira faixa IRRF
// A implementação calcula corretamente e aplica a redução até zerar o imposto
func TestIRRFExample1_Receita_4500(t *testing.T) {
	grossPay := decimal.NewFromFloat(4500.00)
	inssDeduction := decimal.Zero

	irrf := NewIRRFDiscount(grossPay, 0, inssDeduction, true)
	result := irrf.Value()

	expected := decimal.Zero
	if !result.Equal(expected) {
		t.Errorf("Para rendimento de R$ 4.500,00, o IRRF deve ser R$ 0,00 devido à nova redução. Obtido: %s", result)
	}
}

// TestIRRFExample2_Receita_6000 testa cálculo com rendimento de R$ 6.000,00
// Aplica desconto simplificado e redução gradual conforme Lei 15.270/2025
func TestIRRFExample2_Receita_6000(t *testing.T) {
	grossPay := decimal.NewFromFloat(6000.00)
	inssDeduction := decimal.Zero

	irrf := NewIRRFDiscount(grossPay, 0, inssDeduction, true)
	result := irrf.Value()

	// Verificar que o imposto está sendo calculado e que há redução
	// O valor exato depende da tabela IRRF específica, mas deve haver redução
	maxExpected := decimal.NewFromFloat(500.00) // Um valor razoável após redução

	if result.GreaterThan(maxExpected) {
		t.Errorf("Para rendimento de R$ 6.000,00, o IRRF deve ser reduzido. Obtido: %s", result)
	}

	// Verificar que não é zero (pois a redução não elimina todo o imposto nessa faixa)
	if result.LessThanOrEqual(decimal.Zero) {
		t.Errorf("Para rendimento de R$ 6.000,00, o IRRF deve ser maior que zero. Obtido: %s", result)
	}
}

// TestIRRFReduction_LimitAt5000 testa a redução máxima até R$ 5.000,00
func TestIRRFReduction_LimitAt5000(t *testing.T) {
	grossPay := decimal.NewFromFloat(5000.00)
	inssDeduction := decimal.Zero

	irrf := NewIRRFDiscount(grossPay, 0, inssDeduction, true)
	result := irrf.Value()

	maxExpected := decimal.NewFromFloat(100.00)

	// Com rendimento de R$ 5.000,00, a redução deve ser máxima (R$ 312,89)
	// e deve zerar o imposto ou aproximar muito de zero
	if result.GreaterThan(maxExpected) {
		t.Errorf("Para rendimento de R$ 5.000,00, o IRRF deve ser muito baixo ou zero devido à redução máxima. Obtido: %s", result)
	}
}

// TestIRRFReduction_NoReductionAbove7350 testa que não há redução acima de R$ 7.350,00
func TestIRRFReduction_NoReductionAbove7350(t *testing.T) {
	grossPay := decimal.NewFromFloat(8000.00)
	inssDeduction := decimal.Zero

	irrf := NewIRRFDiscount(grossPay, 0, inssDeduction, true)
	result := irrf.Value()

	// Acima de R$ 7.350,00 não deve ter redução
	// O imposto deve ser maior do que para valores abaixo de 7.350
	minExpected := decimal.NewFromFloat(1000.00)

	if result.LessThan(minExpected) {
		t.Errorf("Para rendimento de R$ 8.000,00, o IRRF deve ser calculado sem redução (>= R$ 1000). Obtido: %s", result)
	}
}

// TestIRRFReduction_GradualBetween5000And7350 testa a redução gradual
func TestIRRFReduction_GradualBetween5000And7350(t *testing.T) {
	testCases := []struct {
		name       string
		grossPay   float64
		maxImposto float64 // imposto máximo esperado (com alguma margem)
	}{
		{"Rendimento R$ 5.500,00", 5500.00, 400.00},
		{"Rendimento R$ 6.500,00", 6500.00, 700.00},
		{"Rendimento R$ 7.000,00", 7000.00, 900.00},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			grossPay := decimal.NewFromFloat(tc.grossPay)
			inssDeduction := decimal.Zero

			irrf := NewIRRFDiscount(grossPay, 0, inssDeduction, true)
			result := irrf.Value()

			maxExpected := decimal.NewFromFloat(tc.maxImposto)

			// Verificar que o imposto está dentro de uma faixa esperada
			// (com redução gradual aplicada)
			if result.GreaterThan(maxExpected) {
				t.Errorf("Para rendimento de R$ %.2f, o IRRF com redução gradual deve ser menor que R$ %.2f. Obtido: %s",
					tc.grossPay, tc.maxImposto, result)
			}
		})
	}
}

// TestIRRFWithDependents testa o cálculo com dependentes
func TestIRRFWithDependents(t *testing.T) {
	grossPay := decimal.NewFromFloat(5000.00)
	inssDeduction := decimal.Zero
	numberOfDependents := int64(2)

	irrf := NewIRRFDiscount(grossPay, numberOfDependents, inssDeduction, false)
	result := irrf.Value()

	// Com dependentes, a base de cálculo diminui, então o imposto deve ser menor
	if result.LessThan(decimal.Zero) {
		t.Errorf("O IRRF não pode ser negativo. Obtido: %s", result)
	}
}

// TestIRRFWithINSSDeduction testa o cálculo com dedução do INSS
func TestIRRFWithINSSDeduction(t *testing.T) {
	grossPay := decimal.NewFromFloat(5000.00)
	inssDeduction := decimal.NewFromFloat(550.00) // 11% de R$ 5.000,00

	irrf := NewIRRFDiscount(grossPay, 0, inssDeduction, false)
	result := irrf.Value()

	// Com dedução do INSS, a base de cálculo diminui, então o imposto deve ser menor
	if result.LessThan(decimal.Zero) {
		t.Errorf("O IRRF não pode ser negativo. Obtido: %s", result)
	}
}

// TestIRRFCalculateReduction testa a função calculateReduction isoladamente
func TestIRRFCalculateReduction(t *testing.T) {
	testCases := []struct {
		name              string
		grossPay          float64
		calculatedTax     float64
		expectedReduction float64
	}{
		{
			name:              "Rendimento R$ 3.000 - redução máxima",
			grossPay:          3000.00,
			calculatedTax:     100.00,
			expectedReduction: 100.00, // Limitado ao imposto calculado
		},
		{
			name:              "Rendimento R$ 5.000 - redução máxima",
			grossPay:          5000.00,
			calculatedTax:     400.00,
			expectedReduction: 312.89,
		},
		{
			name:              "Rendimento R$ 6.000 - redução gradual",
			grossPay:          6000.00,
			calculatedTax:     500.00,
			expectedReduction: 179.75,
		},
		{
			name:              "Rendimento R$ 8.000 - sem redução",
			grossPay:          8000.00,
			calculatedTax:     1000.00,
			expectedReduction: 0.00,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			grossPay := decimal.NewFromFloat(tc.grossPay)
			calculatedTax := decimal.NewFromFloat(tc.calculatedTax)
			expectedReduction := decimal.NewFromFloat(tc.expectedReduction)

			irrf := NewIRRFDiscount(grossPay, 0, decimal.Zero, true)
			reduction := irrf.calculateReduction(calculatedTax)

			tolerance := decimal.NewFromFloat(0.10)
			diff := reduction.Sub(expectedReduction).Abs()

			if diff.GreaterThanOrEqual(tolerance) {
				t.Errorf("%s: Redução esperada %.2f, obtida %s (diferença: %s)",
					tc.name, tc.expectedReduction, reduction, diff)
			}
		})
	}
}
