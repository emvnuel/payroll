package controllers

import (
	"net/http"
	"strconv"

	"github.com/emvnuel/payroll/models"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type PayrollResponse struct {
	GrossPay      float64            `json:"grossPay"`
	NetPay        float64            `json:"netPay"`
	TotalDiscount float64            `json:"totalDiscount"`
	Discounts     []DiscountResponse `json:"discounts"`
}

type DiscountResponse struct {
	Value float64 `json:"value"`
	Name  string  `json:"name"`
}

func NewPayrollResponse(p models.Payroll) PayrollResponse {
	var discountsResponse []DiscountResponse

	for _, discount := range p.Discounts {
		discountsResponse = append(discountsResponse, DiscountResponse{Value: discount.Value().InexactFloat64(), Name: discount.Name()})
	}

	return PayrollResponse{GrossPay: p.GrossPay.InexactFloat64(), NetPay: p.NetPay().InexactFloat64(), TotalDiscount: p.TotalDiscount().InexactFloat64(), Discounts: discountsResponse}
}

func GetPayroll(c *gin.Context) {

	grossPay, err1 := strconv.ParseFloat(c.Query("grossPay"), 64)
	numberOfDependents, err2 := strconv.Atoi(c.Query("numberOfDependents"))
	fixedAmountDiscountValue, err3 := strconv.ParseFloat(c.Query("fixedAmountDiscount"), 64)
	percentangeDiscountValue, err4 := strconv.ParseFloat(c.Query("percentangeDiscount"), 64)

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil ||
		numberOfDependents < 0 || fixedAmountDiscountValue < 0 || percentangeDiscountValue < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Campos inválidos"})
		return
	}

	if grossPay < 1212 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Salário bruto deve ser maior ou igual a R$1.212,00"})
		return
	}

	fixedDiscount := models.FixedAmountDiscount{Amount: decimal.NewFromFloat(fixedAmountDiscountValue)}
	percentangeDiscount := models.PercentangeDiscount{Amount: decimal.NewFromFloat(grossPay), Percentange: decimal.NewFromFloat(percentangeDiscountValue)}

	payroll := models.NewPayroll(decimal.NewFromFloat(grossPay), int64(numberOfDependents), fixedDiscount, percentangeDiscount)

	c.JSON(http.StatusOK, gin.H{"data": NewPayrollResponse(payroll)})
}
