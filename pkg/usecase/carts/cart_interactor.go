package carts

import (
	"errors"
	"time"

	"github.com/yauritux/cartsvc/pkg/domain/aggregate"
	"github.com/yauritux/cartsvc/pkg/domain/entity"
	"github.com/yauritux/cartsvc/pkg/domain/repository"
	vo "github.com/yauritux/cartsvc/pkg/domain/valueobject"
	. "github.com/yauritux/cartsvc/pkg/sharedkernel/enum"
	e "github.com/yauritux/cartsvc/pkg/sharedkernel/error"
	prodUsecase "github.com/yauritux/cartsvc/pkg/usecase/products"
)

type CartUsecase struct {
	cartRepo repository.CartRepository
	prodRepo repository.ProductRepository
}

type Cart struct {
	ID         string
	UserID     string
	Status     CartStatus
	Items      []*CartItem
	CreatedAt  time.Time
	CanceledAt *time.Time
}

type CartItem struct {
	ID    string
	Name  string
	Qty   int
	Price float64
	Disc  float64
}

func NewCartUsecase(r1 repository.CartRepository, r2 repository.ProductRepository) *CartUsecase {
	return &CartUsecase{r1, r2}
}

func (this *CartUsecase) FetchUserCart(userID string) (interface{}, error) {
	if userID == "" {
		return nil, e.NewErrNoData("cannot fetch user cart, 'user_id' is missing")
	}

	cart, err := this.cartRepo.FetchUserCart(userID)
	if err != nil {
		return nil, err
	}

	ucCart, ok := cart.(*Cart)
	if !ok {
		return nil, e.NewErrConversion("cannot fetch user cart, invalid type of cart usecase model")
	}

	return ucCart, nil
}

func (this *CartUsecase) AddToCart(userID string, item interface{}) error {
	userCart, err := this.cartRepo.FetchUserCart(userID)
	if err != nil {
		return err
	}
	currentCart, ok := userCart.(*Cart)
	if !ok {
		return errors.New("conversion failed, invalid type of cart usecase model")
	}

	prodItem, ok := item.(*CartItem)
	if !ok {
		return errors.New("conversion failed, invalid type of product item usecase model")
	}

	product, err := this.prodRepo.FindByProductID(prodItem.ID)
	if err != nil {
		return err
	}

	ucProduct, ok := product.(*prodUsecase.Product)
	if !ok {
		return errors.New("conversion failed, invalid type of product usecase model")
	}

	cart := aggregate.NewUserCart(
		&entity.User{
			UserID: currentCart.UserID,
		}, &entity.Cart{
			ID:        currentCart.ID,
			UserID:    currentCart.UserID,
			Status:    currentCart.Status,
			Items:     buildCartVOItems(currentCart.Items),
			CreatedAt: currentCart.CreatedAt,
		})

	addedItem, err := cart.AddItemToCart(&entity.Product{
		ID:    ucProduct.ID,
		Name:  ucProduct.Name,
		Stock: ucProduct.Stock,
		Price: ucProduct.Price,
		Disc:  ucProduct.Disc,
	}, prodItem.Qty)
	if err != nil {
		switch err.(type) {
		case *e.ErrDuplicateData:
			return this.cartRepo.UpdateItem(cart.FetchCartInfo().ID, buildCartUsecaseItem(addedItem))
		default:
			return err
		}
	}

	return this.cartRepo.AddToCart(cart.FetchCartInfo().ID, buildCartUsecaseItem(addedItem))
}

func buildCartUsecaseItem(item interface{}) *CartItem {
	var ucCartItem *CartItem
	switch item.(type) {
	case *vo.CartItem:
		cartItem := item.(*vo.CartItem)
		ucCartItem = &CartItem{
			ID:    cartItem.ProdID,
			Name:  cartItem.ProdName,
			Qty:   cartItem.Qty,
			Price: cartItem.Price,
			Disc:  cartItem.Disc,
		}
	}
	return ucCartItem
}

func buildCartVOItems(items interface{}) []*vo.CartItem {
	voCartItems := make([]*vo.CartItem, 0)
	switch items.(type) {
	case []*CartItem:
		cartItems := items.([]*CartItem)
		for _, v := range cartItems {
			voCartItems = append(voCartItems, &vo.CartItem{
				ProdID:   v.ID,
				ProdName: v.Name,
				Qty:      v.Qty,
				Price:    v.Price,
				Disc:     v.Disc,
			})
		}
	}
	return voCartItems
}
