package models

import "github.com/shopspring/decimal"

type FixedAmountDiscount struct {
	Amount decimal.Decimal
}

func (fd FixedAmountDiscount) Value() decimal.Decimal {
	return fd.Amount
}

func (fd FixedAmountDiscount) Name() string {
	return "Valor fixo"
}
