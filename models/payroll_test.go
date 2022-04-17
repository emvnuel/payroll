package models

import (
	"math/rand"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestWhenGrossPayIs1212AndNumberOfDependetsIsAnyNetPayShouldBe1121And10Cents(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	randomInt0To100 := rand.Intn(100)
	got := NewPayroll(decimal.NewFromInt(1212), int64(randomInt0To100)).NetPay()
	want, _ := decimal.NewFromString("1121.10")

	if !got.Equal(want) {
		t.Errorf("got %q, wanted %q", got, want)
	}

}

func TestWhenGrossPayIs3000AndNumberOfDependentsIs0NetPayShouldBe2668And98Cents(t *testing.T) {
	got := NewPayroll(decimal.NewFromInt(3000), 0).NetPay()
	want, _ := decimal.NewFromString("2668.98")

	if !got.Equal(want) {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestWhenGrossPayIs4000AndNumberOfDependentsIs2NetPayShouldBe3474And92Cents(t *testing.T) {
	got := NewPayroll(decimal.NewFromInt(4000), 2).NetPay()
	want, _ := decimal.NewFromString("3474.92")

	if !got.Equal(want) {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestWhenGrossPayIs4000AndNumberOfDependentsIs0NetPayShouldBe4729And13Cents(t *testing.T) {
	got := NewPayroll(decimal.NewFromInt(6000), 0).NetPay()
	want, _ := decimal.NewFromString("4729.13")

	if !got.Equal(want) {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestWhenGrossPayIs8000AndNumberOfDependentsIs0NetPayShouldBe6068And78Cents(t *testing.T) {
	got := NewPayroll(decimal.NewFromInt(8000), 0).NetPay()
	want, _ := decimal.NewFromString("6068.78")

	if !got.Equal(want) {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestWhenGrossPayIs2500AndNumberOfDependentsIs0AndAdditionalDiscountIs20PercentNetPayShouldBe1761And98Cents(t *testing.T) {
	percentangeDiscount := PercentangeDiscount{Amount: decimal.NewFromInt(2500), Percentange: decimal.NewFromFloat(0.2)}

	got := NewPayroll(decimal.NewFromInt(2500), 0, percentangeDiscount).NetPay()
	want, _ := decimal.NewFromString("1761.98")

	if !got.Equal(want) {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestWhenGrossPayIs6500AndNumberOfDependentsIs0AndAdditionalDiscountIs40And88CentsFixedNetPayShouldBe5000(t *testing.T) {
	percentangeDiscount := FixedAmountDiscount{Amount: decimal.NewFromFloat(40.88)}

	got := NewPayroll(decimal.NewFromInt(6500), 0, percentangeDiscount).NetPay()
	want, _ := decimal.NewFromString("5000")

	if !got.Equal(want) {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
