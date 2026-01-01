package models

import (
	"encoding/json"
	"os"

	"github.com/shopspring/decimal"

	_ "github.com/joho/godotenv/autoload"
)

var (
	DEPENDENT_DEDUCTION_AMOUNT           = os.Getenv("DEPENDENT_DEDUCTION_AMOUNT")
	dependentDeductionAmount, _          = decimal.NewFromString(DEPENDENT_DEDUCTION_AMOUNT)
	IRRFRanges                           = loadIRRFRangesFromEnv()
	IRRF_SIMPLIFIED_DEDUCTION_PERCENTAGE = getEnvOrDefault("IRRF_SIMPLIFIED_DEDUCTION_PERCENTAGE", "0.25")
	simplifiedDeductionPercentage, _     = decimal.NewFromString(IRRF_SIMPLIFIED_DEDUCTION_PERCENTAGE)

	// Nova regra de redução do IRRF 2026 (Lei nº 15.270/2025)
	// Ampliação da faixa de isenção para rendimentos até R$ 5.000,00
	IRRF_MAX_REDUCTION_AMOUNT  = getEnvOrDefault("IRRF_MAX_REDUCTION_AMOUNT", "312.89")
	maxReductionAmount, _      = decimal.NewFromString(IRRF_MAX_REDUCTION_AMOUNT)
	IRRF_REDUCTION_THRESHOLD   = getEnvOrDefault("IRRF_REDUCTION_THRESHOLD", "5000.00")
	reductionThreshold, _      = decimal.NewFromString(IRRF_REDUCTION_THRESHOLD)
	IRRF_REDUCTION_UPPER_LIMIT = getEnvOrDefault("IRRF_REDUCTION_UPPER_LIMIT", "7350.00")
	reductionUpperLimit, _     = decimal.NewFromString(IRRF_REDUCTION_UPPER_LIMIT)
	IRRF_REDUCTION_CONSTANT    = getEnvOrDefault("IRRF_REDUCTION_CONSTANT", "978.62")
	reductionConstant, _       = decimal.NewFromString(IRRF_REDUCTION_CONSTANT)
	IRRF_REDUCTION_MULTIPLIER  = getEnvOrDefault("IRRF_REDUCTION_MULTIPLIER", "0.133145")
	reductionMultiplier, _     = decimal.NewFromString(IRRF_REDUCTION_MULTIPLIER)
)

type IRRFDiscount struct {
	GrossPay            decimal.Decimal
	NumberOfDependents  int64
	INSSDeductionAmount decimal.Decimal
}

type IRRFRange struct {
	StartingValue decimal.Decimal `json:"init_value"`
	EndingValue   decimal.Decimal `json:"end_value"`
	Aliquot       decimal.Decimal `json:"aliquot"`
	Deduction     decimal.Decimal `json:"deduction"`
}

func NewIRRFDiscount(grossPay decimal.Decimal, numberOfDependents int64, inssDeductionAmount decimal.Decimal) *IRRFDiscount {
	return &IRRFDiscount{
		GrossPay:            grossPay,
		NumberOfDependents:  numberOfDependents,
		INSSDeductionAmount: inssDeductionAmount,
	}
}

func NewIRRFRange(startingValue, endingValue, aliquot, deduction decimal.Decimal) *IRRFRange {
	return &IRRFRange{
		StartingValue: startingValue,
		EndingValue:   endingValue,
		Aliquot:       aliquot,
		Deduction:     deduction,
	}
}

func loadIRRFRangesFromEnv() []IRRFRange {
	data := os.Getenv("IRRF_RANGES")
	var irrfRanges []IRRFRange
	if err := json.Unmarshal([]byte(data), &irrfRanges); err != nil {
		// Consider logging the error or handling it appropriately
		return nil
	}
	return irrfRanges
}

// getEnvOrDefault retorna o valor da variável de ambiente ou um valor padrão se não estiver definida
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func (i *IRRFDiscount) dependentsDeduction() decimal.Decimal {
	return dependentDeductionAmount.Mul(decimal.NewFromInt(i.NumberOfDependents))
}

func (i *IRRFDiscount) simplifiedDeductionAmount() decimal.Decimal {
	if len(IRRFRanges) == 0 {
		return decimal.Zero
	}
	return IRRFRanges[0].EndingValue.Mul(simplifiedDeductionPercentage)
}

