package model

import (
	. "github.com/yauritux/cartsvc/pkg/sharedkernel/enum"

	"time"
)

type Cart struct {
	ID         string
	UserID     string
	Status     CartStatus
	Items      []*CartItem
	CreatedAt  time.Time
	CanceledAt *time.Time
}

type CartItem struct {
	ID    string
	Name  string
	Qty   int
	Price float64
	Disc  float64
}
