package models

import "github.com/shopspring/decimal"

func NewPercentangeDiscount(amount decimal.Decimal, percentange decimal.Decimal) *PercentangeDiscount {
	return &PercentangeDiscount{
		Amount:      amount,
		Percentange: percentange,
	}
}

type PercentangeDiscount struct {
	Amount      decimal.Decimal
	Percentange decimal.Decimal
}

func (pd PercentangeDiscount) Value() decimal.Decimal {
	return pd.Amount.Mul(pd.Percentange)
}

func (pd PercentangeDiscount) Name() string {
	return "Porcentagem"
}
