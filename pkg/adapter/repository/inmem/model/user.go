package model

import . "github.com/yauritux/cartsvc/pkg/sharedkernel/enum"

type User struct {
	ID              string
	Name            string
	Phone           string
	Email           string
	BillingAddress  *Address
	ShippingAddress *Address
}

type Address struct {
	StreetName  string
	City        string
	Postal      string
	Province    string
	Region      string
	Country     string
	AddressType AddressType
}
