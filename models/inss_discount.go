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

func NewINSSRange(index int, aliquot decimal.Decimal, initValue decimal.Decimal, endValue decimal.Decimal) *INSSRange {
	return &INSSRange{
		Index:     index,
		Aliquot:   aliquot,
		InitValue: initValue,
		EndValue:  endValue,
	}
}

var (
	INSS_RANGE_5_DISCOUNT_AMOUNT, _ = decimal.NewFromString(os.Getenv("INSS_RANGE_5_DISCOUNT_AMOUNT"))
	INSSRanges                      = loadINSSRangesFromEnv()
)

func loadINSSRangesFromEnv() []INSSRange {
	data := os.Getenv("INSS_RANGES")
	var ranges []INSSRange
	if err := json.Unmarshal([]byte(data), &ranges); err != nil {
		log.Printf("Error loading INSS ranges: %v", err)
	}
	log.Printf("Loaded INSS ranges: %+v", ranges)
	return ranges
}

func findINSSRangeByGrossPay(grossPay decimal.Decimal) INSSRange {
	for _, inssRange := range INSSRanges {
		if grossPay.GreaterThanOrEqual(inssRange.InitValue) && grossPay.LessThanOrEqual(inssRange.EndValue) {
			return inssRange
		}
	}
	return INSSRange{}
}

func (ir INSSRange) calculateRangeDiscount() decimal.Decimal {
	if ir.Index == 1 {
		return ir.Aliquot.Mul(ir.EndValue)
	}

	rangeDifference := ir.EndValue.Sub(ir.InitValue.Sub(decimal.NewFromFloat(0.01)))
	return rangeDifference.Mul(ir.Aliquot)
}

func (ir INSSRange) calculatePreviousRangesDiscount() decimal.Decimal {
	totalDiscount := decimal.Zero
	for _, inssRange := range INSSRanges {
		if inssRange.Index < ir.Index {
			totalDiscount = totalDiscount.Add(inssRange.calculateRangeDiscount())
		}
	}
	return totalDiscount
}

func (i INSSDiscount) Value() decimal.Decimal {
	inssRange := findINSSRangeByGrossPay(i.GrossPay)

	if inssRange.Index == 1 {
		return inssRange.calculateRangeDiscount()
	}

	if inssRange.Index == len(INSSRanges) {
		return INSS_RANGE_5_DISCOUNT_AMOUNT
	}

	currentRangeAmount := i.GrossPay.Sub(inssRange.InitValue.Sub(decimal.NewFromFloat(0.01)))
	currentRangeDiscount := currentRangeAmount.Mul(inssRange.Aliquot)

	return currentRangeDiscount.Add(inssRange.calculatePreviousRangesDiscount()).RoundBank(2)
}

func (i INSSDiscount) Name() string {
	return "INSS"
}
