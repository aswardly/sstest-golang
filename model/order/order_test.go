//order_test provides unit tests for business domain model of order and order item
package order_test

import (
	"fmt"
	"os"
	"sstest/model/coupon"
	"sstest/model/order"
	"sstest/model/product"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

var newOrder, emptyItemDraftOrder, validItemDraftOrder *order.Order
var newItem *order.Item
var newProd, availableProd, discontinuedProd, outOfStockProd *product.Product

var activeCoupon, inactiveCoupon, suspendedCoupon, noStockCoupon, earlyCoupon, expiredCoupon *coupon.Coupon
var currentTime = time.Now()
var beginningOfToday = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
var beginningOfTomorrow = beginningOfToday.AddDate(0, 0, 1)
var beginningOf90DaysFromToday = beginningOfToday.AddDate(0, 0, 90)
var beginningOfYesterday = beginningOfToday.AddDate(0, -1, 0)
var beginningOfLastWeek = beginningOfToday.AddDate(0, -7, 0)

func TestMain(m *testing.M) {
	//test setup for new order and new order item
	newProd = product.New("newProd", "New Product")
	newOrder = order.New("newOrder")
	newItem = order.NewItem("newItem", newOrder, newProd)

	newItem.SetOrder(newOrder)
	newItem.SetProduct(newProd)

	//test setup for coupon
	activeCoupon = coupon.New("activeCoupon")
	activeCoupon.SetStatus(coupon.StatusActive)
	activeCoupon.SetStock(100)
	activeCoupon.SetKind(coupon.KindValue)
	activeCoupon.SetValue(decimal.New(10000, 0))

	inactiveCoupon = coupon.New("inactiveCoupon")
	inactiveCoupon.SetStatus(coupon.StatusInactive)
	inactiveCoupon.SetStock(50)
	inactiveCoupon.SetKind(coupon.KindPercentage)
	inactiveCoupon.SetValue(decimal.New(50, 0))

	suspendedCoupon = coupon.New("suspendedCoupon")
	suspendedCoupon.SetStatus(coupon.StatusSuspended)
	suspendedCoupon.SetStock(10)
	suspendedCoupon.SetKind(coupon.KindValue)
	suspendedCoupon.SetValue(decimal.New(25, 0))

	noStockCoupon = coupon.New("noStockCoupon")
	noStockCoupon.SetStatus(coupon.StatusActive)
	noStockCoupon.SetStock(0)
	noStockCoupon.SetKind(coupon.KindPercentage)
	noStockCoupon.SetValue(decimal.New(30, 0))

	earlyCoupon = coupon.New("earlyCoupon")
	earlyCoupon.SetStatus(coupon.StatusActive)
	earlyCoupon.SetStock(20)
	earlyCoupon.SetStartDate(beginningOfTomorrow)
	earlyCoupon.SetEndDate(beginningOf90DaysFromToday)
	earlyCoupon.SetKind(coupon.KindValue)
	earlyCoupon.SetValue(decimal.New(10000, 0))

	expiredCoupon = coupon.New("expiredCoupon")
	expiredCoupon.SetStatus(coupon.StatusActive)
	expiredCoupon.SetStock(20)
	expiredCoupon.SetStartDate(beginningOfLastWeek)
	expiredCoupon.SetEndDate(beginningOfYesterday)
	expiredCoupon.SetKind(coupon.KindPercentage)
	expiredCoupon.SetValue(decimal.New(80, 0))

	//run tests
	exitCode := m.Run()
	//test teardown
	os.Exit(exitCode)
}

func TestNewlyCreatedOrder(t *testing.T) {
	var now = time.Now()

	t.Run("Items must be empty", func(t *testing.T) {
		if 0 != newItem.Quantity() {
			t.Errorf("expected %v but got %v", 0, newItem.Quantity())
		}
	})
	t.Run("Status must be draft", func(t *testing.T) {
		if order.StatusDraft != newOrder.Status() {
			t.Errorf("expected %v but got %v", order.StatusDraft, newOrder.Status())
		}
	})
	t.Run("Ship Status must be None", func(t *testing.T) {
		if order.ShipStatusNone != newOrder.ShippingStatus() {
			t.Errorf("expected %v but got %v", order.ShipStatusNone, newOrder.ShippingStatus())
		}
	})
	t.Run("Created Date must be now", func(t *testing.T) {
		if false == now.Equal(newOrder.CreatedDate()) {
			t.Errorf("expected %v but got %v", now.String(), newOrder.CreatedDate().String())
		}
	})
	t.Run("Submitted Date must be unix epoch", func(t *testing.T) {
		if false == time.Unix(0, 0).Equal(newOrder.SubmittedDate()) {
			t.Errorf("expected %v but got %v", now.String(), newOrder.SubmittedDate().String())
		}
	})
	t.Run("Processed Date must be unix epoch", func(t *testing.T) {
		if false == time.Unix(0, 0).Equal(newOrder.ProcessedDate()) {
			t.Errorf("expected %v but got %v", now.String(), newOrder.ProcessedDate().String())
		}
	})
}

func TestAddProduct(t *testing.T) {
	newOrder = order.New("newOrder")

	//test setup for product
	availableProd = product.New("availableProd", "Available Product")
	availableProd.SetStatus(product.StatusAvailable)
	availableProd.SetStock(100)

	discontinuedProd = product.New("discontinuedProd", "Discontinued Product")
	discontinuedProd.SetStatus(product.StatusDiscontinued)
	discontinuedProd.SetStock(200)

	outOfStockProd = product.New("outOfStockProd", "Out Of Stock Product")
	outOfStockProd.SetStatus(product.StatusAvailable)
	outOfStockProd.SetStock(0)

	addAvailableProdResult, errAddAvailableProdresult := newOrder.AddProduct(availableProd, 10)
	addNotEnoughStockAvailableProdResult, errAddNotEnoughStockAvailableProdresult := newOrder.AddProduct(availableProd, 91)
	addDiscontinuedProdResult, errAddDiscontinuedProdresult := newOrder.AddProduct(discontinuedProd, 10)
	addNoStockProdResult, errAddNoStockProdresult := newOrder.AddProduct(outOfStockProd, 5)

	var addProductTests = []struct {
		testCase           string
		expectedValue      bool
		actualValue        bool
		expectedErrIsNil   bool
		actualErrIsNil     bool
		expectedHasProduct bool
		actualhasProduct   bool
	}{
		{"Add Available Product", true, addAvailableProdResult, true, errAddAvailableProdresult == nil, true, newOrder.HasProduct(availableProd)},
		{"Add Not Enough Stock Available Product", false, addNotEnoughStockAvailableProdResult, false, errAddNotEnoughStockAvailableProdresult == nil, true, newOrder.HasProduct(availableProd)},
		{"Add Discontinued Product Result", false, addDiscontinuedProdResult, false, errAddDiscontinuedProdresult == nil, true, newOrder.HasProduct(discontinuedProd)},
		{"Add No Stock Product Result", false, addNoStockProdResult, false, errAddNoStockProdresult == nil, true, newOrder.HasProduct(outOfStockProd)},
	}

	for _, test := range addProductTests {
		t.Run(fmt.Sprintf("%s", test.testCase), func(t *testing.T) {
			if test.expectedValue != test.actualValue {
				t.Errorf("want %v for value, got %v", test.expectedValue, test.actualValue)
			}
			if test.expectedErrIsNil != test.actualErrIsNil {
				t.Errorf("want %v for error, got %v", test.expectedErrIsNil, test.actualErrIsNil)
			}
		})
	}
}

func TestEditProduct(t *testing.T) {
	newOrder = order.New("newOrder")

	//test setup for product
	availableProd = product.New("availableProd", "Available Product")
	availableProd.SetStatus(product.StatusAvailable)
	availableProd.SetStock(100)

	anotherAvailableProd := product.New("anotherAvailableProd", "Another Available Product")
	anotherAvailableProd.SetStatus(product.StatusAvailable)
	anotherAvailableProd.SetStock(10)

	discontinuedProd = product.New("discontinuedProd", "Discontinued Product")
	discontinuedProd.SetStatus(product.StatusDiscontinued)
	discontinuedProd.SetStock(200)

	newOrder.AddProduct(availableProd, 5)
	newOrder.AddProduct(anotherAvailableProd, 8)

	validEditAvailableProdResult, errValidEditAvailableProd := newOrder.EditProduct(availableProd, 20)
	notEnoughStockEditAvailableProdResult, errNotEnoughStockValidEditAvailableProd := newOrder.EditProduct(anotherAvailableProd, 300)
	editNonExistingProdResult, errEditNonExistingProdResult := newOrder.AddProduct(availableProd, 91)

	var editProductTests = []struct {
		testCase           string
		expectedValue      bool
		actualValue        bool
		expectedErrIsNil   bool
		actualErrIsNil     bool
		expectedHasProduct bool
		actualhasProduct   bool
	}{
		{"Valid Edit Available Product", true, validEditAvailableProdResult, true, errValidEditAvailableProd == nil, true, newOrder.HasProduct(availableProd)},
		{"Add Not Enough Stock Available Product", false, notEnoughStockEditAvailableProdResult, false, errNotEnoughStockValidEditAvailableProd == nil, true, newOrder.HasProduct(availableProd)},
		{"Add Discontinued Product Result", false, editNonExistingProdResult, false, errEditNonExistingProdResult == nil, false, newOrder.HasProduct(discontinuedProd)},
	}

	for _, test := range editProductTests {
		t.Run(fmt.Sprintf("%s", test.testCase), func(t *testing.T) {
			if test.expectedValue != test.actualValue {
				t.Errorf("want %v for value, got %v", test.expectedValue, test.actualValue)
			}
			if test.expectedErrIsNil != test.actualErrIsNil {
				t.Errorf("want %v for error, got %v", test.expectedErrIsNil, test.actualErrIsNil)
			}
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	newOrder = order.New("newOrder")

	//test setup for product
	availableProd = product.New("availableProd", "Available Product")
	availableProd.SetStatus(product.StatusAvailable)
	availableProd.SetStock(100)

	discontinuedProd = product.New("discontinuedProd", "Discontinued Product")
	discontinuedProd.SetStatus(product.StatusDiscontinued)
	discontinuedProd.SetStock(200)

	newOrder.AddProduct(availableProd, 5)

	validDeleteProdResult, errValidDeleteProd := newOrder.DeleteProduct(availableProd)
	deleteNonExistingProdResult, errDeleteNonExistingProdResult := newOrder.DeleteProduct(discontinuedProd)

	var deleteProductTests = []struct {
		testCase           string
		expectedValue      bool
		actualValue        bool
		expectedErrIsNil   bool
		actualErrIsNil     bool
		expectedHasProduct bool
		actualHasProduct   bool
	}{
		{"Valid Delete Product", true, validDeleteProdResult, true, errValidDeleteProd == nil, false, newOrder.HasProduct(availableProd)},
		{"Delete Non-Existing Product", false, deleteNonExistingProdResult, false, errDeleteNonExistingProdResult == nil, false, newOrder.HasProduct(discontinuedProd)},
	}

	for _, test := range deleteProductTests {
		t.Run(fmt.Sprintf("%s", test.testCase), func(t *testing.T) {
			if test.expectedValue != test.actualValue {
				t.Errorf("want %v for value, got %v", test.expectedValue, test.actualValue)
			}
			if test.expectedErrIsNil != test.actualErrIsNil {
				t.Errorf("want %v for error, got %v", test.expectedErrIsNil, test.actualErrIsNil)
			}
			if test.expectedHasProduct != test.actualHasProduct {
				t.Errorf("want %v for error, got %v", test.expectedHasProduct, test.actualHasProduct)
			}
		})
	}
}

func TestSubmitOrder(t *testing.T) {
	//test setup for product
	availableProd = product.New("availableProd", "Available Product")
	availableProd.SetStatus(product.StatusAvailable)
	availableProd.SetStock(100)
	availableProd.SetPrice(decimal.New(100, 0))

	anotherAvailableProd := product.New("anotherAvailableProd", "Another Available Product")
	anotherAvailableProd.SetStatus(product.StatusAvailable)
	anotherAvailableProd.SetStock(50)
	anotherAvailableProd.SetPrice(decimal.New(150, 0))

	//test setup for coupon
	activeValueCoupon := coupon.New("activeCoupon")
	activeValueCoupon.SetStatus(coupon.StatusActive)
	activeValueCoupon.SetStock(100)
	activeValueCoupon.SetKind(coupon.KindValue)
	activeValueCoupon.SetValue(decimal.New(100, 0))

	activePercentageCoupon := coupon.New("activePercentage")
	activePercentageCoupon.SetStatus(coupon.StatusActive)
	activePercentageCoupon.SetStock(10)
	activePercentageCoupon.SetKind(coupon.KindPercentage)
	activePercentageCoupon.SetValue(decimal.New(20, 0))

	inactiveCoupon = coupon.New("inactiveCoupon")
	inactiveCoupon.SetStatus(coupon.StatusInactive)
	inactiveCoupon.SetStock(50)
	inactiveCoupon.SetKind(coupon.KindPercentage)
	inactiveCoupon.SetValue(decimal.New(50, 0))

	//test setup for order
	validOrder1 := order.New("validOrder1")
	validOrder1.AddProduct(availableProd, 5)
	validOrder1.AddProduct(anotherAvailableProd, 5)

	validOrder2 := order.New("validOrder2")
	validOrder2.AddProduct(availableProd, 5)
	validOrder2.AddProduct(anotherAvailableProd, 5)

	validOrder3 := order.New("validOrder3")
	validOrder3.AddProduct(availableProd, 5)
	validOrder3.AddProduct(anotherAvailableProd, 5)

	validOrder4 := order.New("validOrder4")
	validOrder4.AddProduct(availableProd, 5)
	validOrder4.AddProduct(anotherAvailableProd, 5)

	noItemOrder := order.New("noItemOrder")

	validOrderValueCouponSubmitResult, errValidOrderValueCouponSubmit := validOrder1.Submit("ship name1", "ship address1", activeValueCoupon)
	validOrderPercentageCouponSubmitResult, errValidOrderPercentageCouponSubmit := validOrder2.Submit("ship name2", "ship address2", activePercentageCoupon)
	validOrderInactiveCouponSubmitResult, errValidOrderInactiveCouponSubmit := validOrder3.Submit("ship name3", "ship address3", inactiveCoupon)
	validOrderNoCouponSubmitResult, errValidOrderNoCouponSubmit := validOrder4.Submit("ship name4", "ship address4", nil)
	noItemOrderValueCouponSubmitResult, errNoItemOrderValueCouponSubmit := noItemOrder.Submit("ship name3", "ship address3", activeValueCoupon)
	noItemOrderInactiveCouponSubmitResult, errNoItemOrderInactiveCouponSubmit := noItemOrder.Submit("ship name3", "ship address3", activeValueCoupon)

	var SubmitOrderTests = []struct {
		testCase              string
		expectedValue         bool
		actualValue           bool
		expectedErrIsNil      bool
		actualErrIsNil        bool
		expectedOrderStatus   string
		actualOrderStatus     string
		expectedOrderAmount   decimal.Decimal
		actualOrderAmount     decimal.Decimal
		expectedCouponStock   int64
		actualCouponStock     int64
		expectedProduct1Stock int64
		actualProduct1Stock   int64
		expectedProduct2Stock int64
		actualProduct2Stock   int64
	}{
		{"Valid Order Value Coupon Submit", true, validOrderValueCouponSubmitResult, true, errValidOrderValueCouponSubmit == nil, "S", validOrder1.Status(), decimal.New(1150, 0), validOrder1.Amount(), 99, activeValueCoupon.Stock(), 95, availableProd.Stock(), anotherAvailableProd.Stock(), 45},
		{"Valid Order Percentage Coupon Submit", true, validOrderPercentageCouponSubmitResult, true, errValidOrderPercentageCouponSubmit == nil, "S", validOrder2.Status(), decimal.New(1000, 0), validOrder2.Amount(), 9, activePercentageCoupon.Stock(), 90, availableProd.Stock(), anotherAvailableProd.Stock(), 40},
		{"Valid Order Inactive Coupon Submit", false, validOrderInactiveCouponSubmitResult, false, errValidOrderInactiveCouponSubmit == nil, "D", validOrder3.Status(), decimal.New(0, 0), validOrder3.Amount(), 50, inactiveCoupon.Stock(), 90, availableProd.Stock(), anotherAvailableProd.Stock(), 40},
		{"Valid Order No Coupon Submit", true, validOrderNoCouponSubmitResult, true, errValidOrderNoCouponSubmit == nil, "S", validOrder4.Status(), decimal.New(1250, 0), validOrder4.Amount(), 0, 0, 90, availableProd.Stock(), anotherAvailableProd.Stock(), 40},
		{"No Item Order Value Coupon Submit", false, noItemOrderValueCouponSubmitResult, false, errNoItemOrderValueCouponSubmit == nil, "D", noItemOrder.Status(), decimal.New(0, 0), noItemOrder.Amount(), 99, activeValueCoupon.Stock(), 90, availableProd.Stock(), anotherAvailableProd.Stock(), 40},
		{"No Item Order Inactive Coupon Submit", false, noItemOrderInactiveCouponSubmitResult, false, errNoItemOrderInactiveCouponSubmit == nil, "D", noItemOrder.Status(), decimal.New(0, 0), noItemOrder.Amount(), 50, inactiveCoupon.Stock(), 90, availableProd.Stock(), anotherAvailableProd.Stock(), 40},
	}

	for _, test := range SubmitOrderTests {
		t.Run(fmt.Sprintf("%s", test.testCase), func(t *testing.T) {
			if test.expectedValue != test.actualValue {
				t.Errorf("want %v for value, got %v", test.expectedValue, test.actualValue)
			}
			if test.expectedErrIsNil != test.actualErrIsNil {
				t.Errorf("want %v for error, got %v", test.expectedErrIsNil, test.actualErrIsNil)
			}
			if false == test.expectedOrderAmount.Equal(test.actualOrderAmount) {
				t.Errorf("want %v for amount, got %v", test.expectedOrderAmount.String(), test.actualOrderAmount.String())
			}
			if test.expectedCouponStock != test.actualCouponStock {
				t.Errorf("want %v for coupon stock, got %v", test.expectedCouponStock, test.actualCouponStock)
			}
		})
	}
}

func TestProcessOrder(t *testing.T) {
	//setup orders
	draftOrder := order.New("draftOrder")

	submittedOrder := order.New("submittedOrder")
	submittedOrder.SetStatus(order.StatusSubmitted)

	//sleep for 100 ms, so at least order processed time and order creation time will be 100ms apart
	time.Sleep(100 * time.Millisecond)

	draftOrderProcessResult, errDraftOrderProcess := draftOrder.Process()
	submittedOrderProcessResult, errSubmittedOrderProcess := submittedOrder.Process()

	var processOrderTests = []struct {
		testCase         string
		expectedValue    bool
		actualValue      bool
		expectedErrIsNil bool
		actualErrIsNil   bool
	}{
		{"Draft Order Process", false, draftOrderProcessResult, false, errDraftOrderProcess == nil},
		{"Submitted Order Process", true, submittedOrderProcessResult, true, errSubmittedOrderProcess == nil},
	}

	for _, test := range processOrderTests {
		t.Run(fmt.Sprintf("%s", test.testCase), func(t *testing.T) {
			if test.expectedValue != test.actualValue {
				t.Errorf("want %v for value, got %v", test.expectedValue, test.actualValue)
			}
			if test.expectedErrIsNil != test.actualErrIsNil {
				t.Errorf("want %v for error, got %v", test.expectedErrIsNil, test.actualErrIsNil)
			}
		})
	}

	createdTimeNano := submittedOrder.CreatedDate().UnixNano()
	processedTimeNano := submittedOrder.ProcessedDate().UnixNano()

	t.Run("Submitted Time must be later than Created Time", func(t *testing.T) {
		if processedTimeNano <= createdTimeNano {
			t.Errorf("Submitted time nanosecond %v must be later than Created time nanosecond %v\n", processedTimeNano, createdTimeNano)
		}
	})
}

func TestCancelOrder(t *testing.T) {
	//setup orders
	draftOrder := order.New("draftOrder")

	submittedOrder := order.New("submittedOrder")
	submittedOrder.SetStatus(order.StatusSubmitted)

	draftOrderCancelResult, errDraftOrderCancel := draftOrder.Cancel()
	submittedOrderCancelResult, errSubmittedOrderCancel := submittedOrder.Cancel()

	var cancelOrderTests = []struct {
		testCase         string
		expectedValue    bool
		actualValue      bool
		expectedErrIsNil bool
		actualErrIsNil   bool
	}{
		{"Draft Order Cancel", false, draftOrderCancelResult, false, errDraftOrderCancel == nil},
		{"Submitted Order Cancel", true, submittedOrderCancelResult, true, errSubmittedOrderCancel == nil},
	}

	for _, test := range cancelOrderTests {
		t.Run(fmt.Sprintf("%s", test.testCase), func(t *testing.T) {
			if test.expectedValue != test.actualValue {
				t.Errorf("want %v for value, got %v", test.expectedValue, test.actualValue)
			}
			if test.expectedErrIsNil != test.actualErrIsNil {
				t.Errorf("want %v for error, got %v", test.expectedErrIsNil, test.actualErrIsNil)
			}
		})
	}
}

func TestProcessShipping(t *testing.T) {
	//setup orders
	draftOrder := order.New("draftOrder")

	processedOrder := order.New("processedOrder")
	processedOrder.SetStatus(order.StatusProcessed)

	draftOrderProcessShippingResult, errDraftOrderProcessShipping := draftOrder.ProcessShipping("dummyTrackingNo")
	processedOrderProcessShippingResult, errProcessedOrderProcessShipping := processedOrder.ProcessShipping("dummyTrackingNo")

	var processShippingOrderTests = []struct {
		testCase           string
		expectedValue      bool
		actualValue        bool
		expectedErrIsNil   bool
		actualErrIsNil     bool
		expectedShipStatus string
		actualShipStatus   string
		expectedTrackingNo string
		actualTrackingNo   string
	}{
		{"Draft Order Process Shipping", false, draftOrderProcessShippingResult, false, errDraftOrderProcessShipping == nil, order.ShipStatusNone, draftOrder.ShippingStatus(), "", draftOrder.ShippingTrackingID()},
		{"Processed Order Process Shipping", true, processedOrderProcessShippingResult, true, errProcessedOrderProcessShipping == nil, order.ShipStatusOnProcess, processedOrder.ShippingStatus(), "dummyTrackingNo", processedOrder.ShippingTrackingID()},
	}

	for _, test := range processShippingOrderTests {
		t.Run(fmt.Sprintf("%s", test.testCase), func(t *testing.T) {
			if test.expectedValue != test.actualValue {
				t.Errorf("want %v for value, got %v", test.expectedValue, test.actualValue)
			}
			if test.expectedErrIsNil != test.actualErrIsNil {
				t.Errorf("want %v for error, got %v", test.expectedErrIsNil, test.actualErrIsNil)
			}
			if test.expectedShipStatus != test.actualShipStatus {
				t.Errorf("want %v for ship status, got %v", test.expectedShipStatus, test.actualShipStatus)
			}
			if test.expectedTrackingNo != test.actualTrackingNo {
				t.Errorf("want %v for tracking no, got %v", test.expectedTrackingNo, test.actualTrackingNo)
			}
		})
	}
}

func TestFinishOrder(t *testing.T) {
	//setup orders
	draftOrder := order.New("draftOrder")

	processedOrder := order.New("processedOrder")
	processedOrder.SetStatus(order.StatusProcessed)
	processedOrder.SetShippingTrackingID("dummyTrackingId")
	processedOrder.SetShippingStatus(order.ShipStatusOnProcess)

	draftOrderFinishResult, errDraftOrderFinish := draftOrder.FinishOrder()
	processedOrderFinishResult, errProcessedOrderFinish := processedOrder.FinishOrder()

	var processShippingOrderTests = []struct {
		testCase           string
		expectedValue      bool
		actualValue        bool
		expectedErrIsNil   bool
		actualErrIsNil     bool
		expectedStatus     string
		actualStatus       string
		expectedShipStatus string
		actualShipStatus   string
	}{
		{"Draft Order Finish", false, draftOrderFinishResult, false, errDraftOrderFinish == nil, order.StatusDraft, draftOrder.Status(), order.ShipStatusNone, draftOrder.ShippingStatus()},
		{"Processed Order Finish", true, processedOrderFinishResult, true, errProcessedOrderFinish == nil, order.StatusDelivered, processedOrder.Status(), order.ShipStatusDelivered, processedOrder.ShippingStatus()},
	}

	for _, test := range processShippingOrderTests {
		t.Run(fmt.Sprintf("%s", test.testCase), func(t *testing.T) {
			if test.expectedValue != test.actualValue {
				t.Errorf("want %v for value, got %v", test.expectedValue, test.actualValue)
			}
			if test.expectedErrIsNil != test.actualErrIsNil {
				t.Errorf("want %v for error, got %v", test.expectedErrIsNil, test.actualErrIsNil)
			}
			if test.expectedStatus != test.actualStatus {
				t.Errorf("want %v for status, got %v", test.expectedShipStatus, test.actualShipStatus)
			}
			if test.expectedShipStatus != test.actualShipStatus {
				t.Errorf("want %v for ship status, got %v", test.expectedShipStatus, test.actualShipStatus)
			}
		})
	}
}

func TestNewlyCreatedItem(t *testing.T) {
	t.Run("Quantity Must Be Zero", func(t *testing.T) {
		if 0 != newItem.Quantity() {
			t.Errorf("expected %v but got %v", 0, newItem.Quantity())
		}
	})
}

func TestAddQuantity(t *testing.T) {
	t.Run("Add Quantity by 100", func(t *testing.T) {
		newItem.AddQuantity(100)
		if 100 != newItem.Quantity() {
			t.Errorf("expected %v but got %v", 0, newItem.Quantity())
		}
	})
	t.Run("Add Quantity by -50", func(t *testing.T) {
		newItem.AddQuantity(-50)
		if 50 != newItem.Quantity() {
			t.Errorf("expected %v but got %v", 0, newItem.Quantity())
		}
	})
	t.Run("Subtract Quantity by 25", func(t *testing.T) {
		newItem.SubtractQuantity(25)
		if 25 != newItem.Quantity() {
			t.Errorf("expected %v but got %v", 0, newItem.Quantity())
		}
	})
}
