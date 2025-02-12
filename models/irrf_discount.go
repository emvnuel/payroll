package models

import (
	"encoding/json"
	"os"

	"github.com/shopspring/decimal"

	_ "github.com/joho/godotenv/autoload"
)

var (
	DEPENDENT_DEDUCTION_AMOUNT    = os.Getenv("DEPENDENT_DEDUCTION_AMOUNT")
	dependentDeductionAmount, _   = decimal.NewFromString(DEPENDENT_DEDUCTION_AMOUNT)
	IRRFRanges                    = loadIRRFRangesFromEnv()
	simplifiedDeductionPercentage = decimal.NewFromFloat32(0.25)
)

type IRRFDiscount struct {
	GrossPay            decimal.Decimal
	NumberOfDependents  int64
	INSSDeductionAmount decimal.Decimal
	SimplifiedDeduction bool
}

type IRRFRange struct {
	StartingValue decimal.Decimal `json:"init_value"`
	EndingValue   decimal.Decimal `json:"end_value"`
	Aliquot       decimal.Decimal `json:"aliquot"`
	Deduction     decimal.Decimal `json:"deduction"`
}

func NewIRRFDiscount(grossPay decimal.Decimal, numberOfDependents int64, inssDeductionAmount decimal.Decimal, simplifiedDeduction bool) *IRRFDiscount {
	return &IRRFDiscount{
		GrossPay:            grossPay,
		NumberOfDependents:  numberOfDependents,
		INSSDeductionAmount: inssDeductionAmount,
		SimplifiedDeduction: simplifiedDeduction,
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

func (i *IRRFDiscount) dependentsDeduction() decimal.Decimal {
	return dependentDeductionAmount.Mul(decimal.NewFromInt(i.NumberOfDependents))
}

func (i *IRRFDiscount) totalDeduction() decimal.Decimal {
	if i.SimplifiedDeduction {
		return i.simplifiedDeductionAmount()
	}
	return i.dependentsDeduction().Add(i.INSSDeductionAmount)
}

func (i *IRRFDiscount) simplifiedDeductionAmount() decimal.Decimal {
	if len(IRRFRanges) == 0 {
		return decimal.Zero
	}
	return IRRFRanges[0].EndingValue.Mul(simplifiedDeductionPercentage)
}

func (i *IRRFDiscount) taxableBase() decimal.Decimal {
	return i.GrossPay.Sub(i.totalDeduction())
}

func (i *IRRFDiscount) Value() decimal.Decimal {
	matchingRange := i.findMatchingRange()
	if matchingRange == nil {
		return decimal.Zero
	}

	taxableBaseProduct := i.taxableBase().Mul(matchingRange.Aliquot)
	return taxableBaseProduct.Sub(matchingRange.Deduction).RoundBank(2)
}

func (i *IRRFDiscount) findMatchingRange() *IRRFRange {
	taxBase := i.taxableBase()
	for _, irrfRange := range IRRFRanges {
		if taxBase.GreaterThanOrEqual(irrfRange.StartingValue) && taxBase.LessThanOrEqual(irrfRange.EndingValue) {
			return &irrfRange
		}
	}
	return nil
}

func (i *IRRFDiscount) Name() string {
	return "IRRF"
}
