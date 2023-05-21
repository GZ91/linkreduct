// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/GZ91/linkreduct/internal/models"
	mock "github.com/stretchr/testify/mock"
)

// Storeger is an autogenerated mock type for the Storeger type
type Storeger struct {
	mock.Mock
}

type Storeger_Expecter struct {
	mock *mock.Mock
}

func (_m *Storeger) EXPECT() *Storeger_Expecter {
	return &Storeger_Expecter{mock: &_m.Mock}
}

// AddBatchLink provides a mock function with given fields: _a0, _a1
func (_m *Storeger) AddBatchLink(_a0 context.Context, _a1 []models.IncomingBatchURL) ([]models.ReleasedBatchURL, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []models.ReleasedBatchURL
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []models.IncomingBatchURL) ([]models.ReleasedBatchURL, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []models.IncomingBatchURL) []models.ReleasedBatchURL); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.ReleasedBatchURL)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []models.IncomingBatchURL) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Storeger_AddBatchLink_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddBatchLink'
type Storeger_AddBatchLink_Call struct {
	*mock.Call
}

// AddBatchLink is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 []models.IncomingBatchURL
func (_e *Storeger_Expecter) AddBatchLink(_a0 interface{}, _a1 interface{}) *Storeger_AddBatchLink_Call {
	return &Storeger_AddBatchLink_Call{Call: _e.mock.On("AddBatchLink", _a0, _a1)}
}

func (_c *Storeger_AddBatchLink_Call) Run(run func(_a0 context.Context, _a1 []models.IncomingBatchURL)) *Storeger_AddBatchLink_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]models.IncomingBatchURL))
	})
	return _c
}

func (_c *Storeger_AddBatchLink_Call) Return(_a0 []models.ReleasedBatchURL, _a1 error) *Storeger_AddBatchLink_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Storeger_AddBatchLink_Call) RunAndReturn(run func(context.Context, []models.IncomingBatchURL) ([]models.ReleasedBatchURL, error)) *Storeger_AddBatchLink_Call {
	_c.Call.Return(run)
	return _c
}

// AddURL provides a mock function with given fields: _a0, _a1
func (_m *Storeger) AddURL(_a0 context.Context, _a1 string) (string, error) {
	ret := _m.Called(_a0, _a1)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Storeger_AddURL_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddURL'
type Storeger_AddURL_Call struct {
	*mock.Call
}

// AddURL is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 string
func (_e *Storeger_Expecter) AddURL(_a0 interface{}, _a1 interface{}) *Storeger_AddURL_Call {
	return &Storeger_AddURL_Call{Call: _e.mock.On("AddURL", _a0, _a1)}
}

func (_c *Storeger_AddURL_Call) Run(run func(_a0 context.Context, _a1 string)) *Storeger_AddURL_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *Storeger_AddURL_Call) Return(_a0 string, _a1 error) *Storeger_AddURL_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Storeger_AddURL_Call) RunAndReturn(run func(context.Context, string) (string, error)) *Storeger_AddURL_Call {
	_c.Call.Return(run)
	return _c
}

// FindLongURL provides a mock function with given fields: _a0, _a1
func (_m *Storeger) FindLongURL(_a0 context.Context, _a1 string) (string, bool, error) {
	ret := _m.Called(_a0, _a1)

	var r0 string
	var r1 bool
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, bool, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) bool); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Get(1).(bool)
	}

	if rf, ok := ret.Get(2).(func(context.Context, string) error); ok {
		r2 = rf(_a0, _a1)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Storeger_FindLongURL_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindLongURL'
type Storeger_FindLongURL_Call struct {
	*mock.Call
}

// FindLongURL is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 string
func (_e *Storeger_Expecter) FindLongURL(_a0 interface{}, _a1 interface{}) *Storeger_FindLongURL_Call {
	return &Storeger_FindLongURL_Call{Call: _e.mock.On("FindLongURL", _a0, _a1)}
}

func (_c *Storeger_FindLongURL_Call) Run(run func(_a0 context.Context, _a1 string)) *Storeger_FindLongURL_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *Storeger_FindLongURL_Call) Return(_a0 string, _a1 bool, _a2 error) *Storeger_FindLongURL_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *Storeger_FindLongURL_Call) RunAndReturn(run func(context.Context, string) (string, bool, error)) *Storeger_FindLongURL_Call {
	_c.Call.Return(run)
	return _c
}

// GetURL provides a mock function with given fields: _a0, _a1
func (_m *Storeger) GetURL(_a0 context.Context, _a1 string) (string, bool, error) {
	ret := _m.Called(_a0, _a1)

	var r0 string
	var r1 bool
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, bool, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) bool); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Get(1).(bool)
	}

	if rf, ok := ret.Get(2).(func(context.Context, string) error); ok {
		r2 = rf(_a0, _a1)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Storeger_GetURL_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetURL'
type Storeger_GetURL_Call struct {
	*mock.Call
}

// GetURL is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 string
func (_e *Storeger_Expecter) GetURL(_a0 interface{}, _a1 interface{}) *Storeger_GetURL_Call {
	return &Storeger_GetURL_Call{Call: _e.mock.On("GetURL", _a0, _a1)}
}

func (_c *Storeger_GetURL_Call) Run(run func(_a0 context.Context, _a1 string)) *Storeger_GetURL_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *Storeger_GetURL_Call) Return(_a0 string, _a1 bool, _a2 error) *Storeger_GetURL_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *Storeger_GetURL_Call) RunAndReturn(run func(context.Context, string) (string, bool, error)) *Storeger_GetURL_Call {
	_c.Call.Return(run)
	return _c
}

// Ping provides a mock function with given fields: _a0
func (_m *Storeger) Ping(_a0 context.Context) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Storeger_Ping_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Ping'
type Storeger_Ping_Call struct {
	*mock.Call
}

// Ping is a helper method to define mock.On call
//   - _a0 context.Context
func (_e *Storeger_Expecter) Ping(_a0 interface{}) *Storeger_Ping_Call {
	return &Storeger_Ping_Call{Call: _e.mock.On("Ping", _a0)}
}

func (_c *Storeger_Ping_Call) Run(run func(_a0 context.Context)) *Storeger_Ping_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *Storeger_Ping_Call) Return(_a0 error) *Storeger_Ping_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Storeger_Ping_Call) RunAndReturn(run func(context.Context) error) *Storeger_Ping_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewStoreger interface {
	mock.TestingT
	Cleanup(func())
}

// NewStoreger creates a new instance of Storeger. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewStoreger(t mockConstructorTestingTNewStoreger) *Storeger {
	mock := &Storeger{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
