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

// MoveDocumentType Document type
// Example: EXPENSE
//
// swagger:model MoveDocumentType
type MoveDocumentType string

func NewMoveDocumentType(value MoveDocumentType) *MoveDocumentType {
	return &value
}

// Pointer returns a pointer to a freshly-allocated MoveDocumentType.
func (m MoveDocumentType) Pointer() *MoveDocumentType {
	return &m
}

const (

	// MoveDocumentTypeOTHER captures enum value "OTHER"
	MoveDocumentTypeOTHER MoveDocumentType = "OTHER"

	// MoveDocumentTypeWEIGHTTICKET captures enum value "WEIGHT_TICKET"
	MoveDocumentTypeWEIGHTTICKET MoveDocumentType = "WEIGHT_TICKET"

	// MoveDocumentTypeSTORAGEEXPENSE captures enum value "STORAGE_EXPENSE"
	MoveDocumentTypeSTORAGEEXPENSE MoveDocumentType = "STORAGE_EXPENSE"

	// MoveDocumentTypeSHIPMENTSUMMARY captures enum value "SHIPMENT_SUMMARY"
	MoveDocumentTypeSHIPMENTSUMMARY MoveDocumentType = "SHIPMENT_SUMMARY"

	// MoveDocumentTypeEXPENSE captures enum value "EXPENSE"
	MoveDocumentTypeEXPENSE MoveDocumentType = "EXPENSE"

	// MoveDocumentTypeWEIGHTTICKETSET captures enum value "WEIGHT_TICKET_SET"
	MoveDocumentTypeWEIGHTTICKETSET MoveDocumentType = "WEIGHT_TICKET_SET"
)

// for schema
var moveDocumentTypeEnum []interface{}

func init() {
	var res []MoveDocumentType
	if err := json.Unmarshal([]byte(`["OTHER","WEIGHT_TICKET","STORAGE_EXPENSE","SHIPMENT_SUMMARY","EXPENSE","WEIGHT_TICKET_SET"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		moveDocumentTypeEnum = append(moveDocumentTypeEnum, v)
	}
}

func (m MoveDocumentType) validateMoveDocumentTypeEnum(path, location string, value MoveDocumentType) error {
	if err := validate.EnumCase(path, location, value, moveDocumentTypeEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this move document type
func (m MoveDocumentType) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateMoveDocumentTypeEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validates this move document type based on context it is used
func (m MoveDocumentType) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}
