// Code generated by go-swagger; DO NOT EDIT.

package internalmessages

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// PPMEstimateRange p p m estimate range
//
// swagger:model PPMEstimateRange
type PPMEstimateRange struct {

	// High estimate
	// Required: true
	RangeMax *int64 `json:"range_max"`

	// Low estimate
	// Required: true
	RangeMin *int64 `json:"range_min"`
}

// Validate validates this p p m estimate range
func (m *PPMEstimateRange) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateRangeMax(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRangeMin(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PPMEstimateRange) validateRangeMax(formats strfmt.Registry) error {

	if err := validate.Required("range_max", "body", m.RangeMax); err != nil {
		return err
	}

	return nil
}

func (m *PPMEstimateRange) validateRangeMin(formats strfmt.Registry) error {

	if err := validate.Required("range_min", "body", m.RangeMin); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this p p m estimate range based on context it is used
func (m *PPMEstimateRange) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *PPMEstimateRange) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PPMEstimateRange) UnmarshalBinary(b []byte) error {
	var res PPMEstimateRange
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
