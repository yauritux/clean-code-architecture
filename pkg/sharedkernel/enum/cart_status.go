package enum

type CartStatus string

const (
	Open              CartStatus = "open"
	PaymentProcessing CartStatus = "payment_processing"
	Canceled          CartStatus = "canceled"
	Closed            CartStatus = "closed"
)
