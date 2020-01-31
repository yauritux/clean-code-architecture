package port

import (
	uc "github.com/yauritux/cartsvc/pkg/usecase"
)

type CartInputPort interface {
	FetchUserCart(userID string) (interface{}, error)
	AddToCart(userID string, item interface{}) error
}

type CartOutputPort interface {
	BuildCartItemRepositoryModel(*uc.CartItem) interface{}
}
