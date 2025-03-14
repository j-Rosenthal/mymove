// Code generated by go-swagger; DO NOT EDIT.

package webhook_subscriptions

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// NewGetWebhookSubscriptionParams creates a new GetWebhookSubscriptionParams object
//
// There are no default values defined in the spec.
func NewGetWebhookSubscriptionParams() GetWebhookSubscriptionParams {

	return GetWebhookSubscriptionParams{}
}

// GetWebhookSubscriptionParams contains all the bound params for the get webhook subscription operation
// typically these are obtained from a http.Request
//
// swagger:parameters getWebhookSubscription
type GetWebhookSubscriptionParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: path
	*/
	WebhookSubscriptionID strfmt.UUID
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetWebhookSubscriptionParams() beforehand.
func (o *GetWebhookSubscriptionParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rWebhookSubscriptionID, rhkWebhookSubscriptionID, _ := route.Params.GetOK("webhookSubscriptionId")
	if err := o.bindWebhookSubscriptionID(rWebhookSubscriptionID, rhkWebhookSubscriptionID, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindWebhookSubscriptionID binds and validates parameter WebhookSubscriptionID from path.
func (o *GetWebhookSubscriptionParams) bindWebhookSubscriptionID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	// Format: uuid
	value, err := formats.Parse("uuid", raw)
	if err != nil {
		return errors.InvalidType("webhookSubscriptionId", "path", "strfmt.UUID", raw)
	}
	o.WebhookSubscriptionID = *(value.(*strfmt.UUID))

	if err := o.validateWebhookSubscriptionID(formats); err != nil {
		return err
	}

	return nil
}

// validateWebhookSubscriptionID carries on validations for parameter WebhookSubscriptionID
func (o *GetWebhookSubscriptionParams) validateWebhookSubscriptionID(formats strfmt.Registry) error {

	if err := validate.FormatOf("webhookSubscriptionId", "path", "uuid", o.WebhookSubscriptionID.String(), formats); err != nil {
		return err
	}
	return nil
}
