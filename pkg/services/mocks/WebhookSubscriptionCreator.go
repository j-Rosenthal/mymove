// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	appcontext "github.com/transcom/mymove/pkg/appcontext"

	models "github.com/transcom/mymove/pkg/models"

	validate "github.com/gobuffalo/validate/v3"
)

// WebhookSubscriptionCreator is an autogenerated mock type for the WebhookSubscriptionCreator type
type WebhookSubscriptionCreator struct {
	mock.Mock
}

// CreateWebhookSubscription provides a mock function with given fields: appCtx, subscription
func (_m *WebhookSubscriptionCreator) CreateWebhookSubscription(appCtx appcontext.AppContext, subscription *models.WebhookSubscription) (*models.WebhookSubscription, *validate.Errors, error) {
	ret := _m.Called(appCtx, subscription)

	var r0 *models.WebhookSubscription
	var r1 *validate.Errors
	var r2 error
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, *models.WebhookSubscription) (*models.WebhookSubscription, *validate.Errors, error)); ok {
		return rf(appCtx, subscription)
	}
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, *models.WebhookSubscription) *models.WebhookSubscription); ok {
		r0 = rf(appCtx, subscription)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.WebhookSubscription)
		}
	}

	if rf, ok := ret.Get(1).(func(appcontext.AppContext, *models.WebhookSubscription) *validate.Errors); ok {
		r1 = rf(appCtx, subscription)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*validate.Errors)
		}
	}

	if rf, ok := ret.Get(2).(func(appcontext.AppContext, *models.WebhookSubscription) error); ok {
		r2 = rf(appCtx, subscription)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

type mockConstructorTestingTNewWebhookSubscriptionCreator interface {
	mock.TestingT
	Cleanup(func())
}

// NewWebhookSubscriptionCreator creates a new instance of WebhookSubscriptionCreator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewWebhookSubscriptionCreator(t mockConstructorTestingTNewWebhookSubscriptionCreator) *WebhookSubscriptionCreator {
	mock := &WebhookSubscriptionCreator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
