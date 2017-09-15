//Package order provides the business domain models definitions of order and order item
package order

import (
	"fmt"
	"sstest/model/coupon"
	"sstest/model/product"
	"sync"
	"time"

	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

//StatusDraft is const for 'draft' order status
const StatusDraft string = "D"

//StatusSubmitted is const for 'submitted' order status
const StatusSubmitted string = "S"

//StatusCanceled is const for 'canceled' order status
const StatusCanceled string = "C"

//StatusProcessed is const for 'processed' order status
const StatusProcessed string = "P"

//StatusDelivered is const for 'delivered' order status
const StatusDelivered string = "DL"

//statusMap is a map of known status code and its label pairs
var statusMap = map[string]string{
	StatusDraft:     "Draft",
	StatusSubmitted: "Submitted",
	StatusCanceled:  "Canceled",
	StatusProcessed: "Processed",
	StatusDelivered: "Delivered",
}

//ShipStatusNone is const for 'none' order ship status
const ShipStatusNone string = "N"

//ShipStatusOnProcess is const for 'on process' order ship status
const ShipStatusOnProcess string = "O"

//ShipStatusDelivered is const for 'on process' order ship status
const ShipStatusDelivered string = "D"

//shipstatusMap is a map of known shipping status code and its label pairs
var shipstatusMap = map[string]string{
	ShipStatusNone:      "None",
	ShipStatusOnProcess: "On Process",
	ShipStatusDelivered: "Delivered",
}

//Order is business domain model definition of order
type Order struct {
	id                 string
	createdDate        time.Time
	submittedDate      time.Time
	processedDate      time.Time
	status             string
	items              map[string]*Item
	coupon             *coupon.Coupon
	amount             decimal.Decimal
	shippingName       string
	shippingAddress    string
	shippingStatus     string
	shippingTrackingID string
	mu                 sync.Mutex
}

//New creates a new product model struct, initializes it's properties and returns a reference to it
func New(id string) *Order {
	return &Order{
		id,
		time.Now(),
		time.Unix(0, 0),
		time.Unix(0, 0),
		StatusDraft,
		make(map[string]*Item, 5),
		nil,
		decimal.New(0, 0),
		"",
		"",
		ShipStatusNone,
		"",
		*new(sync.Mutex),
	}
}

//ID is a getter function for returning an order's id
func (o *Order) ID() string {
	return o.id
}

//CreatedDate is a getter function for returning an order's created date
func (o *Order) CreatedDate() time.Time {
	return o.createdDate
}

//SubmittedDate is a getter function for returning an order's submitted date
func (o *Order) SubmittedDate() time.Time {
	return o.submittedDate
}

//ProcessedDate is a getter function for returning an order's processed date
func (o *Order) ProcessedDate() time.Time {
	return o.processedDate
}

//Status is a getter function for returning an order's status
func (o *Order) Status() string {
	return o.status
}

//Items is a getter function for returning an order's items
func (o *Order) Items() map[string]*Item {
	return o.items
}

//Amount is a getter function for returning an order's amount
func (o *Order) Amount() decimal.Decimal {
	return o.amount
}

//ShippingName is a getter function for returning an order's shipping name
func (o *Order) ShippingName() string {
	return o.shippingName
}

//ShippingAddress is a getter function for returning an order's shipping address
func (o *Order) ShippingAddress() string {
	return o.shippingAddress
}

//ShippingStatus is a getter function for returning an order's shipping status
func (o *Order) ShippingStatus() string {
	return o.shippingStatus
}

//ShippingTrackingID is a getter function for returning an order's shipping tracking id
func (o *Order) ShippingTrackingID() string {
	return o.shippingTrackingID
}

//SetID is a setter function for setting an order's id
func (o *Order) SetID(id string) *Order {
	o.id = id
	return o
}

//SetCreatedDate is a setter function for setting an order's created date
func (o *Order) SetCreatedDate(createdDate time.Time) *Order {
	o.createdDate = createdDate
	return o
}

//SetSubmittedDate is a setter function for setting an order's submitted date
func (o *Order) SetSubmittedDate(submittedDate time.Time) *Order {
	o.submittedDate = submittedDate
	return o
}

//SetProcessedDate is a setter function for setting an order's processed date
func (o *Order) SetProcessedDate(processedDate time.Time) *Order {
	o.processedDate = processedDate
	return o
}

//SetStatus is a setter function for setting an order's status
func (o *Order) SetStatus(status string) (*Order, *errors.Error) {
	if _, ok := statusMap[status]; false == ok {
		//note: defensive code, return nil when error is encountered (possible runtime error on caller code when method chaining)
		return nil, errors.Wrap(fmt.Errorf("Can't set unknown status type: %v", status), 0)
	}
	o.status = status
	return o, nil
}

//SetCoupon is a setter function for setting an order's coupon
func (o *Order) SetCoupon(coupon *coupon.Coupon) *Order {
	o.coupon = coupon
	return o
}

//SetAmount is a setter function for setting an order's amount
func (o *Order) SetAmount(amount decimal.Decimal) *Order {
	o.amount = amount
	return o
}

//SetShippingName is a setter function for setting an order's shipping name
func (o *Order) SetShippingName(name string) *Order {
	o.shippingName = name
	return o
}

//SetShippingAddress is a setter function for setting an order's shipping name
func (o *Order) SetShippingAddress(address string) *Order {
	o.shippingAddress = address
	return o
}

//SetShippingStatus is a setter function for setting an order's shipping status
func (o *Order) SetShippingStatus(status string) (*Order, *errors.Error) {
	if _, ok := shipstatusMap[status]; false == ok {
		//note: defensive code, return nil when error is encountered (possible runtime error on caller code when method chaining)
		return nil, errors.Wrap(fmt.Errorf("Can't set unknown shipping status type: %v", status), 0)
	}
	o.shippingStatus = status
	return o, nil
}

//SetShippingTrackingID is a setter function for setting an order's shipping tracking id
func (o *Order) SetShippingTrackingID(trackNo string) *Order {
	o.shippingTrackingID = trackNo
	return o
}

//Business logic methods

//AddProduct is a function for adding a product to an order (as order item) with a specified quantity for the purpose of ordering
//Returns true if product addition is successful or false and an error describing the failure
func (o *Order) AddProduct(product *product.Product, quantity int) (bool, *errors.Error) {
	if StatusDraft != o.status {
		return false, errors.Wrap(fmt.Errorf("Can't add product to order %v, status is %v (not draft)", o.id, statusMap[o.status]), 0)
	}
	if false == o.HasProduct(product) {
		//product doesn't exist in order item
		if ok, err := product.CanBeOrdered(quantity); false == ok {
			return false, errors.Wrap(fmt.Errorf("Can't add product %v to order %v with quantity %d: %v", product.ID(), o.id, quantity, err), 0)
		}
		newItem := NewItem(uuid.New().String(), o, product)
		newItem.quantity = quantity
		o.items[product.ID()] = newItem
	} else {
		//product exists in order item
		existingItem := o.items[product.ID()]
		if ok, err := product.CanBeOrdered(existingItem.quantity + quantity); false == ok {
			return false, errors.Wrap(fmt.Errorf("Can't order product %v with quantity %d: %v", product.ID(), existingItem.quantity+quantity, err), 0)
		}
		o.items[product.ID()].AddQuantity(quantity)
	}
	return true, nil
}

//EditProduct is a function for editing a product in an order (as order item) to a specified quantity for the purpose of ordering
//Returns true if product editing is successful or false and an error describing the failure
func (o *Order) EditProduct(product *product.Product, quantity int) (bool, *errors.Error) {
	if StatusDraft != o.status {
		return false, errors.Wrap(fmt.Errorf("Can't edit product in order %v, status is %v (not draft)", o.id, statusMap[o.status]), 0)
	}
	if false == o.HasProduct(product) {
		return false, errors.Wrap(fmt.Errorf("Can't edit, order %v has no product with id: %v", o.id, product.ID()), 0)
	}
	//else product exists in order
	existingItem := o.items[product.ID()]
	if ok, err := product.CanBeOrdered(existingItem.quantity + quantity); false == ok {
		return false, errors.Wrap(fmt.Errorf("Can't edit, can't order product %v with quantity %d: %v", product.ID(), existingItem.quantity+quantity, err), 0)
	}
	existingItem.quantity += quantity
	return true, nil
}

//DeleteProduct is a function for removing a product from an order (as order item) for the purpose of ordering
//Returns true if product deletion is successful or false and an error describing the failure
func (o *Order) DeleteProduct(product *product.Product) (bool, *errors.Error) {
	if StatusDraft != o.status {
		return false, errors.Wrap(fmt.Errorf("Can't delete product in order %v, status is %v (not draft)", o.id, statusMap[o.status]), 0)
	}
	if false == o.HasProduct(product) {
		return false, errors.Wrap(fmt.Errorf("Can't delete, order %v has no product with id: %v", o.id, product.ID()), 0)
	}
	delete(o.items, product.ID())
	return true, nil
}

//HasProduct is a function for checking whether order has a given product (as order item)
//Returns true if product exists in order items, otherwise returns false
func (o *Order) HasProduct(product *product.Product) bool {
	_, ok := o.items[product.ID()]
	return ok
}

//applyCoupon is a function for applying a coupon to an order
//Returns true if coupon application is successful or false and an error describing the failure
func (o *Order) applyCoupon(coupon *coupon.Coupon) (bool, *errors.Error) {
	if StatusDraft != o.status {
		return false, errors.Wrap(fmt.Errorf("Can't apply coupon in order %v, status is %v (not draft)", o.id, statusMap[o.status]), 0)
	}
	if 0 == len(o.items) {
		return false, errors.Wrap(fmt.Errorf("Can't apply coupon: order %v has no item", o.id), 0)
	}
	if _, err := coupon.CanBeApplied(); err != nil {
		return false, errors.Wrap(fmt.Errorf("Can't apply coupon with id %v in order %v: %v", coupon.ID(), o.id, err), 0)
	}
	if _, err := o.calculateAmount(coupon); err != nil {
		return false, errors.Wrap(fmt.Errorf("Can't apply coupon with id %v in order %v: %v", coupon.ID(), o.id, err), 0)
	}
	o.SetCoupon(coupon)
	return true, nil
}

//calculateAmount is a function for calculating the order's amount (subtracted with discount from a given coupon)
//(regardless of the order's status, order item's product status, and the coupon status)
func (o *Order) calculateAmount(coupon *coupon.Coupon) (bool, *errors.Error) {
	amount := decimal.New(0, 0)
	for _, val := range o.items {
		amount = amount.Add(val.Product().Price().Mul(decimal.New(int64(val.Quantity()), 0)))
	}
	if coupon != nil {
		discountAmount := coupon.GetDiscountAmount(amount)
		amount = amount.Sub(discountAmount)
		if decimal.New(0, 0).GreaterThanOrEqual(amount) {
			return false, errors.Wrap(fmt.Errorf("Zero or less calculated amount of order with id %v (applied with coupon with id %v)", o.id, o.coupon.ID()), 0)
		}
	}
	o.amount = amount
	return true, nil
}

//Submit is a function for submitting order
func (o *Order) Submit(shippingName, shippingAddress string, coupon *coupon.Coupon) (bool, *errors.Error) {
	if StatusDraft != o.status {
		return false, errors.Wrap(fmt.Errorf("Can't submit: order %v status is %v (not draft)", o.id, o.status), 0)
	}
	if 0 == len(o.items) {
		return false, errors.Wrap(fmt.Errorf("Can't submit: order %v has no item", o.id), 0)
	}
	for _, val := range o.items {
		if canOrder, err := val.Product().CanBeOrdered(val.Quantity()); false == canOrder {
			return false, errors.Wrap(fmt.Errorf("Can't submit: order %v for item with product id %v: %v", o.id, val.Product().ID(), err), 0)
		}
	}
	//try applying coupon if exist
	if coupon != nil {
		_, err := o.applyCoupon(coupon)
		if err != nil {
			return false, errors.Wrap(fmt.Errorf("Can't submit: order %v can't apply coupon %v: %v", o.id, coupon.ID(), err), 0)
		}
		coupon.SetStock(coupon.Stock() - 1)
	} else {
		o.calculateAmount(nil)
	}
	o.status = StatusSubmitted
	o.submittedDate = time.Now()
	o.shippingStatus = ShipStatusNone

	for _, val := range o.items {
		o.mu.Lock()
		val.Product().SetStock(val.Product().Stock() - int64(val.Quantity()))
		o.mu.Unlock()
	}

	return true, nil
}

//Process is a function for processing order
func (o *Order) Process() (bool, *errors.Error) {
	if StatusSubmitted != o.status {
		return false, errors.Wrap(fmt.Errorf("Can't process: order %v status is %v (not submitted)", o.id, o.status), 0)
	}
	o.status = StatusProcessed
	o.processedDate = time.Now()
	return true, nil
}

//Cancel is a function for canceling order
func (o *Order) Cancel() (bool, *errors.Error) {
	if StatusSubmitted != o.status {
		return false, errors.Wrap(fmt.Errorf("Can't cancel: order %v status is %v (not submitted)", o.id, o.status), 0)
	}
	o.status = StatusCanceled
	return true, nil
}

//ProcessShipping is a function for processing order shipping
func (o *Order) ProcessShipping(trackingNo string) (bool, *errors.Error) {
	if StatusProcessed != o.status {
		return false, errors.Wrap(fmt.Errorf("Can't process: order %v status is %v (not processed)", o.id, o.status), 0)
	}
	o.shippingTrackingID = trackingNo
	o.shippingStatus = ShipStatusOnProcess
	return true, nil
}

//FinishOrder is a function for finishing order
func (o *Order) FinishOrder() (bool, *errors.Error) {
	if StatusProcessed != o.status {
		return false, errors.Wrap(fmt.Errorf("Can't finish: order %v status is %v (not processed)", o.id, o.status), 0)
	}
	o.status = StatusDelivered
	o.shippingStatus = ShipStatusDelivered
	return true, nil
}
