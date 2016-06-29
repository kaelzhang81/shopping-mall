package interfaces

import (
	"net/http"
)

import (
	"fmt"
	"usecases"
)

type OrderInteractor interface {
	Items(userId, orderId int) ([]usercases.Item, error)
	Add(userId, orderId, itemId int) error
}

type webServiceHandler struct {
	OrderInteractor OrderInteractor
}

func (handler webServiceHandler) ShowOrder(res http.ResponseWriter, req *http.Request) {

}
