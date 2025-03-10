// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	appcontext "github.com/transcom/mymove/pkg/appcontext"

	models "github.com/transcom/mymove/pkg/models"

	uuid "github.com/gofrs/uuid"
)

// SITExtensionDenier is an autogenerated mock type for the SITExtensionDenier type
type SITExtensionDenier struct {
	mock.Mock
}

// DenySITExtension provides a mock function with given fields: appCtx, shipmentID, sitExtensionID, officeRemarks, eTag
func (_m *SITExtensionDenier) DenySITExtension(appCtx appcontext.AppContext, shipmentID uuid.UUID, sitExtensionID uuid.UUID, officeRemarks *string, eTag string) (*models.MTOShipment, error) {
	ret := _m.Called(appCtx, shipmentID, sitExtensionID, officeRemarks, eTag)

	var r0 *models.MTOShipment
	var r1 error
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, uuid.UUID, uuid.UUID, *string, string) (*models.MTOShipment, error)); ok {
		return rf(appCtx, shipmentID, sitExtensionID, officeRemarks, eTag)
	}
	if rf, ok := ret.Get(0).(func(appcontext.AppContext, uuid.UUID, uuid.UUID, *string, string) *models.MTOShipment); ok {
		r0 = rf(appCtx, shipmentID, sitExtensionID, officeRemarks, eTag)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.MTOShipment)
		}
	}

	if rf, ok := ret.Get(1).(func(appcontext.AppContext, uuid.UUID, uuid.UUID, *string, string) error); ok {
		r1 = rf(appCtx, shipmentID, sitExtensionID, officeRemarks, eTag)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewSITExtensionDenier interface {
	mock.TestingT
	Cleanup(func())
}

// NewSITExtensionDenier creates a new instance of SITExtensionDenier. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSITExtensionDenier(t mockConstructorTestingTNewSITExtensionDenier) *SITExtensionDenier {
	mock := &SITExtensionDenier{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
