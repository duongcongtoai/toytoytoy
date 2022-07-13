// Code generated by mockery v2.14.0. DO NOT EDIT.

package services

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	sql "database/sql"
)

// DBX is an autogenerated mock type for the DBX type
type DBX struct {
	mock.Mock
}

// Begin provides a mock function with given fields:
func (_m *DBX) Begin() (*sql.Tx, error) {
	ret := _m.Called()

	var r0 *sql.Tx
	if rf, ok := ret.Get(0).(func() *sql.Tx); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sql.Tx)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BeginTx provides a mock function with given fields: ctx, opts
func (_m *DBX) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	ret := _m.Called(ctx, opts)

	var r0 *sql.Tx
	if rf, ok := ret.Get(0).(func(context.Context, *sql.TxOptions) *sql.Tx); ok {
		r0 = rf(ctx, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sql.Tx)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *sql.TxOptions) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ExecContext provides a mock function with given fields: _a0, _a1, _a2
func (_m *DBX) ExecContext(_a0 context.Context, _a1 string, _a2 ...interface{}) (sql.Result, error) {
	var _ca []interface{}
	_ca = append(_ca, _a0, _a1)
	_ca = append(_ca, _a2...)
	ret := _m.Called(_ca...)

	var r0 sql.Result
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) sql.Result); ok {
		r0 = rf(_a0, _a1, _a2...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(sql.Result)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, ...interface{}) error); ok {
		r1 = rf(_a0, _a1, _a2...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PrepareContext provides a mock function with given fields: _a0, _a1
func (_m *DBX) PrepareContext(_a0 context.Context, _a1 string) (*sql.Stmt, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *sql.Stmt
	if rf, ok := ret.Get(0).(func(context.Context, string) *sql.Stmt); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sql.Stmt)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QueryContext provides a mock function with given fields: _a0, _a1, _a2
func (_m *DBX) QueryContext(_a0 context.Context, _a1 string, _a2 ...interface{}) (*sql.Rows, error) {
	var _ca []interface{}
	_ca = append(_ca, _a0, _a1)
	_ca = append(_ca, _a2...)
	ret := _m.Called(_ca...)

	var r0 *sql.Rows
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) *sql.Rows); ok {
		r0 = rf(_a0, _a1, _a2...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sql.Rows)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, ...interface{}) error); ok {
		r1 = rf(_a0, _a1, _a2...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QueryRowContext provides a mock function with given fields: _a0, _a1, _a2
func (_m *DBX) QueryRowContext(_a0 context.Context, _a1 string, _a2 ...interface{}) *sql.Row {
	var _ca []interface{}
	_ca = append(_ca, _a0, _a1)
	_ca = append(_ca, _a2...)
	ret := _m.Called(_ca...)

	var r0 *sql.Row
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) *sql.Row); ok {
		r0 = rf(_a0, _a1, _a2...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sql.Row)
		}
	}

	return r0
}

type mockConstructorTestingTNewDBX interface {
	mock.TestingT
	Cleanup(func())
}

// NewDBX creates a new instance of DBX. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDBX(t mockConstructorTestingTNewDBX) *DBX {
	mock := &DBX{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}