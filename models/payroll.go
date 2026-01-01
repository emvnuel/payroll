package models

import (
	"github.com/shopspring/decimal"
)

type Payroll struct {
	GrossPay  decimal.Decimal
	Discounts []Discount
}

func NewPayroll(grossPay decimal.Decimal, numberOfDependents int64, additionalDiscounts ...Discount) *Payroll {
	payroll := &Payroll{
		GrossPay:  grossPay,
		Discounts: make([]Discount, 0),
	}

	payroll.addMandatoryDiscounts(numberOfDependents)
	payroll.addOptionalDiscounts(additionalDiscounts...)

	return payroll
}

func (p *Payroll) addMandatoryDiscounts(numberOfDependents int64) {
	inss := NewINSSDiscount(p.GrossPay)
	irrf := NewIRRFDiscount(p.GrossPay, numberOfDependents, inss.Value())
	p.Discounts = append(p.Discounts, inss, irrf)
}

func (p *Payroll) addOptionalDiscounts(discounts ...Discount) {
	if len(discounts) > 0 {
		p.Discounts = append(p.Discounts, discounts...)
	}
}

func (p *Payroll) NetPay() decimal.Decimal {
	return p.GrossPay.Sub(p.TotalDiscount())
}

func (p *Payroll) TotalDiscount() decimal.Decimal {
	totalDiscount := decimal.Zero
	for _, discount := range p.Discounts {
		totalDiscount = totalDiscount.Add(discount.Value())
	}
	return totalDiscount
}

func (p *Payroll) AddDiscount(discount Discount) {
	p.Discounts = append(p.Discounts, discount)
}
