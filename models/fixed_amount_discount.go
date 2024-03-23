package models

import "github.com/shopspring/decimal"

func NewFixedAmountDiscount(amount decimal.Decimal) *FixedAmountDiscount {
	return &FixedAmountDiscount{
		Amount: amount,
	}
}

type FixedAmountDiscount struct {
	Amount decimal.Decimal
}

func (fd FixedAmountDiscount) Value() decimal.Decimal {
	return fd.Amount
}

func (fd FixedAmountDiscount) Name() string {
	return "Valor fixo"
}
