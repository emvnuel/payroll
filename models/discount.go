package models

import "github.com/shopspring/decimal"

type Discount interface {
	Value() decimal.Decimal
	Name() string
}
