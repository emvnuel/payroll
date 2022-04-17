package models

import (
	"github.com/shopspring/decimal"
)

type INSSDiscount struct {
	GrossPay decimal.Decimal
}

type INSSRange struct {
	Index     int
	Aliquot   decimal.Decimal
	InitValue decimal.Decimal
	EndValue  decimal.Decimal
	Prev      *INSSRange
}

func NewINSSRange(
	index int,
	aliquot decimal.Decimal,
	initValue decimal.Decimal,
	endValue decimal.Decimal,
	prev *INSSRange) INSSRange {

	inssRange := INSSRange{Index: index, Aliquot: aliquot, InitValue: initValue, EndValue: endValue, Prev: prev}
	return inssRange
}

const INSS_RANGE_5_DISCOUNT_AMOUNT = 828.38

var (
	INSS_RANGE_1 = NewINSSRange(1, decimal.NewFromFloat(0.075), decimal.NewFromFloat(1212.00), decimal.NewFromFloat(1212.0), nil)
	INSS_RANGE_2 = NewINSSRange(2, decimal.NewFromFloat(0.09), decimal.NewFromFloat(1212.01), decimal.NewFromFloat(2427.35), &INSS_RANGE_1)
	INSS_RANGE_3 = NewINSSRange(3, decimal.NewFromFloat(0.12), decimal.NewFromFloat(2427.36), decimal.NewFromFloat(3641.03), &INSS_RANGE_2)
	INSS_RANGE_4 = NewINSSRange(4, decimal.NewFromFloat(0.14), decimal.NewFromFloat(3641.03), decimal.NewFromFloat(7087.22), &INSS_RANGE_3)
	INSS_RANGE_5 = NewINSSRange(5, decimal.NewFromFloat(0.14), decimal.NewFromFloat(7087.23), decimal.NewFromFloat(1000000), &INSS_RANGE_4)
)

var inssRanges = []INSSRange{
	INSS_RANGE_1, INSS_RANGE_2, INSS_RANGE_3, INSS_RANGE_4, INSS_RANGE_5,
}

func INSSRangeByGrossPay(grossPay decimal.Decimal) INSSRange {
	item := INSSRange{}
	for _, inssRange := range inssRanges {
		if grossPay.GreaterThanOrEqual(inssRange.InitValue) && grossPay.LessThanOrEqual(inssRange.EndValue) {
			item = inssRange
		}
	}
	return item
}

func (ir INSSRange) RangeDiscount() decimal.Decimal {
	if ir == INSS_RANGE_1 {
		return ir.Aliquot.Mul(ir.EndValue)
	}

	return ir.EndValue.Sub(ir.Prev.EndValue).Mul(ir.Aliquot)
}

func (ir INSSRange) PrevRangesDiscounts() decimal.Decimal {
	discount := decimal.Zero

	for _, inssRange := range inssRanges {
		if inssRange.Index < ir.Index {
			discount = discount.Add(inssRange.RangeDiscount())
		}
	}

	return discount
}

func (i INSSDiscount) Value() decimal.Decimal {
	inssRange := INSSRangeByGrossPay(i.GrossPay)

	if inssRange == INSS_RANGE_1 {
		return inssRange.RangeDiscount()
	}
	if inssRange == INSS_RANGE_5 {
		return decimal.NewFromFloat(INSS_RANGE_5_DISCOUNT_AMOUNT)
	}

	currentRangeDiscount := i.GrossPay.Sub(inssRange.Prev.EndValue).Mul(inssRange.Aliquot)

	return currentRangeDiscount.Add(inssRange.PrevRangesDiscounts()).RoundBank(2)
}

func (i INSSDiscount) Name() string {
	return "INSS"
}
