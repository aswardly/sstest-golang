//Package coupon provides the business domain models definitions of coupon
package coupon

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-errors/errors"
	"github.com/shopspring/decimal"
)

//StatusActive is const for 'active' coupon status
const StatusActive string = "A"

//StatusInactive is const for 'inactive' coupon status
const StatusInactive string = "I"

//StatusExpired is const for 'expired' coupon status
const StatusExpired string = "E"

//StatusSuspended is const for 'suspended' coupon status
const StatusSuspended string = "S"

//statusMap is a map of known status code and its label pairs
var statusMap = map[string]string{
	StatusActive:    "Active",
	StatusInactive:  "Inactive",
	StatusExpired:   "Expired",
	StatusSuspended: "Suspended",
}

//KindValue is const for 'value' coupon kind/type
const KindValue string = "V"

//KindPercentage is const for 'percentage' coupon kind/type
const KindPercentage string = "P"

//kindMap is a map of known kind/type code and its label pairs
var kindMap = map[string]string{
	KindValue:      "Value",
	KindPercentage: "Percentage",
}

//Coupon is business domain model definition of product
type Coupon struct {
	id        string
	status    string
	stock     int64
	kind      string
	value     decimal.Decimal
	startDate time.Time
	endDate   time.Time
	mu        sync.Mutex
}

//New creates a new coupon model struct, initializes it's properties and returns a reference to it
func New(id string) *Coupon {
	//coupon's start date defaults to today's date at beginning of day (at 00:00:00)
	//coupon's expiry date defaults to the date in 90 days from coupon's start date
	currentTime := time.Now()
	couponStartDate := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
	couponEndDate := couponStartDate.AddDate(0, 0, 90)

	return &Coupon{
		id,
		StatusInactive,
		0,
		KindPercentage,     //default kind/type is percentage
		decimal.New(10, 0), //default value is 10 (10 percent)
		couponStartDate,
		couponEndDate,
		*new(sync.Mutex),
	}
}

//ID is a getter function for returning a coupon's id
func (c *Coupon) ID() string {
	return c.id
}

//Status is a getter function for returning a coupon's status
func (c *Coupon) Status() string {
	return c.status
}

//Stock is a getter function for returning a coupon's stock
func (c *Coupon) Stock() int64 {
	return c.stock
}

//Kind is a getter function for returning a coupon's kind/type
func (c *Coupon) Kind() string {
	return c.kind
}

//Value is a getter function for returning a coupon's value
func (c *Coupon) Value() decimal.Decimal {
	return c.value
}

//StartDate is a getter function for returning a coupon's start date
func (c *Coupon) StartDate() time.Time {
	return c.startDate
}

//EndDate is a getter function for returning a coupon's end date
func (c *Coupon) EndDate() time.Time {
	return c.endDate
}

//SetID is a setter function for setting a coupon's id
func (c *Coupon) SetID(id string) *Coupon {
	c.id = id
	return c
}

//SetStock is a setter function for setting a coupon's stock
func (c *Coupon) SetStock(stock int64) *Coupon {
	c.stock = stock
	return c
}

//SetStatus is a setter function for setting a coupon's status
func (c *Coupon) SetStatus(status string) (*Coupon, *errors.Error) {
	if _, ok := statusMap[status]; false == ok {
		//note: defensive code, return nil when error is encountered (possible runtime error on caller code when method chaining)
		return nil, errors.Wrap(fmt.Errorf("Can't set unknown status type: %v", status), 0)
	}
	c.status = status
	return c, nil
}

//SetKind is a setter function for setting a coupon's kind/type
func (c *Coupon) SetKind(kind string) (*Coupon, *errors.Error) {
	if _, ok := kindMap[kind]; false == ok {
		//note: defensive code, return nil when error is encountered (possible runtime error on caller code when method chaining)
		return nil, errors.Wrap(fmt.Errorf("Can't set unknown kind type: %v", kind), 0)
	}
	//if kind is changed to percentage and current value is more than 100, return false and an error
	if KindPercentage == kind && c.value.GreaterThanOrEqual(decimal.New(100, 0)) {
		return nil, errors.Wrap(fmt.Errorf("Can't set kind to %v when value is %v", kindMap[c.kind], c.value.String()), 0)
	}
	c.kind = kind

	return c, nil
}

