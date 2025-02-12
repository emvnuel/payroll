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

type Error struct {
	Message string `json:"message"`
}

func NewPayrollResponse(p *models.Payroll) *PayrollResponse {
	var discountsResponse []DiscountResponse

	for _, discount := range p.Discounts {
		discountsResponse = append(discountsResponse, DiscountResponse{Value: discount.Value().InexactFloat64(), Name: discount.Name()})
	}

	return &PayrollResponse{GrossPay: p.GrossPay.InexactFloat64(), NetPay: p.NetPay().InexactFloat64(), TotalDiscount: p.TotalDiscount().InexactFloat64(), Discounts: discountsResponse}
}

// @Summary Calculate Payroll
// @Description This endpoint calculates the net pay based on gross pay, number of dependents, and applied discounts.
// @Tags payroll
// @Param grossPay query number true "Gross pay of the employee" minimum(1412)
// @Param numberOfDependents query integer true "Number of dependents of the employee" minimum(0)
// @Param fixedAmountDiscount query number true "Value of the fixed amount discount" minimum(0)
// @Param percentangeDiscount query number true "Percentage discount value (between 0 and 1)" minimum(0) maximum(1)
// @Param simplifiedDeduction query boolean true "Simplified deduction" default(false)
// @Produce  json
// @Success 200 {object} controllers.PayrollResponse "Payroll information"
// @Failure 400 {object} controllers.Error "Invalid fields provided"
// @Router /payroll [get]
func GetPayroll(c *gin.Context) {

	grossPay, err1 := strconv.ParseFloat(c.Query("grossPay"), 64)
	numberOfDependents, err2 := strconv.Atoi(c.Query("numberOfDependents"))
	fixedAmountDiscountValue, err3 := strconv.ParseFloat(c.Query("fixedAmountDiscount"), 64)
	percentangeDiscountValue, err4 := strconv.ParseFloat(c.Query("percentangeDiscount"), 64)
	simplifiedDeduction, err5 := strconv.ParseBool(c.Query("simplifiedDeduction"))

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil ||
		numberOfDependents < 0 || fixedAmountDiscountValue < 0 || percentangeDiscountValue < 0 {
		c.JSON(http.StatusBadRequest, Error{Message: "Campos inválidos"})
		return
	}

	if grossPay < 1412 {
		c.JSON(http.StatusBadRequest, Error{Message: "Salário bruto deve ser maior ou igual a R$1.412,00"})
		return
	}

	if percentangeDiscountValue < 0 || percentangeDiscountValue > 1 {
		c.JSON(http.StatusBadRequest, Error{Message: "Porcentagem deve ser entre 0 e 1"})
		return
	}

	if fixedAmountDiscountValue < 0 {
		c.JSON(http.StatusBadRequest, Error{Message: "Valor fixo não pode ser negativo"})
		return
	}

	fixedDiscount := models.NewFixedAmountDiscount(decimal.NewFromFloat(fixedAmountDiscountValue))
	percentangeDiscount := models.NewPercentageDiscount(decimal.NewFromFloat(grossPay), decimal.NewFromFloat(percentangeDiscountValue))

	payroll := models.NewPayroll(decimal.NewFromFloat(grossPay), int64(numberOfDependents), simplifiedDeduction, fixedDiscount, percentangeDiscount)

	c.JSON(http.StatusOK, NewPayrollResponse(payroll))
}
