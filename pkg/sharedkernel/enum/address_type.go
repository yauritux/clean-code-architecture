package enum

type AddressType string

const (
	BillingAddress  AddressType = "billing_address"
	ShippingAddress AddressType = "shipping_address"
)
