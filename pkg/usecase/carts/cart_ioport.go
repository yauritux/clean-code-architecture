package carts

type CartInputPort interface {
	FetchUserCart(userID string) (interface{}, error)
	AddToCart(userID string, item interface{}) error
}

type CartOutputPort interface {
	BuildCartItemRepositoryModel(*CartItem) interface{}
}
