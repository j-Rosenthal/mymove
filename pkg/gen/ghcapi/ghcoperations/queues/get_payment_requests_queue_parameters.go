// Code generated by go-swagger; DO NOT EDIT.

package queues

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// NewGetPaymentRequestsQueueParams creates a new GetPaymentRequestsQueueParams object
//
// There are no default values defined in the spec.
func NewGetPaymentRequestsQueueParams() GetPaymentRequestsQueueParams {

	return GetPaymentRequestsQueueParams{}
}

// GetPaymentRequestsQueueParams contains all the bound params for the get payment requests queue operation
// typically these are obtained from a http.Request
//
// swagger:parameters getPaymentRequestsQueue
type GetPaymentRequestsQueueParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  In: query
	*/
	Branch *string
	/*
	  In: query
	*/
	DestinationDutyLocation *string
	/*
	  In: query
	*/
	DodID *string
	/*
	  In: query
	*/
	LastName *string
	/*
	  In: query
	*/
	Locator *string
	/*direction of sort order if applied
	  In: query
	*/
	Order *string
	/*
	  In: query
	*/
	OriginDutyLocation *string
	/*requested page of results
	  In: query
	*/
	Page *int64
	/*number of records to include per page
	  In: query
	*/
	PerPage *int64
	/*field that results should be sorted by
	  In: query
	*/
	Sort *string
	/*Filtering for the status.
	  Unique: true
	  In: query
	*/
	Status []string
	/*Start of the submitted at date in the user's local time zone converted to UTC
	  In: query
	*/
	SubmittedAt *strfmt.DateTime
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetPaymentRequestsQueueParams() beforehand.
func (o *GetPaymentRequestsQueueParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qBranch, qhkBranch, _ := qs.GetOK("branch")
	if err := o.bindBranch(qBranch, qhkBranch, route.Formats); err != nil {
		res = append(res, err)
	}

	qDestinationDutyLocation, qhkDestinationDutyLocation, _ := qs.GetOK("destinationDutyLocation")
	if err := o.bindDestinationDutyLocation(qDestinationDutyLocation, qhkDestinationDutyLocation, route.Formats); err != nil {
		res = append(res, err)
	}

	qDodID, qhkDodID, _ := qs.GetOK("dodID")
	if err := o.bindDodID(qDodID, qhkDodID, route.Formats); err != nil {
		res = append(res, err)
	}

	qLastName, qhkLastName, _ := qs.GetOK("lastName")
	if err := o.bindLastName(qLastName, qhkLastName, route.Formats); err != nil {
		res = append(res, err)
	}

	qLocator, qhkLocator, _ := qs.GetOK("locator")
	if err := o.bindLocator(qLocator, qhkLocator, route.Formats); err != nil {
		res = append(res, err)
	}

	qOrder, qhkOrder, _ := qs.GetOK("order")
	if err := o.bindOrder(qOrder, qhkOrder, route.Formats); err != nil {
		res = append(res, err)
	}

	qOriginDutyLocation, qhkOriginDutyLocation, _ := qs.GetOK("originDutyLocation")
	if err := o.bindOriginDutyLocation(qOriginDutyLocation, qhkOriginDutyLocation, route.Formats); err != nil {
		res = append(res, err)
	}

	qPage, qhkPage, _ := qs.GetOK("page")
	if err := o.bindPage(qPage, qhkPage, route.Formats); err != nil {
		res = append(res, err)
	}

	qPerPage, qhkPerPage, _ := qs.GetOK("perPage")
	if err := o.bindPerPage(qPerPage, qhkPerPage, route.Formats); err != nil {
		res = append(res, err)
	}

	qSort, qhkSort, _ := qs.GetOK("sort")
	if err := o.bindSort(qSort, qhkSort, route.Formats); err != nil {
		res = append(res, err)
	}

	qStatus, qhkStatus, _ := qs.GetOK("status")
	if err := o.bindStatus(qStatus, qhkStatus, route.Formats); err != nil {
		res = append(res, err)
	}

	qSubmittedAt, qhkSubmittedAt, _ := qs.GetOK("submittedAt")
	if err := o.bindSubmittedAt(qSubmittedAt, qhkSubmittedAt, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindBranch binds and validates parameter Branch from query.
func (o *GetPaymentRequestsQueueParams) bindBranch(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}
	o.Branch = &raw

	return nil
}

// bindDestinationDutyLocation binds and validates parameter DestinationDutyLocation from query.
func (o *GetPaymentRequestsQueueParams) bindDestinationDutyLocation(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}
	o.DestinationDutyLocation = &raw

	return nil
}

