// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	appcontext "github.com/transcom/mymove/pkg/appcontext"

	models "github.com/transcom/mymove/pkg/models"

	services "github.com/transcom/mymove/pkg/services"

	validate "github.com/gobuffalo/validate/v3"
)

// AdminUserCreator is an autogenerated mock type for the AdminUserCreator type
type AdminUserCreator struct {
	mock.Mock
}

// CreateAdminUser provides a mock function with given fields: appCtx, user, organizationIDFilter
func (_m *AdminUserCreator) CreateAdminUser(appCtx appcontext.AppContext, user *models.AdminUser, organizationIDFilter []services.QueryFilter) (*models.AdminUser, *validate.Errors, error) {
	ret := _m.Called(appCtx, user, organizationIDFilter)

	var r0 *models.AdminUser
	var r1 *validate.Errors
	var r2 error
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, *models.AdminUser, []services.QueryFilter) (*models.AdminUser, *validate.Errors, error)); ok {
		return rf(appCtx, user, organizationIDFilter)
	}
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, *models.AdminUser, []services.QueryFilter) *models.AdminUser); ok {
		r0 = rf(appCtx, user, organizationIDFilter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.AdminUser)
		}
	}

	if rf, ok := ret.Get(1).(func(appcontext.AppContext, *models.AdminUser, []services.QueryFilter) *validate.Errors); ok {
		r1 = rf(appCtx, user, organizationIDFilter)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*validate.Errors)
		}
	}

	if rf, ok := ret.Get(2).(func(appcontext.AppContext, *models.AdminUser, []services.QueryFilter) error); ok {
		r2 = rf(appCtx, user, organizationIDFilter)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

type mockConstructorTestingTNewAdminUserCreator interface {
	mock.TestingT
	Cleanup(func())
}

// NewAdminUserCreator creates a new instance of AdminUserCreator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAdminUserCreator(t mockConstructorTestingTNewAdminUserCreator) *AdminUserCreator {
	mock := &AdminUserCreator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
