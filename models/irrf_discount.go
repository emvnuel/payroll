package models

import "github.com/shopspring/decimal"

type IRRFDiscount struct {
	GrossPay            decimal.Decimal
	NumberOfDependents  int64
	INSSDeductionAmount decimal.Decimal
}

type IRRFRange struct {
	InitValue decimal.Decimal
	EndValue  decimal.Decimal
	Aliquot   decimal.Decimal
	Deduction decimal.Decimal
}

func NewIRRFRange(
	initValue decimal.Decimal,
	endValue decimal.Decimal,
	aliquot decimal.Decimal,
	deduction decimal.Decimal) IRRFRange {
	irrfRange := IRRFRange{InitValue: initValue, EndValue: endValue, Aliquot: aliquot, Deduction: deduction}
	return irrfRange
}

const DEPENDENT_DEDUCTION_AMOUNT = 189.59

var (
	IRRF_RANGE_1 = NewIRRFRange(decimal.Zero, decimal.NewFromFloat(2112.00), decimal.Zero, decimal.Zero)
	IRRF_RANGE_2 = NewIRRFRange(decimal.NewFromFloat(2112.01), decimal.NewFromFloat(2826.65), decimal.NewFromFloat(0.075), decimal.NewFromFloat(142.80))
	IRRF_RANGE_3 = NewIRRFRange(decimal.NewFromFloat(2826.66), decimal.NewFromFloat(3751.05), decimal.NewFromFloat(0.150), decimal.NewFromFloat(354.80))
	IRRF_RANGE_4 = NewIRRFRange(decimal.NewFromFloat(3751.06), decimal.NewFromFloat(4664.68), decimal.NewFromFloat(0.225), decimal.NewFromFloat(636.13))
	IRRF_RANGE_5 = NewIRRFRange(decimal.NewFromFloat(4664.68), decimal.NewFromFloat(1000000), decimal.NewFromFloat(0.275), decimal.NewFromFloat(869.36))
)

var irrfRanges = []IRRFRange{
	IRRF_RANGE_1, IRRF_RANGE_2, IRRF_RANGE_3, IRRF_RANGE_4, IRRF_RANGE_5,
}

func IRRFRangeByBaseAmount(grossPay decimal.Decimal) IRRFRange {
	item := IRRFRange{}
	for _, irrfRange := range irrfRanges {
		if grossPay.GreaterThanOrEqual(irrfRange.InitValue) && grossPay.LessThanOrEqual(irrfRange.EndValue) {
			item = irrfRange
		}
	}
	return item
}

func (i IRRFDiscount) DependentsDeduction() decimal.Decimal {
	return decimal.NewFromFloat(DEPENDENT_DEDUCTION_AMOUNT).Mul(decimal.NewFromInt(i.NumberOfDependents))
}

func (i IRRFDiscount) TotalDeduction() decimal.Decimal {
	return i.DependentsDeduction().Add(i.INSSDeductionAmount)
}

func (i IRRFDiscount) BaseIRRFAmount() decimal.Decimal {
	return i.GrossPay.Sub(i.TotalDeduction())
}

func (i IRRFDiscount) Value() decimal.Decimal {
	irrfRange := IRRFRangeByBaseAmount(i.BaseIRRFAmount())
	return i.BaseIRRFAmount().Mul(irrfRange.Aliquot).Sub(irrfRange.Deduction).RoundBank(2)

}

func (i IRRFDiscount) Name() string {
	return "IRRF"
}
