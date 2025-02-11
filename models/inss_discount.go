package models

import (
	"encoding/json"
	"log"
	"os"

	"github.com/shopspring/decimal"

	_ "github.com/joho/godotenv/autoload"
)

type INSSDiscount struct {
	GrossPay decimal.Decimal
}

type INSSRange struct {
	Index     int             `json:"index"`
	Aliquot   decimal.Decimal `json:"aliquot"`
	InitValue decimal.Decimal `json:"init_value"`
	EndValue  decimal.Decimal `json:"end_value"`
}

func NewINSSDiscount(grosspay decimal.Decimal) *INSSDiscount {
	return &INSSDiscount{GrossPay: grosspay}
}

func NewINSSRange(
	index int,
	aliquot decimal.Decimal,
	initValue decimal.Decimal,
	endValue decimal.Decimal) *INSSRange {

	return &INSSRange{Index: index, Aliquot: aliquot, InitValue: initValue, EndValue: endValue}
}

var (
	INSS_RANGE_5_DISCOUNT_AMOUNT, _ = decimal.NewFromString(os.Getenv("INSS_RANGE_5_DISCOUNT_AMOUNT"))
)

var INSSRanges = loadINSSRangesFromEnv()

func loadINSSRangesFromEnv() []INSSRange {
	data := os.Getenv("INSS_RANGES")
	var INSSRanges []INSSRange
	json.Unmarshal([]byte(data), &INSSRanges)
	log.Println("INSS_RANGES", INSSRanges)
	return INSSRanges
}

func inssRangeByGrossPay(grossPay decimal.Decimal) INSSRange {
	item := INSSRange{}
	for _, inssRange := range INSSRanges {
		if grossPay.GreaterThanOrEqual(inssRange.InitValue) && grossPay.LessThanOrEqual(inssRange.EndValue) {
			item = inssRange
		}
	}
	return item
}

func (ir INSSRange) rangeDiscount() decimal.Decimal {
	if ir == INSSRanges[0] {
		return ir.Aliquot.Mul(ir.EndValue)
	}
	prev := ir.InitValue.Sub(decimal.NewFromFloat(0.01))
	previousDiscounts := ir.EndValue.Sub(prev)
	return previousDiscounts.Mul(ir.Aliquot)
}

func (ir INSSRange) prevRangesDiscounts() decimal.Decimal {
	discount := decimal.Zero

	for _, inssRange := range INSSRanges {
		if inssRange.Index < ir.Index {
			discount = discount.Add(inssRange.rangeDiscount())
		}
	}

	return discount
}

func (i INSSDiscount) Value() decimal.Decimal {
	inssRange := inssRangeByGrossPay(i.GrossPay)

	if inssRange == INSSRanges[0] {
		return inssRange.rangeDiscount()
	}
	if inssRange == INSSRanges[len(INSSRanges)-1] {
		return INSS_RANGE_5_DISCOUNT_AMOUNT
	}

	prev := inssRange.InitValue.Sub(decimal.NewFromFloat(0.01))
	previousDiscounts := i.GrossPay.Sub(prev)
	currentRangeDiscount := previousDiscounts.Mul(inssRange.Aliquot)
	return currentRangeDiscount.Add(inssRange.prevRangesDiscounts()).RoundBank(2)
}

func (i INSSDiscount) Name() string {
	return "INSS"
}
