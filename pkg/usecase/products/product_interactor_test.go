package products

import (
	"errors"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
	"github.com/yauritux/cartsvc/pkg/domain/entity"
	mockProductRepo "github.com/yauritux/cartsvc/pkg/sharedkernel/mock/repository"
)

func TestProductUsecase(t *testing.T) {

	Convey("1. Given a user is searching for a product by ID", t, func() {

		prodRepo := &mockProductRepo.MockProductRepository{}

		Convey("-> Negative Scenarios", func() {
			Convey("-> Some errors occured within the system repository", func() {
				Convey("-> Should return an error with a message related to a repository error", func() {
					prodRepo.On("FindByProductID", mock.Anything).Return(nil, errors.New("Database error"))
					uc := NewProductUsecase(prodRepo)
					res, err := uc.FindByProductID("001")
					So(res, ShouldBeNil)
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, "Database error")
				})
			})
			Convey("-> Some errors occured due to invalid product system type", func() {
				Convey("-> Should return an error with an error message of invalid product system type", func() {
					prodRepo.On("FindByProductID", mock.Anything).Return(&entity.Product{}, nil)
					uc := NewProductUsecase(prodRepo)
					res, err := uc.FindByProductID("001")
					So(res, ShouldBeNil)
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, "cannot find product with ID 001, got an invalid product type returned from the repository")
				})
			})
		})

		Convey("-> Positive Scenarios", func() {
			Convey("-> Found the product", func() {
				Convey("-> Should return the product to the user without no error", func() {
					sProduct := &Product{
						ID:    "001",
						Name:  "Shuriken",
						Stock: 999,
						Price: 1250,
						Disc:  0,
					}
					prodRepo.On("FindByProductID", "001").Return(sProduct, nil)
					uc := NewProductUsecase(prodRepo)
					res, err := uc.FindByProductID("001")
					So(res, ShouldNotBeNil)
					So(err, ShouldBeNil)
					So(reflect.DeepEqual(res, sProduct), ShouldBeTrue)
				})
			})
		})
	})
}
