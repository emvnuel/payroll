package models

import (
	"encoding/json"
	"os"

	"github.com/shopspring/decimal"
)

const DEPENDENT_DEDUCTION_AMOUNT = 189.59

var IRRFRanges = loadIRRFRangesFromFile("./resources/irrf_ranges.json")

type IRRFDiscount struct {
	GrossPay            decimal.Decimal
	NumberOfDependents  int64
	INSSDeductionAmount decimal.Decimal
}

type IRRFRange struct {
	StartingValue decimal.Decimal `json:"starting_value"`
	EndingValue   decimal.Decimal `json:"ending_value"`
	Aliquot       decimal.Decimal `json:"aliquot"`
	Deduction     decimal.Decimal `json:"deduction"`
}

func NewIRRFRange(startingValue, endingValue, aliquot, deduction decimal.Decimal) IRRFRange {
	return IRRFRange{
		StartingValue: startingValue,
		EndingValue:   endingValue,
		Aliquot:       aliquot,
		Deduction:     deduction,
	}
}

func loadIRRFRangesFromFile(filename string) []IRRFRange {
	data, err := os.ReadFile(filename)

	if err != nil {
		panic(err)
	}
	var irrfRanges []IRRFRange
	err = json.Unmarshal(data, &irrfRanges)

	if err != nil {
		panic(err)
	}
	return irrfRanges
}

func (i IRRFDiscount) dependentsDeduction() decimal.Decimal {
	return decimal.NewFromFloat(DEPENDENT_DEDUCTION_AMOUNT).Mul(decimal.NewFromInt(i.NumberOfDependents))
}

func (i IRRFDiscount) totalDeduction() decimal.Decimal {
	return i.dependentsDeduction().Add(i.INSSDeductionAmount)
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
