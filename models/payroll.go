package models

import (
	"github.com/shopspring/decimal"
)

type Payroll struct {
	GrossPay  decimal.Decimal
	Discounts []Discount
}

func NewPayroll(grosspay decimal.Decimal, numberOfDependents int64, simplifiedDeduction bool, additionalDiscounts ...Discount) *Payroll {

	inss := NewINSSDiscount(grosspay)
	irrf := NewIRRFDiscount(grosspay, numberOfDependents, inss.Value(), simplifiedDeduction)
	payroll := Payroll{GrossPay: grosspay, Discounts: []Discount{inss, irrf}}
	payroll.Discounts = append(payroll.Discounts, additionalDiscounts...)

	return &payroll
}

func (p Payroll) NetPay() decimal.Decimal {
	return p.GrossPay.Sub(p.TotalDiscount())
}

func (p Payroll) TotalDiscount() decimal.Decimal {
	totalDiscount := decimal.Zero

	for _, discount := range p.Discounts {
		totalDiscount = totalDiscount.Add(discount.Value())
	}

	return totalDiscount
}

func (p *Payroll) AddDiscount(discount Discount) {
	p.Discounts = append(p.Discounts, discount)
}