// bindDodID binds and validates parameter DodID from query.
func (o *GetPaymentRequestsQueueParams) bindDodID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}
	o.DodID = &raw

	return nil
}

// bindLastName binds and validates parameter LastName from query.
func (o *GetPaymentRequestsQueueParams) bindLastName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}
	o.LastName = &raw

	return nil
}

// bindLocator binds and validates parameter Locator from query.
func (o *GetPaymentRequestsQueueParams) bindLocator(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}
	o.Locator = &raw

	return nil
}

// bindOrder binds and validates parameter Order from query.
func (o *GetPaymentRequestsQueueParams) bindOrder(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}
	o.Order = &raw

	if err := o.validateOrder(formats); err != nil {
		return err
	}

	return nil
}

// validateOrder carries on validations for parameter Order
func (o *GetPaymentRequestsQueueParams) validateOrder(formats strfmt.Registry) error {

	if err := validate.EnumCase("order", "query", *o.Order, []interface{}{"asc", "desc"}, true); err != nil {
		return err
	}

	return nil
}

// bindOriginDutyLocation binds and validates parameter OriginDutyLocation from query.
func (o *GetPaymentRequestsQueueParams) bindOriginDutyLocation(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}
	o.OriginDutyLocation = &raw

	return nil
}

// bindPage binds and validates parameter Page from query.
func (o *GetPaymentRequestsQueueParams) bindPage(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("page", "query", "int64", raw)
	}
	o.Page = &value

	return nil
}

// bindPerPage binds and validates parameter PerPage from query.
func (o *GetPaymentRequestsQueueParams) bindPerPage(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("perPage", "query", "int64", raw)
	}
	o.PerPage = &value

	return nil
}

// bindSort binds and validates parameter Sort from query.
func (o *GetPaymentRequestsQueueParams) bindSort(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}
	o.Sort = &raw

	if err := o.validateSort(formats); err != nil {
		return err
	}

	return nil
}

// validateSort carries on validations for parameter Sort
func (o *GetPaymentRequestsQueueParams) validateSort(formats strfmt.Registry) error {

	if err := validate.EnumCase("sort", "query", *o.Sort, []interface{}{"lastName", "locator", "submittedAt", "branch", "status", "dodID", "age", "originDutyLocation"}, true); err != nil {
		return err
	}

	return nil
}

// bindStatus binds and validates array parameter Status from query.
//
// Arrays are parsed according to CollectionFormat: "" (defaults to "csv" when empty).
func (o *GetPaymentRequestsQueueParams) bindStatus(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var qvStatus string
	if len(rawData) > 0 {
		qvStatus = rawData[len(rawData)-1]
	}

	// CollectionFormat:
	statusIC := swag.SplitByFormat(qvStatus, "")
	if len(statusIC) == 0 {
		return nil
	}

	var statusIR []string
	for i, statusIV := range statusIC {
		statusI := statusIV

		if err := validate.EnumCase(fmt.Sprintf("%s.%v", "status", i), "query", statusI, []interface{}{"Payment requested", "Reviewed", "Rejected", "Paid"}, true); err != nil {
			return err
		}

		statusIR = append(statusIR, statusI)
	}

	o.Status = statusIR
	if err := o.validateStatus(formats); err != nil {
		return err
	}

	return nil
}

// validateStatus carries on validations for parameter Status
func (o *GetPaymentRequestsQueueParams) validateStatus(formats strfmt.Registry) error {

	// uniqueItems: true
	if err := validate.UniqueItems("status", "query", o.Status); err != nil {
		return err
	}
	return nil
}

// bindSubmittedAt binds and validates parameter SubmittedAt from query.
func (o *GetPaymentRequestsQueueParams) bindSubmittedAt(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}

	// Format: date-time
	value, err := formats.Parse("date-time", raw)
	if err != nil {
		return errors.InvalidType("submittedAt", "query", "strfmt.DateTime", raw)
	}
	o.SubmittedAt = (value.(*strfmt.DateTime))

	if err := o.validateSubmittedAt(formats); err != nil {
		return err
	}

	return nil
}

// validateSubmittedAt carries on validations for parameter SubmittedAt
func (o *GetPaymentRequestsQueueParams) validateSubmittedAt(formats strfmt.Registry) error {

	if err := validate.FormatOf("submittedAt", "query", "date-time", o.SubmittedAt.String(), formats); err != nil {
		return err
	}
	return nil
}
