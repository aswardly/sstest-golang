//product_test provides unit tests for business domain model of product
package product_test

import (
	"fmt"
	"os"
	"sstest/model/product"
	"testing"

	"github.com/shopspring/decimal"
)

var newProd, availableProd, discontinuedProd, outOfStockProd *product.Product

func TestMain(m *testing.M) {
	//test setup
	newProd = product.New("newProd", "New Product")

	availableProd = product.New("availableProd", "Available Product")
	availableProd.SetStatus(product.StatusAvailable)
	availableProd.SetStock(100)
	availableProd.SetPrice(decimal.New(100, 0))

	discontinuedProd = product.New("discontinuedProd", "Discontinued Product")
	discontinuedProd.SetStatus(product.StatusDiscontinued)
	discontinuedProd.SetStock(200)
	discontinuedProd.SetPrice(decimal.New(250, 0))

	outOfStockProd = product.New("outOfStockProd", "Out Of Stock Product")
	outOfStockProd.SetStatus(product.StatusAvailable)
	outOfStockProd.SetStock(0)
	outOfStockProd.SetPrice(decimal.New(60, 0))

	//run tests
	exitCode := m.Run()
	//test teardown
	os.Exit(exitCode)
}

func TestNewlyCreatedProduct(t *testing.T) {
	t.Run("Status Must Be Prototype", func(t *testing.T) {
		if product.StatusPrototype != newProd.Status() {
			t.Errorf("expected %v but got %v", product.StatusPrototype, newProd.Status())
		}
	})
	t.Run("Stock Must Be Zero", func(t *testing.T) {
		if 0 != int(newProd.Stock()) {
			t.Errorf("expected %v but got %v", 0, int(newProd.Stock()))
		}
	})
	t.Run("Price Must Be Zero", func(t *testing.T) {
		if false == newProd.Price().Equal(decimal.New(0, 0)) {
			t.Errorf("expected %v but got %v", decimal.New(0, 0).String(), newProd.Price().String())
		}
	})
}

func TestSettingUnknownStatusValue(t *testing.T) {
	_, err := newProd.SetStatus("UnknownStatus")
	if err == nil {
		t.Error("expected error but got none\n")
	}
}

func TestSetPrice(t *testing.T) {
	t.Run("Set Positive Price", func(t *testing.T) {
		result, err := newProd.SetPrice(decimal.New(100, 0))
		if err != nil {
			t.Errorf("expected %v but got %v\n", true, result)
		}
	})
	t.Run("Set Negative Price", func(t *testing.T) {
		_, err := newProd.SetPrice(decimal.New(-100, 0))
		if err == nil {
			t.Error("expected error but got none\n")
		}
	})
}

func TestCanBeOrdered(t *testing.T) {
	validProductOrder, _ := availableProd.CanBeOrdered(10)
	notEnoughStockProductOrder, _ := availableProd.CanBeOrdered(9999)
	discontinuedProductOrder, _ := discontinuedProd.CanBeOrdered(2)
	outOfStockProductOrder, _ := outOfStockProd.CanBeOrdered(10)

	var canBeOrderedTests = []struct {
		testCase string
		expected bool
		actual   bool
	}{
		{"Valid Product Order", true, validProductOrder},
		{"Not Enough Stock Product Order", false, notEnoughStockProductOrder},
		{"Discontinued Product Order", false, discontinuedProductOrder},
		{"Out of Stock Product Order", false, outOfStockProductOrder},
	}

	for _, test := range canBeOrderedTests {
		t.Run(fmt.Sprintf("%s", test.testCase), func(t *testing.T) {
			if test.expected != test.actual {
				t.Errorf("want %v got %v", test.expected, test.actual)
			}
		})
	}
}
