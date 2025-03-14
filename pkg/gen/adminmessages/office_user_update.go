// Code generated by go-swagger; DO NOT EDIT.

package adminmessages

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// OfficeUserUpdate office user update
//
// swagger:model OfficeUserUpdate
type OfficeUserUpdate struct {

	// active
	Active *bool `json:"active,omitempty"`

	// First Name
	FirstName *string `json:"firstName,omitempty"`

	// Last Name
	LastName *string `json:"lastName,omitempty"`

	// Middle Initials
	// Example: Q.
	MiddleInitials *string `json:"middleInitials,omitempty"`

	// roles
	Roles []*OfficeUserRole `json:"roles"`

	// telephone
	// Example: 212-555-5555
	// Pattern: ^[2-9]\d{2}-\d{3}-\d{4}$
	Telephone *string `json:"telephone,omitempty"`

	// transportation office Id
	// Example: c56a4180-65aa-42ec-a945-5fd21dec0538
	// Format: uuid
	TransportationOfficeID strfmt.UUID `json:"transportationOfficeId,omitempty"`
}

// Validate validates this office user update
func (m *OfficeUserUpdate) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateRoles(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTelephone(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTransportationOfficeID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *OfficeUserUpdate) validateRoles(formats strfmt.Registry) error {
	if swag.IsZero(m.Roles) { // not required
		return nil
	}

	for i := 0; i < len(m.Roles); i++ {
		if swag.IsZero(m.Roles[i]) { // not required
			continue
		}

		if m.Roles[i] != nil {
			if err := m.Roles[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("roles" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("roles" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *OfficeUserUpdate) validateTelephone(formats strfmt.Registry) error {
	if swag.IsZero(m.Telephone) { // not required
		return nil
	}

	if err := validate.Pattern("telephone", "body", *m.Telephone, `^[2-9]\d{2}-\d{3}-\d{4}$`); err != nil {
		return err
	}

	return nil
}

func (m *OfficeUserUpdate) validateTransportationOfficeID(formats strfmt.Registry) error {
	if swag.IsZero(m.TransportationOfficeID) { // not required
		return nil
	}

	if err := validate.FormatOf("transportationOfficeId", "body", "uuid", m.TransportationOfficeID.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this office user update based on the context it is used
func (m *OfficeUserUpdate) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateRoles(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *OfficeUserUpdate) contextValidateRoles(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Roles); i++ {

		if m.Roles[i] != nil {

			if swag.IsZero(m.Roles[i]) { // not required
				return nil
			}

			if err := m.Roles[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("roles" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("roles" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *OfficeUserUpdate) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *OfficeUserUpdate) UnmarshalBinary(b []byte) error {
	var res OfficeUserUpdate
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
