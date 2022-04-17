package models

import "github.com/shopspring/decimal"

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
