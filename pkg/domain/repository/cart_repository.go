package repository

type CartRepository interface {
	FetchUserCart(userID string) (interface{}, error)
	AddToCart(cartID string, item interface{}) error
	RemoveItem(cartID string, itemID string) error
	UpdateItem(cartID string, item interface{}) error
	Checkout(cartID string) interface{}
	Canceled(cartID string) error
	Close(cartID string) error
}
