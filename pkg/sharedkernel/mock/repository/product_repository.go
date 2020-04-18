package repository

import (
	"github.com/stretchr/testify/mock"
)

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) FindByProductID(id string) (interface{}, error) {
	call := m.Called(id)
	res := call.Get(0)
	if res == nil {
		return nil, call.Error(1)
	}
	return res, nil
}
