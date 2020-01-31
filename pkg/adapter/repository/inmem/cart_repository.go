package inmem

import (
	"errors"
	"fmt"
	"time"

	"github.com/lucsky/cuid"

	"github.com/yauritux/cartsvc/pkg/adapter/repository/inmem/model"
	uc "github.com/yauritux/cartsvc/pkg/usecase"
)

var carts []*model.Cart

type CartRepository struct {
	currentCart *model.Cart
}

func NewCartRepository(uid string) *CartRepository {
	if carts == nil {
		carts = []*model.Cart{
			newCart(uid),
		}
		return &CartRepository{carts[0]}
	}

	for i, v := range carts {
		if v.UserID == uid && v.Status == "open" {
			return &CartRepository{carts[i]}
		}
	}
	carts = append(carts, newCart(uid))
	return &CartRepository{carts[len(carts)-1]}
}

func (r *CartRepository) FetchUserCart(userID string) (interface{}, error) {
	for i, v := range carts {
		if v.UserID == userID {
			return buildCartUsecaseModel(carts[i]), nil
		}
	}

	return nil, nil
}

func (r *CartRepository) AddToCart(cartID string, item interface{}) error {
	currUserCart, err := r.getCurrentUserCart(cartID)
	if err != nil {
		return err
	}

	if currUserCart.Status != "open" {
		return fmt.Errorf("cart with ID of %s is %s, please create a new cart",
			cartID, currUserCart.Status)
	}

	cartItem, ok := item.(*uc.CartItem)
	if !ok {
		return errors.New("failed to add item, invalid type of cart item")
	}

	currUserCart.Items = append(currUserCart.Items, r.BuildCartItemRepositoryModel(cartItem).(*model.CartItem))
	r.currentCart.Items = currUserCart.Items
	return nil
}

func (r *CartRepository) RemoveItem(id string, itemID string) error {
	currUserCart, err := r.getCurrentUserCart(id)
	if err != nil {
		return err
	}

	if currUserCart.Status != "open" {
		return fmt.Errorf("cart with ID of %s is %s, please create a new cart",
			id, currUserCart.Status)
	}

	var newCartItems []*model.CartItem

	for i, v := range currUserCart.Items {
		if v.ID != itemID {
			continue
		}
		newCartItems = make([]*model.CartItem, 0)
		newCartItems = append(newCartItems[:i], newCartItems[i+1:]...)
		currUserCart.Items = newCartItems
	}
	return nil
}

func (r *CartRepository) UpdateItem(id string, item interface{}) error {
	currUserCart, err := r.getCurrentUserCart(id)
	if err != nil {
		return err
	}

	if currUserCart.Status != "open" {
		return fmt.Errorf("cart with ID of %s is %s, please create a new cart",
			id, currUserCart.Status)
	}

	cartItem, ok := item.(*uc.CartItem)
	if !ok {
		return errors.New("failed to update item, invalid type of cart item")
	}

	updatedCartItem := r.BuildCartItemRepositoryModel(cartItem).(*model.CartItem)
	for i, v := range currUserCart.Items {
		if v.ID == updatedCartItem.ID {
			currUserCart.Items[i] = updatedCartItem
			break
		}
	}

	return nil
}

func (r *CartRepository) Checkout(id string) interface{} {
	currUserCart, err := r.getCurrentUserCart(id)
	if err != nil {
		return err
	}

	if currUserCart.Status != "open" {
		return fmt.Errorf("failed to checkout cart with ID of %s, it is already in %s",
			id, currUserCart.Status)
	}

	return currUserCart
}

func (r *CartRepository) Canceled(id string) error {
	currUserCart, err := r.getCurrentUserCart(id)
	if err != nil {
		return err
	}

	if currUserCart.Status != "open" {
		return fmt.Errorf("cannot cancel the cart with status of %s", currUserCart.Status)
	}

	currUserCart.Status = "canceled"
	return nil
}

func (r *CartRepository) Close(id string) error {
	currUserCart, err := r.getCurrentUserCart(id)
	if err != nil {
		return err
	}

	if currUserCart.Status != "payment_processing" {
		return fmt.Errorf("failed to close cart with ID of %s. It is %s, please settle the payment first",
			id, currUserCart.Status)
	}

	currUserCart.Status = "closed"

	return nil
}

func (r *CartRepository) BuildCartItemRepositoryModel(item *uc.CartItem) interface{} {
	return &model.CartItem{
		ID:    item.ID,
		Name:  item.Name,
		Qty:   item.Qty,
		Price: item.Price,
		Disc:  item.Disc,
	}
}

func newCart(uid string) *model.Cart {
	return &model.Cart{
		ID:        cuid.New(),
		UserID:    uid,
		Status:    "open",
		Items:     make([]*model.CartItem, 0),
		CreatedAt: time.Now(),
	}
}

func (r *CartRepository) getCurrentUserCart(cartID string) (*model.Cart, error) {
	for i, v := range carts {
		if v.ID == cartID {
			return carts[i], nil
		}
	}
	return nil, fmt.Errorf("cannot find cart with id %s", cartID)
}

func buildCartUsecaseModel(cart *model.Cart) *uc.Cart {
	ucCart := &uc.Cart{
		ID:        cart.ID,
		UserID:    cart.UserID,
		Status:    cart.Status,
		CreatedAt: cart.CreatedAt,
	}
	ucCartItems := make([]*uc.CartItem, 0)
	for _, v := range cart.Items {
		ucCartItems = append(ucCartItems, &uc.CartItem{
			ID:    v.ID,
			Name:  v.Name,
			Qty:   v.Qty,
			Price: v.Price,
			Disc:  v.Disc,
		})
	}
	ucCart.Items = ucCartItems
	return ucCart
}
