package domain

import (
	"errors"
)

type CustomerRepository interface {
	Store(customer Customer)
	FindById(id int) Customer
}

type ItemRepository interface {
	Store(item Item)
	FindById(id int) Item
}

type OrderRepository interface {
	Store(order Order)
	FindById(id int) Order
}

type Customer struct {
	Id   int
	Name string
}

type Item struct {
	Id        int
	Name      string
	Value     float64
	Available bool
}

type Order struct {
	Id       int
	Customer Customer
	Items    []Item
}

func (order *Order) Add(item Item) error {
	if !item.Available {
		return errors.New("Can not add unavalible item to order")
	}
	if order.value()+item.Value > 250 {
		return errors.New("an order may not exceed a total value of $250")
	}
	order.Items = append(order.Items, item)
	return nil
}

func (order *Order) value() float64 {
	sum := 0.0
	for i := range order.Items {
		sum += order.Items[i].Value
	}
	return sum
}
