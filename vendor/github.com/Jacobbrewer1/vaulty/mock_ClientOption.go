// Code generated by mockery. DO NOT EDIT.

package vaulty

import mock "github.com/stretchr/testify/mock"

// MockClientOption is an autogenerated mock type for the ClientOption type
type MockClientOption struct {
	mock.Mock
}

// Execute provides a mock function with given fields: c
func (_m *MockClientOption) Execute(c *client) {
	_m.Called(c)
}

// NewMockClientOption creates a new instance of MockClientOption. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockClientOption(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockClientOption {
	mock := &MockClientOption{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
