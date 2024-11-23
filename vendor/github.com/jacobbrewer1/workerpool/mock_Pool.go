// Code generated by mockery. DO NOT EDIT.

package workerpool

import mock "github.com/stretchr/testify/mock"

// MockPool is an autogenerated mock type for the Pool type
type MockPool struct {
	mock.Mock
}

// BlockingSchedule provides a mock function with given fields: task
func (_m *MockPool) BlockingSchedule(task Runnable) error {
	ret := _m.Called(task)

	if len(ret) == 0 {
		panic("no return value specified for BlockingSchedule")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(Runnable) error); ok {
		r0 = rf(task)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MustSchedule provides a mock function with given fields: task
func (_m *MockPool) MustSchedule(task Runnable) {
	_m.Called(task)
}

// Schedule provides a mock function with given fields: task
func (_m *MockPool) Schedule(task Runnable) error {
	ret := _m.Called(task)

	if len(ret) == 0 {
		panic("no return value specified for Schedule")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(Runnable) error); ok {
		r0 = rf(task)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Stop provides a mock function with given fields:
func (_m *MockPool) Stop() {
	_m.Called()
}

// StopAsync provides a mock function with given fields:
func (_m *MockPool) StopAsync() {
	_m.Called()
}

// NewMockPool creates a new instance of MockPool. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockPool(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockPool {
	mock := &MockPool{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
