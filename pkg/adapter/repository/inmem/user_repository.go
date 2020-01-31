package inmem

import (
	"github.com/yauritux/cartsvc/pkg/adapter/repository/inmem/model"
	uc "github.com/yauritux/cartsvc/pkg/usecase"
)

type UserRepository struct {
	data []*model.User
}

func NewUserRepository() *UserRepository {
	userRecords := make([]*model.User, 0)
	user := &model.User{
		ID:    "yauritux",
		Name:  "Yauri Attamimi",
		Phone: "+62822xxxxxx",
		Email: "yauritux@gmail.com",
		BillingAddress: &model.Address{
			StreetName:  "Kalibata Raya No.1",
			City:        "Jakarta",
			Region:      "South Jakarta",
			Province:    "DKI Jakarta",
			Postal:      "12750",
			Country:     "Indonesia",
			AddressType: "billing_address",
		},
		ShippingAddress: &model.Address{
			StreetName:  "Kalibata Raya No.1",
			City:        "Jakarta",
			Region:      "DKI Jakarta",
			Postal:      "12750",
			Country:     "Indonesia",
			AddressType: "shipping_address",
		},
	}
	userRecords = append(userRecords, user)
	return &UserRepository{data: userRecords}
}

func (r *UserRepository) FindByUserID(uid string) (interface{}, error) {
	if uid == "" {
		return nil, nil
	}
	for i, u := range r.data {
		if u.ID != uid {
			continue
		}
		return r.BuildUserUsecaseModel(r.data[i]), nil
	}
	return nil, nil
}

func (r *UserRepository) BuildUserUsecaseModel(user interface{}) *uc.User {
	switch user.(type) {
	case *model.User:
		u := user.(*model.User)
		return &uc.User{
			ID:       u.ID,
			Username: u.Name,
			Phone:    u.Phone,
			Email:    u.Email,
			BillingAddr: &uc.Address{
				Street:      u.BillingAddress.StreetName,
				City:        u.BillingAddress.City,
				Postal:      u.BillingAddress.Postal,
				Province:    u.BillingAddress.Province,
				Region:      u.BillingAddress.Region,
				Country:     u.BillingAddress.Country,
				AddressType: u.BillingAddress.AddressType,
			},
			ShippingAddr: &uc.Address{
				Street:      u.ShippingAddress.StreetName,
				City:        u.ShippingAddress.City,
				Postal:      u.ShippingAddress.Postal,
				Province:    u.ShippingAddress.Province,
				Region:      u.ShippingAddress.Region,
				Country:     u.ShippingAddress.Country,
				AddressType: u.ShippingAddress.AddressType,
			},
		}
	default:
		return nil
	}
}
