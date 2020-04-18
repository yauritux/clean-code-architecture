package users

import (
	"errors"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
	"github.com/yauritux/cartsvc/pkg/domain/entity"
	mockUserRepo "github.com/yauritux/cartsvc/pkg/sharedkernel/mock/repository"
)

func TestUserUsecase(t *testing.T) {

	Convey("1. When searching a user by ID", t, func() {

		userRepo := &mockUserRepo.MockUserRepository{}

		Convey("-> Negative Scenarios", func() {
			Convey("-> Some errors occures within the system repository", func() {
				Convey("-> Should return an error with a message related to a repository error", func() {
					userRepo.On("FindByUserID", mock.Anything).Return(nil, errors.New("Database error"))
					uc := NewUserUsecase(userRepo)
					res, err := uc.FetchCurrentUser(mock.Anything)
					So(res, ShouldBeNil)
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, "Database error")
				})
			})
			Convey("-> An error occured due to a wrong user system type", func() {
				Convey("-> Should return an error with a message of wrong user type", func() {
					userRepo.On("FindByUserID", mock.Anything).Return(&entity.User{}, nil)
					uc := NewUserUsecase(userRepo)
					res, err := uc.FetchCurrentUser(mock.Anything)
					So(res, ShouldBeNil)
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, "cannot fetch current user, got an invalid user type returned from the repository")
				})
			})
		})

		Convey("-> Positive Scenarios", func() {
			Convey("-> User is found", func() {
				Convey("-> Should return a user with no error", func() {
					sUser := &User{
						ID:       "123",
						Username: "yauritux",
					}
					userRepo.On("FindByUserID", "123").Return(sUser, nil)
					uc := NewUserUsecase(userRepo)
					res, err := uc.FetchCurrentUser("123")
					So(res, ShouldNotBeNil)
					So(err, ShouldBeNil)
					So(reflect.DeepEqual(sUser, res), ShouldBeTrue)
				})
			})
		})
	})
}
