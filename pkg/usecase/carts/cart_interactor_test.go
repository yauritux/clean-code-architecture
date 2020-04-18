package carts

import (
	"errors"
	"reflect"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
	"github.com/yauritux/cartsvc/pkg/domain/entity"
	vo "github.com/yauritux/cartsvc/pkg/domain/valueobject"
	"github.com/yauritux/cartsvc/pkg/sharedkernel/enum"
	mockRepo "github.com/yauritux/cartsvc/pkg/sharedkernel/mock/repository"
	prodUsecase "github.com/yauritux/cartsvc/pkg/usecase/products"
)

func TestCartUsecase(t *testing.T) {

	Convey("1. Given a user get his cart", t, func() {

		cartRepo := &mockRepo.MockCartRepository{}
		prodRepo := &mockRepo.MockProductRepository{}

		Convey("-> Negative Scenarios", func() {
			Convey("-> When userID is empty, then should return an error", func() {
				uc := NewCartUsecase(cartRepo, prodRepo)
				res, err := uc.FetchUserCart("")
				So(res, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "cannot fetch user cart, 'user_id' is missing")
			})
			Convey("-> Should return an error for some errors occured within the system repository", func() {
				cartRepo.On("FetchUserCart", mock.Anything).Return(nil, errors.New("Database error"))
				uc := NewCartUsecase(cartRepo, prodRepo)
				res, err := uc.FetchUserCart("123")
				So(res, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Database error")
			})
			Convey("-> Should return an error for invalid cart system type", func() {
				cartRepo.On("FetchUserCart", "123").Return(&entity.Cart{}, nil)
				uc := NewCartUsecase(cartRepo, prodRepo)
				res, err := uc.FetchUserCart("123")
				So(res, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "cannot fetch user cart, invalid type of cart usecase model")
			})
		})

		Convey("-> Positive Scenarios", func() {
			Convey("-> Cart is found, should return the cart object with no error", func() {
				sCart := &Cart{
					ID:        "001",
					UserID:    "123",
					Status:    enum.Open,
					CreatedAt: time.Now(),
				}
				cartRepo.On("FetchUserCart", "123").Return(sCart, nil)
				uc := NewCartUsecase(cartRepo, prodRepo)
				res, err := uc.FetchUserCart("123")
				So(res, ShouldNotBeNil)
				So(err, ShouldBeNil)
				So(reflect.DeepEqual(res, sCart), ShouldBeTrue)
			})
		})
	})

	Convey("2. Given a user add an item to his cart", t, func() {

		cartRepo := &mockRepo.MockCartRepository{}
		prodRepo := &mockRepo.MockProductRepository{}

		Convey("-> Negative Scenarios", func() {
			Convey("-> Returns an error when cart cannot be fetched due to errors occured in the system repository", func() {
				cartRepo.On("FetchUserCart", mock.Anything).Return(nil, errors.New("Database error"))
				uc := NewCartUsecase(cartRepo, prodRepo)
				err := uc.AddToCart("123", &entity.Cart{})
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Database error")
			})
			Convey("-> Returns an error when a wrong type of cart is returned by the system repository", func() {
				cartRepo.On("FetchUserCart", mock.Anything).Return(&entity.Cart{}, nil)
				uc := NewCartUsecase(cartRepo, prodRepo)
				err := uc.AddToCart("123", &entity.Cart{})
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "conversion failed, invalid type of cart usecase model")
			})
			Convey("-> Returns an error when a wrong type of cart item is added", func() {
				cartRepo.On("FetchUserCart", "123").Return(&Cart{
					ID: "123", UserID: "123", Status: enum.Open, CreatedAt: time.Now(),
				}, nil)
				uc := NewCartUsecase(cartRepo, prodRepo)
				err := uc.AddToCart("123", &vo.CartItem{})
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "conversion failed, invalid type of product item usecase model")
			})
			Convey("-> Returns an error when failed to locate the item from the product repository", func() {
				cartRepo.On("FetchUserCart", "123").Return(&Cart{
					ID: "123", UserID: "123", Status: enum.Open, CreatedAt: time.Now(),
				}, nil)
				prodRepo.On("FindByProductID", "001").Return(nil, errors.New("Product not found"))
				uc := NewCartUsecase(cartRepo, prodRepo)
				err := uc.AddToCart("123", &CartItem{ID: "001", Name: "Shuriken", Price: 1250})
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Product not found")
			})
			Convey("-> Return an error when cart item returned from the product repository is invalid", func() {
				cartRepo.On("FetchUserCart", "123").Return(&Cart{
					ID: "123", UserID: "123", Status: enum.Open, CreatedAt: time.Now(),
				}, nil)
				prodRepo.On("FindByProductID", "001").Return(&entity.Product{}, nil)
				uc := NewCartUsecase(cartRepo, prodRepo)
				err := uc.AddToCart("123", &CartItem{ID: "001", Name: "Shuriken", Price: 1250})
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "conversion failed, invalid type of product usecase model")
			})
		})

		Convey("-> Positive Scenarios", func() {
			Convey("-> When everything goes normal, no error is thrown", func() {
				cartRepo.On("FetchUserCart", "123").Return(&Cart{
					ID: "123", UserID: "123", Status: enum.Open, CreatedAt: time.Now(),
				}, nil)
				prodRepo.On("FindByProductID", "001").Return(&prodUsecase.Product{
					ID: "001", Name: "Shuriken", Price: 1250, Stock: 999, Disc: 0,
				}, nil)
				cartRepo.On("AddToCart", mock.Anything, mock.Anything).Return(nil)
				uc := NewCartUsecase(cartRepo, prodRepo)
				err := uc.AddToCart("123", &CartItem{ID: "001", Name: "Shuriken", Price: 1250, Qty: 1})
				So(err, ShouldBeNil)
			})
		})
	})
}
