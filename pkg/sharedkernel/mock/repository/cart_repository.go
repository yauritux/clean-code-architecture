package repository

import (
	"github.com/stretchr/testify/mock"
)

type MockCartRepository struct {
	mock.Mock
}

func (m *MockCartRepository) FetchUserCart(userID string) (interface{}, error) {
	call := m.Called(userID)
	res := call.Get(0)
	if res == nil {
		return nil, call.Error(1)
	}
	return res, nil
}

func (m *MockCartRepository) AddToCart(cartID string, item interface{}) error {
	call := m.Called(cartID, item)
	return call.Error(0)
}

func (m *MockCartRepository) RemoveItem(cartID string, itemID string) error {
	call := m.Called(cartID, itemID)
	return call.Error(0)
}

func (m *MockCartRepository) UpdateItem(cartID string, item interface{}) error {
	call := m.Called(cartID, item)
	return call.Error(0)
}

func (m *MockCartRepository) Checkout(cartID string) interface{} {
	call := m.Called(cartID)
	return call.Get(0)
}

func (m *MockCartRepository) Canceled(cartID string) error {
	call := m.Called(cartID)
	return call.Error(0)
}

func (m *MockCartRepository) Close(cartID string) error {
	call := m.Called(cartID)
	return call.Error(0)
}
