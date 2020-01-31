package valueobject

import (
	. "github.com/yauritux/cartsvc/pkg/sharedkernel/enum"
)

type BuyerAddress struct {
	StreetName string
	City       string
	Postal     string
	Province   string
	Region     string
	Country    string
	Type       AddressType
}