// totalDeductionWithDependents calcula a dedução usando dependentes + INSS
func (i *IRRFDiscount) totalDeductionWithDependents() decimal.Decimal {
	return i.dependentsDeduction().Add(i.INSSDeductionAmount)
}

// totalDeductionSimplified calcula a dedução usando desconto simplificado
func (i *IRRFDiscount) totalDeductionSimplified() decimal.Decimal {
	return i.simplifiedDeductionAmount()
}

func (i *IRRFDiscount) taxableBaseWithDeduction(deduction decimal.Decimal) decimal.Decimal {
	return i.GrossPay.Sub(deduction)
}

// Value calcula o IRRF usando a opção mais favorável ao contribuinte
// (desconto simplificado vs dedução de dependentes + INSS)
func (i *IRRFDiscount) Value() decimal.Decimal {
	// Calcula IRRF com desconto simplificado
	simplifiedDeduction := i.totalDeductionSimplified()
	irrfSimplified := i.calculateIRRFWithDeduction(simplifiedDeduction)

	// Calcula IRRF com dedução de dependentes + INSS
	dependentsDeduction := i.totalDeductionWithDependents()
	irrfDependents := i.calculateIRRFWithDeduction(dependentsDeduction)

	// Retorna o MENOR valor (mais favorável ao contribuinte)
	if irrfSimplified.LessThan(irrfDependents) {
		return irrfSimplified
	}
	return irrfDependents
}

// calculateIRRFWithDeduction calcula o IRRF para uma dedução específica
func (i *IRRFDiscount) calculateIRRFWithDeduction(deduction decimal.Decimal) decimal.Decimal {
	taxableBase := i.taxableBaseWithDeduction(deduction)

	matchingRange := i.findMatchingRangeForBase(taxableBase)
	if matchingRange == nil {
		return decimal.Zero
	}

	// Calcula o imposto pela tabela progressiva
	taxableBaseProduct := taxableBase.Mul(matchingRange.Aliquot)
	calculatedTax := taxableBaseProduct.Sub(matchingRange.Deduction)

	// Aplica a redução conforme a Lei nº 15.270/2025
	reduction := i.calculateReduction(calculatedTax)
	finalTax := calculatedTax.Sub(reduction)

	// Garante que o imposto final não seja negativo
	if finalTax.LessThan(decimal.Zero) {
		return decimal.Zero
	}

	return finalTax.RoundBank(2)
}

func (i *IRRFDiscount) findMatchingRangeForBase(taxBase decimal.Decimal) *IRRFRange {
	for _, irrfRange := range IRRFRanges {
		if taxBase.GreaterThanOrEqual(irrfRange.StartingValue) && taxBase.LessThanOrEqual(irrfRange.EndingValue) {
			return &irrfRange
		}
	}
	return nil
}

// calculateReduction calcula a redução do imposto conforme a Lei nº 15.270/2025
// Regras:
// - Até R$ 5.000,00: redução de até R$ 312,89 (limitado ao imposto calculado)
// - Entre R$ 5.000,01 e R$ 7.350,00: redução gradual usando fórmula: R$ 978,62 - (0,133145 x rendimento)
// - Acima de R$ 7.350,00: sem redução
func (i *IRRFDiscount) calculateReduction(calculatedTax decimal.Decimal) decimal.Decimal {
	grossPay := i.GrossPay

	// Acima de R$ 7.350,00: sem redução
	if grossPay.GreaterThan(reductionUpperLimit) {
		return decimal.Zero
	}

	var reduction decimal.Decimal

	// Até R$ 5.000,00: redução máxima de R$ 312,89
	if grossPay.LessThanOrEqual(reductionThreshold) {
		reduction = maxReductionAmount
	} else {
		// Entre R$ 5.000,01 e R$ 7.350,00: redução gradual
		// Fórmula: R$ 978,62 - (0,133145 x rendimento)
		reduction = reductionConstant.Sub(reductionMultiplier.Mul(grossPay))

		// Garantir que a redução não seja negativa
		if reduction.LessThan(decimal.Zero) {
			reduction = decimal.Zero
		}
	}

	// A redução não pode ser maior que o imposto calculado
	if reduction.GreaterThan(calculatedTax) {
		return calculatedTax
	}

	return reduction
}

func (i *IRRFDiscount) Name() string {
	return "IRRF"
}
