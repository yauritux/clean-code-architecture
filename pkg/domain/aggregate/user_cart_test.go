package aggregate

import (
	"testing"
	"time"

	"github.com/yauritux/cartsvc/pkg/domain/entity"
	vo "github.com/yauritux/cartsvc/pkg/domain/valueobject"
	"github.com/yauritux/cartsvc/pkg/sharedkernel/enum"

	. "github.com/smartystreets/goconvey/convey"
)

var u *entity.User
var c *entity.Cart
var p *entity.Product

func setup() {
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

func TestSpec(t *testing.T) {

	Convey("1. Given a new cart object", t, func() {
		setup()
		userCart := NewUserCart(u, c)
		Convey("-> Cart initial status should equal to Open", func() {
			So(userCart.cart.Status, ShouldEqual, enum.Open)
		})
		Convey("-> Cart user should not be empty", func() {
			So(userCart.user, ShouldNotBeEmpty)
			So(userCart.user, ShouldEqual, u)
		})
		Convey("-> Cart object should not be empty", func() {
			So(userCart.cart, ShouldNotBeEmpty)
			So(userCart.cart, ShouldEqual, c)
		})
	})

	Convey("2. Given adding an item to cart", t, func() {

		Convey("-> Negative Scenarios", func() {
			setup()
			userCart := NewUserCart(u, c)
			Convey("-> When item is out of stock", func() {
				Convey("-> Should return error", func() {
					res, err := userCart.AddItemToCart(p, 105)
					So(res, ShouldBeNil)
					So(err, ShouldNotBeEmpty)
				})
			})
			Convey("-> When cart session id is missing", func() {
				Convey("-> Should return error", func() {
					c.ID = ""
					res, err := userCart.AddItemToCart(p, 10)
					So(res, ShouldBeNil)
					So(err, ShouldNotBeEmpty)
				})
			})
			Convey("-> When user id is missing", func() {
				Convey("-> Should return error", func() {
					c.UserID = ""
					res, err := userCart.AddItemToCart(p, 5)
					So(res, ShouldBeNil)
					So(err, ShouldNotBeEmpty)
				})
			})
			Convey("-> When cart is closed, means all items within cart is settled/paid", func() {
				Convey("-> Should return error", func() {
					c.Status = enum.Closed
					res, err := userCart.AddItemToCart(p, 5)
					So(res, ShouldBeNil)
					So(err, ShouldNotBeEmpty)
				})
			})
		})

		Convey("-> Positive Scenarios", func() {
			setup()
			userCart := NewUserCart(u, c)
			Convey("-> When new item is added", func() {
				Convey("-> requested item should exist within the cart", func() {
					addedItem, err := userCart.AddItemToCart(p, 5)
					So(err, ShouldBeEmpty)
					So(addedItem.ProdID, ShouldEqual, p.ID)
					So(addedItem.Qty == 5, ShouldBeTrue)
				})
			})
			Convey("-> When the same item already exist inside the cart", func() {
				userCart.AddItemToCart(p, 5)
				Convey("-> should update the existing cart's item quantity with the one being added", func() {
					addedItem, err := userCart.AddItemToCart(p, 2)
					So(err, ShouldNotBeEmpty)
					So(addedItem.Qty == 7, ShouldBeTrue)
				})
			})
		})
	})

	Convey("3. Given update the cart item", t, func() {

		Convey("-> Negative Scenarios", func() {
			setup()
			userCart := NewUserCart(u, c)

			Convey("-> When cart is empty", func() {
				Convey("-> Should return error with message cart still empty", func() {
					err := userCart.UpdateItemInCart(
						&vo.CartItem{
							ProdID:   "001",
							ProdName: "Shuriken",
							Qty:      5,
							Price:    125.5,
							Disc:     0.0,
						},
					)
					So(err, ShouldNotBeEmpty)
					So(err.Error(), ShouldEqual, "cart is still empty")
				})
			})

			Convey("-> Item does not exist in the cart", func() {
				Convey("-> Should return error with message cannot find cart item", func() {
					userCart.AddItemToCart(p, 5)
					err := userCart.UpdateItemInCart(&vo.CartItem{
						ProdID:   "002",
						ProdName: "Katana",
						Qty:      3,
						Price:    750.75,
						Disc:     0.0,
					})
					So(err, ShouldNotBeEmpty)
					So(err.Error(), ShouldEqual, "cannot find cart item 002 - Katana within the cart")
				})
			})
		})

		Convey("-> Positive Scenarios", func() {
			setup()
			userCart := NewUserCart(u, c)

			Convey("-> No error should be thrown", func() {
				userCart.AddItemToCart(p, 5)
				err := userCart.UpdateItemInCart(&vo.CartItem{
					ProdID:   p.ID,
					ProdName: p.Name,
					Qty:      3,
					Price:    p.Price,
					Disc:     p.Disc,
				})
				So(err, ShouldBeEmpty)
			})
		})
	})

	Convey("4. Given remove item from cart", t, func() {

		Convey("-> Negative Scenarios", func() {
			setup()
			userCart := NewUserCart(u, c)

			Convey("-> When cart is empty", func() {
				Convey("-> Should got an error with message cart is still empty", func() {
					err := userCart.RemoveItemFromCart("001")
					So(err, ShouldNotBeEmpty)
					So(err.Error(), ShouldEqual, "cart is still empty")
				})
			})

			Convey("-> When trying to remove an unexisting item from the cart", func() {
				Convey("-> Should got an error with message cannot find the cart item", func() {
					userCart.AddItemToCart(p, 5)
					err := userCart.RemoveItemFromCart("002")
					So(err, ShouldNotBeEmpty)
					So(err.Error(), ShouldEqual, "cannot find cart item with ID 002")
				})
			})
		})

		Convey("-> Positive Scenarios", func() {
			setup()
			userCart := NewUserCart(u, c)

			Convey("-> The item should be removed from the cart", func() {
				userCart.AddItemToCart(p, 5)
				err := userCart.RemoveItemFromCart(p.ID)
				So(err, ShouldBeEmpty)
			})
		})
	})
}
