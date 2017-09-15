//Package order provides the business domain models definitions of order and order item
package order

import (
	"fmt"
	"sstest/model/product"
	"sync"

	"github.com/go-errors/errors"
)

//Item is business domain model definition of order item
type Item struct {
	id       string
	order    *Order
	product  *product.Product
	quantity int
	mu       sync.Mutex
}

//NewItem creates a new order item model struct, initializes it's properties and returns a reference to it
func NewItem(id string, orderID *Order, productID *product.Product) *Item {
	return &Item{id, orderID, productID, 0, *new(sync.Mutex)}
}

//ID is a getter function for returning an order item's id
func (i *Item) ID() string {
	return i.id
}

//Order is a getter function for returning an order item's order
func (i *Item) Order() *Order {
	return i.order
}

//Product is a getter function for returning an order item's product id
func (i *Item) Product() *product.Product {
	return i.product
}

//Quantity is a getter function for returning an order item's quantity
func (i *Item) Quantity() int {
	return i.quantity
}

//SetID is a setter function for setting an order item's id
func (i *Item) SetID(id string) *Item {
	i.id = id
	return i
}

//SetOrder is a setter function for setting an order item's order
func (i *Item) SetOrder(order *Order) *Item {
	i.order = order
	return i
}

//SetProduct is a setter function for setting an order item's product
func (i *Item) SetProduct(product *product.Product) *Item {
	i.product = product
	return i
}

//SetQuantity is a setter function for setting an order item's quantity
func (i *Item) SetQuantity(quantity int) *Item {
	i.quantity = quantity
	return i
}

//Business logic methods

//AddQuantity is a function for adding some quantity to an order item
/*
func (i *Item) AddQuantity(quantity int) (*Item, error) {
	//check whether addition results in quantity less than 0 or overflow
	if quantity > 0 {
		//positive quantity to add to existing item quantity
		if i.quantity > math.MaxInt8-quantity {
			return nil, errors.Wrap(fmt.Errorf("Addition with %v results in integer overflow", quantity), 0)
		}
	} else {
		//negative quantity to add to existing item quantity
		if i.quantity < math.MinInt8-quantity {
			return nil, errors.Wrap(fmt.Errorf("Addition with %v results in integer overflow", quantity), 0)
		}
		if 0 >= i.quantity+quantity {
			return nil, errors.Wrap(fmt.Errorf("Addition with %v makes quantity become 0 or less", quantity), 0)
		}
	}
	//alternatively to check for overflow
	//if ((i.quantity + quantity) < i.quantity) != (quantity < 0)

	i.quantity += quantity
	return i, nil
}
*/

//AddQuantity is a function for adding some quantity to an order item
func (i *Item) AddQuantity(quantity int) (*Item, *errors.Error) {
	//check whether addition results in quantity less than 0 or overflow
	if ((i.quantity + quantity) < i.quantity) != (quantity < 0) {
		return nil, errors.Wrap(fmt.Errorf("Addition with %v results in integer overflow", quantity), 0)
	}
	if 0 >= i.quantity+quantity {
		return nil, errors.Wrap(fmt.Errorf("Addition with %v makes quantity become 0 or less", quantity), 0)
	}

	i.quantity += quantity
	return i, nil
}

//SubtractQuantity is a function for subtracting some quantity from an order item
func (i *Item) SubtractQuantity(quantity int) (*Item, *errors.Error) {
	_, err := i.AddQuantity(0 - quantity)
	if err != nil {
		return nil, err
	}
	return i, nil
}
