package entity

import (
	vo "github.com/yauritux/cartsvc/pkg/domain/valueobject"
)

type User struct {
	UserID          string
	Username        string
	Phone           string
	Email           string
	BillingAddress  vo.BuyerAddress
	ShippingAddress vo.BuyerAddress
}
