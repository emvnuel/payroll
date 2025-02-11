package models

import (
	"encoding/json"
	"os"

	"github.com/shopspring/decimal"

	_ "github.com/joho/godotenv/autoload"
)

var DEPENDENT_DEDUCTION_AMOUNT = os.Getenv("DEPENDENT_DEDUCTION_AMOUNT")
var dependentDeductionAmount, _ = decimal.NewFromString(DEPENDENT_DEDUCTION_AMOUNT)
var IRRFRanges = loadIRRFRangesFromEnv()

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

func NewIRRFDiscount(grossPay decimal.Decimal,
	numberOfDependents int64,
	inssDeductionAmount decimal.Decimal,
	simplifiedDeduction bool) *IRRFDiscount {
	return &IRRFDiscount{GrossPay: grossPay, NumberOfDependents: numberOfDependents, INSSDeductionAmount: inssDeductionAmount, SimplifiedDeduction: simplifiedDeduction}

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
	json.Unmarshal([]byte(data), &irrfRanges)
	return irrfRanges
}

func (i IRRFDiscount) dependentsDeduction() decimal.Decimal {
	return dependentDeductionAmount.Mul(decimal.NewFromInt(i.NumberOfDependents))
}

func (i IRRFDiscount) totalDeduction() decimal.Decimal {
	if i.SimplifiedDeduction {
		return simplifiedDeductionAmount()
	}
	return i.dependentsDeduction().Add(i.INSSDeductionAmount)
}

func simplifiedDeductionAmount() decimal.Decimal {
	return IRRFRanges[0].EndingValue.Mul(decimal.NewFromFloat32(0.25))
}

func (i IRRFDiscount) taxableBase() decimal.Decimal {
	return i.GrossPay.Sub(i.totalDeduction())
}

func (i IRRFDiscount) Value() decimal.Decimal {
	matchingRange := findMatchingRange(i.taxableBase())
	taxableBaseProduct := i.taxableBase().Mul(matchingRange.Aliquot)
	return taxableBaseProduct.Sub(matchingRange.Deduction).RoundBank(2)
}

func findMatchingRange(baseAmount decimal.Decimal) *IRRFRange {
	for _, irrfRange := range IRRFRanges {
		if baseAmount.GreaterThanOrEqual(irrfRange.StartingValue) && baseAmount.LessThanOrEqual(irrfRange.EndingValue) {
			return &irrfRange
		}
	}
	return nil
}

func (i IRRFDiscount) Name() string {
	return "IRRF"
}
