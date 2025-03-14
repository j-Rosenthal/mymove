// Code generated by go-swagger; DO NOT EDIT.

package internalmessages

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// SignedCertificationTypeCreate signed certification type create
//
// swagger:model SignedCertificationTypeCreate
type SignedCertificationTypeCreate string

func NewSignedCertificationTypeCreate(value SignedCertificationTypeCreate) *SignedCertificationTypeCreate {
	return &value
}

// Pointer returns a pointer to a freshly-allocated SignedCertificationTypeCreate.
func (m SignedCertificationTypeCreate) Pointer() *SignedCertificationTypeCreate {
	return &m
}

const (

	// SignedCertificationTypeCreatePPMPAYMENT captures enum value "PPM_PAYMENT"
	SignedCertificationTypeCreatePPMPAYMENT SignedCertificationTypeCreate = "PPM_PAYMENT"

	// SignedCertificationTypeCreateSHIPMENT captures enum value "SHIPMENT"
	SignedCertificationTypeCreateSHIPMENT SignedCertificationTypeCreate = "SHIPMENT"
)

// for schema
var signedCertificationTypeCreateEnum []interface{}

func init() {
	var res []SignedCertificationTypeCreate
	if err := json.Unmarshal([]byte(`["PPM_PAYMENT","SHIPMENT"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		signedCertificationTypeCreateEnum = append(signedCertificationTypeCreateEnum, v)
	}
}

func (m SignedCertificationTypeCreate) validateSignedCertificationTypeCreateEnum(path, location string, value SignedCertificationTypeCreate) error {
	if err := validate.EnumCase(path, location, value, signedCertificationTypeCreateEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this signed certification type create
func (m SignedCertificationTypeCreate) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateSignedCertificationTypeCreateEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validates this signed certification type create based on context it is used
func (m SignedCertificationTypeCreate) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}
