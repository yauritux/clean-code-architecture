package port

import (
	uc "github.com/yauritux/cartsvc/pkg/usecase"
)

type UserInputPort interface {
	FetchCurrentUser(id string) (interface{}, error)
	BuildUserUsecaseModel(interface{}) *uc.User
}

type UserOutputPort interface {
}
