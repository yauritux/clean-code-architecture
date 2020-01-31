package entity

import (
	"time"

	vo "github.com/yauritux/cartsvc/pkg/domain/valueobject"
	. "github.com/yauritux/cartsvc/pkg/sharedkernel/enum"
)

type Cart struct {
	ID         string
	UserID     string
	Status     CartStatus
	Items      []*vo.CartItem
	CreatedAt  time.Time
	CanceledAt *time.Time
}
