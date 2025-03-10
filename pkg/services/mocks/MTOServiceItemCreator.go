// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	appcontext "github.com/transcom/mymove/pkg/appcontext"

	models "github.com/transcom/mymove/pkg/models"

	validate "github.com/gobuffalo/validate/v3"
)

// MTOServiceItemCreator is an autogenerated mock type for the MTOServiceItemCreator type
type MTOServiceItemCreator struct {
	mock.Mock
}

// CreateMTOServiceItem provides a mock function with given fields: appCtx, serviceItem
func (_m *MTOServiceItemCreator) CreateMTOServiceItem(appCtx appcontext.AppContext, serviceItem *models.MTOServiceItem) (*models.MTOServiceItems, *validate.Errors, error) {
	ret := _m.Called(appCtx, serviceItem)

	var r0 *models.MTOServiceItems
	var r1 *validate.Errors
	var r2 error
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, *models.MTOServiceItem) (*models.MTOServiceItems, *validate.Errors, error)); ok {
		return rf(appCtx, serviceItem)
	}
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, *models.MTOServiceItem) *models.MTOServiceItems); ok {
		r0 = rf(appCtx, serviceItem)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.MTOServiceItems)
		}
	}

	if rf, ok := ret.Get(1).(func(appcontext.AppContext, *models.MTOServiceItem) *validate.Errors); ok {
		r1 = rf(appCtx, serviceItem)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*validate.Errors)
		}
	}

	if rf, ok := ret.Get(2).(func(appcontext.AppContext, *models.MTOServiceItem) error); ok {
		r2 = rf(appCtx, serviceItem)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

type mockConstructorTestingTNewMTOServiceItemCreator interface {
	mock.TestingT
	Cleanup(func())
}

// NewMTOServiceItemCreator creates a new instance of MTOServiceItemCreator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMTOServiceItemCreator(t mockConstructorTestingTNewMTOServiceItemCreator) *MTOServiceItemCreator {
	mock := &MTOServiceItemCreator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
