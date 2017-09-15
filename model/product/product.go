//Package product provides the business domain models definitions of product
package product

import (
	"fmt"
	"sync"

	"github.com/go-errors/errors"
	"github.com/shopspring/decimal"
)

//StatusPrototype is const for 'prototype' product status
const StatusPrototype string = "P"

//StatusAvailable is const for 'available' product status
const StatusAvailable string = "A"

//StatusDiscontinued is const for 'discontinued' product status
const StatusDiscontinued string = "D"

//statusSlice is a map of known status code and its label pairs
var statusSlice = map[string]string{
	StatusPrototype:    "Prototype",
	StatusAvailable:    "Available",
	StatusDiscontinued: "Discontinued",
}

//Product is business domain model definition of product
type Product struct {
	id     string
	name   string
	status string
	price  decimal.Decimal
	stock  int64
	mu     sync.Mutex
}

//New creates a new product model struct, initializes it's properties and returns a reference to it
func New(id, name string) *Product {
	return &Product{
		id,
		name,
		StatusPrototype,
		decimal.New(0, 0),
		0,
		*new(sync.Mutex),
	}
}

//ID is a getter function for returning a product's id
func (p *Product) ID() string {
	return p.id
}

//Name is a getter function for returning a product's name
func (p *Product) Name() string {
	return p.name
}

//Status is a getter function for returning a product's status
func (p *Product) Status() string {
	return p.status
}

//Price is a getter function for returning a product's price
func (p *Product) Price() decimal.Decimal {
	return p.price
}

//Stock is a getter function for returning a product's stock
func (p *Product) Stock() int64 {
	return p.stock
}

//SetID is a setter function for setting a product's id
func (p *Product) SetID(id string) *Product {
	p.id = id
	return p
}

//SetName is a setter function for setting a product's name
func (p *Product) SetName(name string) *Product {
	p.name = name
	return p
}

//SetStock is a setter function for setting a product's stock
func (p *Product) SetStock(stock int64) *Product {
	p.stock = stock
	return p
}

//SetPrice is a setter function for setting a product's price
func (p *Product) SetPrice(price decimal.Decimal) (*Product, *errors.Error) {
	if zero := decimal.New(0, 0); zero.GreaterThan(price) {
		return nil, errors.Wrap(fmt.Errorf("Can't set negative value %v for price", price.String()), 0)
	}
	p.price = price
	return p, nil
}

//SetStatus is a setter function for setting a product's status
func (p *Product) SetStatus(status string) (*Product, *errors.Error) {
	if _, ok := statusSlice[status]; false == ok {
		//note: defensive code, return nil when error is encountered (possible runtime error on caller code when method chaining)
		return nil, errors.Wrap(fmt.Errorf("Can't set unknown status type: %v", status), 0)
	}
	p.status = status
	return p, nil
}

//Business logic methods

//CanBeOrdered is a function for inquiring whether product can be ordered or not
//returns boolean and string (reason explaining why product can't be ordered)
func (p *Product) CanBeOrdered(quantity int) (bool, *errors.Error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if quantity <= 0 {
		return false, errors.Wrap(fmt.Errorf("Can't order product (id: %v) with quantity %d", p.id, quantity), 0)
	}
	if p.status != StatusAvailable {
		return false, errors.Wrap(fmt.Errorf("Product (id: %v) status is not available", p.id), 0)
	}
	if p.stock-int64(quantity) < 0 {
		return false, errors.Wrap(fmt.Errorf("Product (id: %v) stock %d does not have enough quantity %d", p.id, p.stock, quantity), 0)
	}
	return true, nil
}
