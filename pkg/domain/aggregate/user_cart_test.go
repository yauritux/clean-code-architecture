package aggregate

import (
	"errors"
	"testing"
	"time"

	"github.com/yauritux/cartsvc/pkg/domain/entity"
	vo "github.com/yauritux/cartsvc/pkg/domain/valueobject"
	"github.com/yauritux/cartsvc/pkg/sharedkernel/enum"
	e "github.com/yauritux/cartsvc/pkg/sharedkernel/error"
)

var u *entity.User
var c *entity.Cart
var p *entity.Product

func init() {
	u = &entity.User{
		UserID:   "yauritux",
		Username: "Yauri Attamimi",
		Phone:    "+6282225251437",
		Email:    "yauritux@gmail.com",
		BillingAddress: vo.BuyerAddress{
			StreetName: "Kalibata Raya",
			City:       "Jakarta",
			Postal:     "12750",
			Province:   "DKI Jakarta",
			Region:     "South Jakarta",
			Country:    "Indonesia",
			Type:       enum.BillingAddress,
		},
		ShippingAddress: vo.BuyerAddress{
			StreetName: "Kalibata Raya",
			City:       "Jakarta",
			Postal:     "12750",
			Province:   "DKI Jakarta",
			Region:     "South Jakarta",
			Country:    "Indonesia",
			Type:       enum.ShippingAddress,
		},
	}
	c = &entity.Cart{
		ID:        "123",
		UserID:    "yauritux",
		Items:     make([]*vo.CartItem, 0),
		CreatedAt: time.Now(),
	}
	p = &entity.Product{
		ID:    "001",
		Name:  "Shuriken",
		Stock: 100,
		Price: 125.5,
		Disc:  0.0,
	}
}

func TestNewUserCart(t *testing.T) {
	userCart := NewUserCart(u, c)
	if userCart.cart.Status != enum.Open {
		t.Errorf("initial cart status should be %s, got %s", enum.Open, userCart.cart.Status)
	}
	if userCart.user != u {
		t.Errorf("expected user %v, got %v", u, userCart.user)
	}
	if userCart.cart != c {
		t.Errorf("expected cart %v, got %v", c, userCart.cart)
	}
}

func TestAddItemToCartOutOfStock(t *testing.T) {
	userCart := NewUserCart(u, c)
	_, err := userCart.AddItemToCart(p, 105)
	if err == nil {
		t.Errorf("expected got error of %v, but no error is thrown", errors.New("out of stock"))
	}
}

func TestAddItemToCartMissingCartSessionID(t *testing.T) {
	c.ID = ""
	userCart := NewUserCart(u, c)
	_, err := userCart.AddItemToCart(p, 10)
	if err == nil {
		t.Error("expected got error, but no error is thrown")
	}
}

func TestAddITemToCartMissingCartUserID(t *testing.T) {
	c.UserID = ""
	userCart := NewUserCart(u, c)
	_, err := userCart.AddItemToCart(p, 5)
	if err == nil {
		t.Error("expected got error, but no error is thrown")
	}
}

func TestAddItemToCartMissingUserID(t *testing.T) {
	u.UserID = ""
	userCart := NewUserCart(u, c)
	_, err := userCart.AddItemToCart(p, 5)
	if err == nil {
		t.Error("expected got error, but no error is thrown")
	}
}

func TestAddItemToCartUnopenedCart(t *testing.T) {
	c.Status = enum.Closed
	userCart := NewUserCart(u, c)
	_, err := userCart.AddItemToCart(p, 5)
	if err == nil {
		t.Error("expected got error, but no error is thrown")
	}
}

func TestAddItemToCart(t *testing.T) {
	userCart := NewUserCart(u, c)
	addedItem, err := userCart.AddItemToCart(p, 5)
	if err != nil {
		t.Errorf("Expected no error, but got error of %v", err)
	}
	if addedItem.ProdID != p.ID {
		t.Errorf("Expected product id of %s, got %s", p.ID, addedItem.ProdID)
	}
	if addedItem.Qty != 5 {
		t.Errorf("Expected added qty = %d, got %d", 5, addedItem.Qty)
	}
}

func TestAddItemToCartExistingItem(t *testing.T) {
	userCart := NewUserCart(u, c)
	userCart.AddItemToCart(p, 5)
	addedItem, err := userCart.AddItemToCart(p, 2)
	if err == nil {
		t.Errorf("Expected got error of %v, but no error is thrown", e.ErrDuplicateData{})
	}
	if addedItem.Qty != 7 {
		t.Errorf("Expected qty of %d, got %d", 7, addedItem.Qty)
	}
}

func TestUpdateItemInCartEmptyCartItems(t *testing.T) {
	userCart := NewUserCart(u, c)
	err := userCart.UpdateItemInCart(&vo.CartItem{
		ProdID:   "001",
		ProdName: "Shuriken",
		Qty:      5,
		Price:    125.5,
		Disc:     0.0,
	})
	if err == nil {
		t.Errorf("Expected got error of %v, but no error is thrown", e.ErrNoData{})
	}
}

func TestUpdateItemInCartNonExistingItem(t *testing.T) {
	userCart := NewUserCart(u, c)
	userCart.AddItemToCart(p, 5)
	err := userCart.UpdateItemInCart(&vo.CartItem{
		ProdID:   "002",
		ProdName: "Katana",
		Qty:      3,
		Price:    750.75,
		Disc:     0.0,
	})
	if err == nil {
		t.Errorf("Expected got error of %v, but no error is thrown", e.ErrNoData{})
	}
}

func TestUpdateItemInCart(t *testing.T) {
	userCart := NewUserCart(u, c)
	userCart.AddItemToCart(p, 5)
	err := userCart.UpdateItemInCart(&vo.CartItem{
		ProdID:   p.ID,
		ProdName: p.Name,
		Qty:      3,
		Price:    p.Price,
		Disc:     p.Disc,
	})
	if err != nil {
		t.Errorf("Expected no error, but got error of %v", err)
	}
}

func TestRemoveItemFromCartEmptyCartItems(t *testing.T) {
	userCart := NewUserCart(u, c)
	err := userCart.RemoveItemFromCart("001")
	if err == nil {
		t.Errorf("Expected got error of %v, but no error is thrown", e.ErrNoData{})
	}
}

func TestRemoveItemFromCartNonExistingItem(t *testing.T) {
	userCart := NewUserCart(u, c)
	userCart.AddItemToCart(p, 5)
	err := userCart.RemoveItemFromCart("002")
	if err == nil {
		t.Errorf("Expected got error of %v, but no error is thrown", e.ErrNoData{})
	}
}

func TestRemoveItemFromCart(t *testing.T) {
	userCart := NewUserCart(u, c)
	userCart.AddItemToCart(p, 5)
	err := userCart.RemoveItemFromCart(p.ID)
	if err != nil {
		t.Errorf("Expected no error, but got error of %v", err)
	}
}

func TestFetchUserInfo(t *testing.T) {
	userCart := NewUserCart(u, c)
	user := userCart.FetchUserInfo()
	if user != u {
		t.Errorf("Expected user of %#v, got %#v", u, user)
	}
}

func TestFetchCartInfo(t *testing.T) {
	userCart := NewUserCart(u, c)
	cart := userCart.FetchCartInfo()
	if cart != c {
		t.Errorf("Expected cart of %#v, got %#v", c, cart)
	}
}
