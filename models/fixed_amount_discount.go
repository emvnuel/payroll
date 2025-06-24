package models

import "github.com/shopspring/decimal"

type FixedAmountDiscount struct {
	amount decimal.Decimal
}

func NewFixedAmountDiscount(amount decimal.Decimal) *FixedAmountDiscount {
	return &FixedAmountDiscount{
		amount: amount,
	}
}

func (fd FixedAmountDiscount) Value() decimal.Decimal {
	return fd.amount.RoundBank(2)
}

func (fd FixedAmountDiscount) Name() string {
	return "Valor fixo"
}
