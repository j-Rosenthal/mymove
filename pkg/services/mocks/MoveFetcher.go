// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	appcontext "github.com/transcom/mymove/pkg/appcontext"

	models "github.com/transcom/mymove/pkg/models"

	services "github.com/transcom/mymove/pkg/services"
)

// MoveFetcher is an autogenerated mock type for the MoveFetcher type
type MoveFetcher struct {
	mock.Mock
}

// FetchMove provides a mock function with given fields: appCtx, locator, searchParams
func (_m *MoveFetcher) FetchMove(appCtx appcontext.AppContext, locator string, searchParams *services.MoveFetcherParams) (*models.Move, error) {
	ret := _m.Called(appCtx, locator, searchParams)

	var r0 *models.Move
	var r1 error
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, string, *services.MoveFetcherParams) (*models.Move, error)); ok {
		return rf(appCtx, locator, searchParams)
	}
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, string, *services.MoveFetcherParams) *models.Move); ok {
		r0 = rf(appCtx, locator, searchParams)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Move)
		}
	}

	if rf, ok := ret.Get(1).(func(appcontext.AppContext, string, *services.MoveFetcherParams) error); ok {
		r1 = rf(appCtx, locator, searchParams)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMoveFetcher interface {
	mock.TestingT
	Cleanup(func())
}

// NewMoveFetcher creates a new instance of MoveFetcher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMoveFetcher(t mockConstructorTestingTNewMoveFetcher) *MoveFetcher {
	mock := &MoveFetcher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
