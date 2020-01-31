package aggregate

import (
	"errors"
	"fmt"

	"github.com/yauritux/cartsvc/pkg/domain/entity"
	vo "github.com/yauritux/cartsvc/pkg/domain/valueobject"
	"github.com/yauritux/cartsvc/pkg/sharedkernel/enum"
	e "github.com/yauritux/cartsvc/pkg/sharedkernel/error"
)

//UserCart is represents aggregate which contains our core business logic
//(enterprise business rules) for user's shopping cart
type UserCart struct {
	user *entity.User
	cart *entity.Cart
}

func NewUserCart(user *entity.User, cart *entity.Cart) *UserCart {
	if cart.Status == "" {
		cart.Status = enum.Open
	}
	if cart.Items == nil {
		cart.Items = make([]*vo.CartItem, 0)
	}
	return &UserCart{
		user: user,
		cart: cart,
	}
}

func (userCart *UserCart) AddItemToCart(prod *entity.Product, qty int) (*vo.CartItem, error) {
	if qty > prod.Stock {
		return nil, errors.New("out of stock")
	}
	addedItem := &vo.CartItem{
		ProdID:   prod.ID,
		ProdName: prod.Name,
		Qty:      qty,
		Price:    prod.Price,
		Disc:     prod.Disc,
	}
	if err := userCart.Validate(); err != nil {
		return nil, err
	}
	if userCart.cart.Status != enum.Open {
		return nil, fmt.Errorf("cannot add item to a cart with status as %s", userCart.cart.Status)
	}

	qtyOverride := false
	for i, v := range userCart.cart.Items {
		if prod.ID == v.ProdID {
			userCart.cart.Items[i].Qty = userCart.cart.Items[i].Qty + qty
			addedItem.Qty = userCart.cart.Items[i].Qty
			qtyOverride = true
			break
		}
	}

	if qtyOverride {
		return addedItem, e.NewErrDuplicateData("item exists, updated amount of cart existing item")
	}
	return addedItem, nil
}

func (userCart *UserCart) UpdateItemInCart(item *vo.CartItem) error {
	if userCart.cart.Items == nil || len(userCart.cart.Items) == 0 {
		return e.NewErrNoData("cart is still empty")
	}
	itemUpdated := false
	for i, v := range userCart.cart.Items {
		if item.ProdID == v.ProdID {
			userCart.cart.Items[i] = item
			itemUpdated = true
			break
		}
	}
	if itemUpdated == false {
		return e.NewErrNoData(fmt.Sprintf(
			"cannot find cart item %s - %s within the cart",
			item.ProdID, item.ProdName,
		))
	}

	return nil
}

func (userCart *UserCart) RemoveItemFromCart(itemID string) error {
	if userCart.cart.Items == nil || len(userCart.cart.Items) == 0 {
		return errors.New("cart is still empty")
	}

	updatedCartItems := make([]*vo.CartItem, 0)

	for i, v := range userCart.cart.Items {
		if v.ProdID == itemID {
			updatedCartItems = append(userCart.cart.Items[:i], userCart.cart.Items[i+1:]...)
			userCart.cart.Items = updatedCartItems
			break
		}
	}
	if len(updatedCartItems) == 0 {
		return fmt.Errorf("cannot find cart item with ID %s", itemID)
	}
	return nil
}

func (userCart *UserCart) FetchUserInfo() *entity.User {
	return userCart.user
}

func (userCart *UserCart) FetchCartInfo() *entity.Cart {
	return userCart.cart
}

func (userCart *UserCart) Validate() error {
	if userCart.cart.ID == "" {
		return errors.New("cart 'session_id' is missing")
	}
	if userCart.cart.UserID == "" {
		return errors.New("cart is orphaned (who's owning this cart ?)")
	}
	if userCart.user.UserID == "" {
		return errors.New("cart 'user_id' is missing")
	}
	return nil
}
