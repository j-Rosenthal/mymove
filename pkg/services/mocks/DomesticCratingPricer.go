// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	appcontext "github.com/transcom/mymove/pkg/appcontext"

	models "github.com/transcom/mymove/pkg/models"

	services "github.com/transcom/mymove/pkg/services"

	time "time"

	unit "github.com/transcom/mymove/pkg/unit"
)

// DomesticCratingPricer is an autogenerated mock type for the DomesticCratingPricer type
type DomesticCratingPricer struct {
	mock.Mock
}

// Price provides a mock function with given fields: appCtx, contractCode, requestedPickupDate, billedCubicFeet, servicesScheduleOrigin
func (_m *DomesticCratingPricer) Price(appCtx appcontext.AppContext, contractCode string, requestedPickupDate time.Time, billedCubicFeet unit.CubicFeet, servicesScheduleOrigin int) (unit.Cents, services.PricingDisplayParams, error) {
	ret := _m.Called(appCtx, contractCode, requestedPickupDate, billedCubicFeet, servicesScheduleOrigin)

	var r0 unit.Cents
	var r1 services.PricingDisplayParams
	var r2 error
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, string, time.Time, unit.CubicFeet, int) (unit.Cents, services.PricingDisplayParams, error)); ok {
		return rf(appCtx, contractCode, requestedPickupDate, billedCubicFeet, servicesScheduleOrigin)
	}
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, string, time.Time, unit.CubicFeet, int) unit.Cents); ok {
		r0 = rf(appCtx, contractCode, requestedPickupDate, billedCubicFeet, servicesScheduleOrigin)
	} else {
		r0 = ret.Get(0).(unit.Cents)
	}

	if rf, ok := ret.Get(1).(func(appcontext.AppContext, string, time.Time, unit.CubicFeet, int) services.PricingDisplayParams); ok {
		r1 = rf(appCtx, contractCode, requestedPickupDate, billedCubicFeet, servicesScheduleOrigin)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(services.PricingDisplayParams)
		}
	}

	if rf, ok := ret.Get(2).(func(appcontext.AppContext, string, time.Time, unit.CubicFeet, int) error); ok {
		r2 = rf(appCtx, contractCode, requestedPickupDate, billedCubicFeet, servicesScheduleOrigin)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// PriceUsingParams provides a mock function with given fields: appCtx, params
func (_m *DomesticCratingPricer) PriceUsingParams(appCtx appcontext.AppContext, params models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error) {
	ret := _m.Called(appCtx, params)

	var r0 unit.Cents
	var r1 services.PricingDisplayParams
	var r2 error
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, models.PaymentServiceItemParams) (unit.Cents, services.PricingDisplayParams, error)); ok {
		return rf(appCtx, params)
	}
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, models.PaymentServiceItemParams) unit.Cents); ok {
		r0 = rf(appCtx, params)
	} else {
		r0 = ret.Get(0).(unit.Cents)
	}

	if rf, ok := ret.Get(1).(func(appcontext.AppContext, models.PaymentServiceItemParams) services.PricingDisplayParams); ok {
		r1 = rf(appCtx, params)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(services.PricingDisplayParams)
		}
	}

	if rf, ok := ret.Get(2).(func(appcontext.AppContext, models.PaymentServiceItemParams) error); ok {
		r2 = rf(appCtx, params)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

type mockConstructorTestingTNewDomesticCratingPricer interface {
	mock.TestingT
	Cleanup(func())
}

// NewDomesticCratingPricer creates a new instance of DomesticCratingPricer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDomesticCratingPricer(t mockConstructorTestingTNewDomesticCratingPricer) *DomesticCratingPricer {
	mock := &DomesticCratingPricer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
