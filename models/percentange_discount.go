package models

import "github.com/shopspring/decimal"

type PercentageDiscount struct {
	amount     decimal.Decimal
	percentage decimal.Decimal
}

func NewPercentageDiscount(amount decimal.Decimal, percentage decimal.Decimal) *PercentageDiscount {
	return &PercentageDiscount{
		amount:     amount,
		percentage: percentage,
	}
}

func (pd PercentageDiscount) Value() decimal.Decimal {
	return pd.amount.Mul(pd.percentage).RoundBank(2)
}

func (pd PercentageDiscount) Name() string {
	return "Porcentagem"
}
