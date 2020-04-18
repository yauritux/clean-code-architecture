package carts

import (
	"errors"
	"reflect"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
	"github.com/yauritux/cartsvc/pkg/domain/entity"
	"github.com/yauritux/cartsvc/pkg/sharedkernel/enum"
	mockRepo "github.com/yauritux/cartsvc/pkg/sharedkernel/mock/repository"
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
}
