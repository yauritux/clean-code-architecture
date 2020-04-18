package users

import (
	"errors"

	"github.com/yauritux/cartsvc/pkg/domain/repository"
	. "github.com/yauritux/cartsvc/pkg/sharedkernel/enum"
)

type UserUsecase struct {
	repo repository.UserRepository
}

type User struct {
	ID           string
	Username     string
	Email        string
	Phone        string
	BillingAddr  *Address
	ShippingAddr *Address
}

type Address struct {
	Street      string
	City        string
	Postal      string
	Province    string
	Region      string
	Country     string
	AddressType AddressType
}

func NewUserUsecase(r repository.UserRepository) *UserUsecase {
	return &UserUsecase{r}
}

func (user *UserUsecase) FetchCurrentUser(id string) (interface{}, error) {
	u, err := user.repo.FindByUserID(id)
	if err != nil {
		return nil, err
	}

	currUser, ok := u.(*User)
	if !ok {
		return nil, errors.New("cannot fetch current user, got an invalid user type returned from the repository")
	}

	return currUser, nil
}