//SetValue is a setter function for setting a coupon's value
func (c *Coupon) SetValue(value decimal.Decimal) (*Coupon, *errors.Error) {
	if value.LessThanOrEqual(decimal.New(0, 0)) {
		return nil, errors.Wrap(fmt.Errorf("Can't set value to 0"), 0)
	}
	//if kind is percentage and value to set is more than 100, return false and an error
	if KindPercentage == c.kind && value.GreaterThanOrEqual(decimal.New(100, 0)) {
		return nil, errors.Wrap(fmt.Errorf("Can't set value to %v when kind is %v", c.value.String(), kindMap[c.kind]), 0)
	}
	c.value = value
	return c, nil
}

//SetStartDate is a setter function for setting a coupon's start date
func (c *Coupon) SetStartDate(startDate time.Time) (*Coupon, *errors.Error) {
	if startIsAfterEnd := c.endDate.Before(startDate); startIsAfterEnd {
		//note: defensive code, return nil when error is encountered (possible runtime error on caller code when method chaining)
		return nil, errors.Wrap(fmt.Errorf("Can't set start date to be later than %v", c.endDate.Format("2000-12-31 23:59:59")), 0)
	}
	c.startDate = startDate
	return c, nil
}

//SetEndDate is a setter function for setting a coupon's end date
func (c *Coupon) SetEndDate(endDate time.Time) (*Coupon, *errors.Error) {
	if EndIsBeforeStart := c.startDate.After(endDate); EndIsBeforeStart {
		//note: defensive code, return nil when error is encountered (possible runtime error on caller code when method chaining)
		return nil, errors.Wrap(fmt.Errorf("Can't set end date to be earlier than %v", c.startDate.Format("2000-12-31 23:59:59")), 0)
	}
	c.endDate = endDate
	return c, nil
}

//Business logic methods

//IsEarly is a function for inquiring whether the coupon start date is in the future (indicating coupon is too early to be applied)
func (c *Coupon) IsEarly() bool {
	if isEarly := c.startDate.After(time.Now()); isEarly {
		//automatically set status to inactive
		c.SetStatus(StatusInactive)
		return true
	}
	return false
}

//IsExpired is a function for inquiring whether the coupon end date is in the past (indicating coupon is expired)
func (c *Coupon) IsExpired() bool {
	if isExpired := c.endDate.Before(time.Now()); isExpired {
		//automatically set status to expired
		c.SetStatus(StatusExpired)
		return true
	}
	return false
}

//CanBeApplied is a function for inquiring whether coupon can be used or not
//returns boolean and string (reason explaining why coupon can't be used)
func (c *Coupon) CanBeApplied() (bool, *errors.Error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.status != StatusActive {
		return false, errors.Wrap(fmt.Errorf("Coupon status is %v (not active)", statusMap[c.status]), 0)
	}
	if c.IsEarly() {
		return false, errors.Wrap(fmt.Errorf("coupon can't be applied before %v", c.startDate.Format("2000-12-31 23:59:59")), 0)
	}
	if c.IsExpired() {
		return false, errors.Wrap(fmt.Errorf("coupon can't be applied after %v", c.endDate.Format("2000-12-31 23:59:59")), 0)
	}
	if c.stock <= 0 {
		return false, errors.Wrap(fmt.Errorf("coupon can't be applied, stock is %d (no stock)", c.stock), 0)
	}

	return true, nil
}

//GetDiscountAmount returns discount amount from a certain given amount when applied by this coupon
func (c *Coupon) GetDiscountAmount(amount decimal.Decimal) decimal.Decimal {
	var retAmount decimal.Decimal
	if KindPercentage == c.kind {
		//coupon kind/type is KindPercentage
		retAmount = amount.Mul(c.value).Div(decimal.New(100, 0))
		return retAmount
	}
	//coupon kind/type is KindValue
	retAmount = c.value
	return retAmount
}
