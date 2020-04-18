package repository

import (
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByUserID(uid string) (interface{}, error) {
	call := m.Called(uid)
	res := call.Get(0)
	if res == nil {
		return nil, call.Error(1)
	}
	return res, nil
}
