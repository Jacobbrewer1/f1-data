// Code generated by mockery v2.46.1. DO NOT EDIT.

package vaulty

import (
	api "github.com/hashicorp/vault/api"
	mock "github.com/stretchr/testify/mock"
)

// mockLoginFunc is an autogenerated mock type for the loginFunc type
type mockLoginFunc struct {
	mock.Mock
}

// Execute provides a mock function with given fields: v
func (_m *mockLoginFunc) Execute(v *api.Client) (*api.Secret, error) {
	ret := _m.Called(v)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 *api.Secret
	var r1 error
	if rf, ok := ret.Get(0).(func(*api.Client) (*api.Secret, error)); ok {
		return rf(v)
	}
	if rf, ok := ret.Get(0).(func(*api.Client) *api.Secret); ok {
		r0 = rf(v)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.Secret)
		}
	}

	if rf, ok := ret.Get(1).(func(*api.Client) error); ok {
		r1 = rf(v)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// newMockLoginFunc creates a new instance of mockLoginFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockLoginFunc(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockLoginFunc {
	mock := &mockLoginFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}