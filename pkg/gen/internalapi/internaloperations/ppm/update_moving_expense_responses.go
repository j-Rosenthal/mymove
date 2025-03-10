// Code generated by go-swagger; DO NOT EDIT.

package ppm

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// UpdateMovingExpenseOKCode is the HTTP code returned for type UpdateMovingExpenseOK
const UpdateMovingExpenseOKCode int = 200

/*
UpdateMovingExpenseOK returns an updated moving expense object

swagger:response updateMovingExpenseOK
*/
type UpdateMovingExpenseOK struct {

	/*
	  In: Body
	*/
	Payload *internalmessages.MovingExpense `json:"body,omitempty"`
}

// NewUpdateMovingExpenseOK creates UpdateMovingExpenseOK with default headers values
func NewUpdateMovingExpenseOK() *UpdateMovingExpenseOK {

	return &UpdateMovingExpenseOK{}
}

// WithPayload adds the payload to the update moving expense o k response
func (o *UpdateMovingExpenseOK) WithPayload(payload *internalmessages.MovingExpense) *UpdateMovingExpenseOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update moving expense o k response
func (o *UpdateMovingExpenseOK) SetPayload(payload *internalmessages.MovingExpense) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateMovingExpenseOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateMovingExpenseBadRequestCode is the HTTP code returned for type UpdateMovingExpenseBadRequest
const UpdateMovingExpenseBadRequestCode int = 400

/*
UpdateMovingExpenseBadRequest The request payload is invalid.

swagger:response updateMovingExpenseBadRequest
*/
type UpdateMovingExpenseBadRequest struct {

	/*
	  In: Body
	*/
	Payload *internalmessages.ClientError `json:"body,omitempty"`
}

// NewUpdateMovingExpenseBadRequest creates UpdateMovingExpenseBadRequest with default headers values
func NewUpdateMovingExpenseBadRequest() *UpdateMovingExpenseBadRequest {

	return &UpdateMovingExpenseBadRequest{}
}

// WithPayload adds the payload to the update moving expense bad request response
func (o *UpdateMovingExpenseBadRequest) WithPayload(payload *internalmessages.ClientError) *UpdateMovingExpenseBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update moving expense bad request response
func (o *UpdateMovingExpenseBadRequest) SetPayload(payload *internalmessages.ClientError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateMovingExpenseBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateMovingExpenseUnauthorizedCode is the HTTP code returned for type UpdateMovingExpenseUnauthorized
const UpdateMovingExpenseUnauthorizedCode int = 401

/*
UpdateMovingExpenseUnauthorized The request was denied.

swagger:response updateMovingExpenseUnauthorized
*/
type UpdateMovingExpenseUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *internalmessages.ClientError `json:"body,omitempty"`
}

// NewUpdateMovingExpenseUnauthorized creates UpdateMovingExpenseUnauthorized with default headers values
func NewUpdateMovingExpenseUnauthorized() *UpdateMovingExpenseUnauthorized {

	return &UpdateMovingExpenseUnauthorized{}
}

// WithPayload adds the payload to the update moving expense unauthorized response
func (o *UpdateMovingExpenseUnauthorized) WithPayload(payload *internalmessages.ClientError) *UpdateMovingExpenseUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update moving expense unauthorized response
func (o *UpdateMovingExpenseUnauthorized) SetPayload(payload *internalmessages.ClientError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateMovingExpenseUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateMovingExpenseForbiddenCode is the HTTP code returned for type UpdateMovingExpenseForbidden
const UpdateMovingExpenseForbiddenCode int = 403

/*
UpdateMovingExpenseForbidden The request was denied.

swagger:response updateMovingExpenseForbidden
*/
type UpdateMovingExpenseForbidden struct {

	/*
	  In: Body
	*/
	Payload *internalmessages.ClientError `json:"body,omitempty"`
}

// NewUpdateMovingExpenseForbidden creates UpdateMovingExpenseForbidden with default headers values
func NewUpdateMovingExpenseForbidden() *UpdateMovingExpenseForbidden {

	return &UpdateMovingExpenseForbidden{}
}

// WithPayload adds the payload to the update moving expense forbidden response
func (o *UpdateMovingExpenseForbidden) WithPayload(payload *internalmessages.ClientError) *UpdateMovingExpenseForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update moving expense forbidden response
func (o *UpdateMovingExpenseForbidden) SetPayload(payload *internalmessages.ClientError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateMovingExpenseForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateMovingExpenseNotFoundCode is the HTTP code returned for type UpdateMovingExpenseNotFound
const UpdateMovingExpenseNotFoundCode int = 404

/*
UpdateMovingExpenseNotFound The requested resource wasn't found.

swagger:response updateMovingExpenseNotFound
*/
type UpdateMovingExpenseNotFound struct {

	/*
	  In: Body
	*/
	Payload *internalmessages.ClientError `json:"body,omitempty"`
}

// NewUpdateMovingExpenseNotFound creates UpdateMovingExpenseNotFound with default headers values
func NewUpdateMovingExpenseNotFound() *UpdateMovingExpenseNotFound {

	return &UpdateMovingExpenseNotFound{}
}

// WithPayload adds the payload to the update moving expense not found response
func (o *UpdateMovingExpenseNotFound) WithPayload(payload *internalmessages.ClientError) *UpdateMovingExpenseNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update moving expense not found response
func (o *UpdateMovingExpenseNotFound) SetPayload(payload *internalmessages.ClientError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateMovingExpenseNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateMovingExpensePreconditionFailedCode is the HTTP code returned for type UpdateMovingExpensePreconditionFailed
const UpdateMovingExpensePreconditionFailedCode int = 412

/*
UpdateMovingExpensePreconditionFailed Precondition failed, likely due to a stale eTag (If-Match). Fetch the request again to get the updated eTag value.

swagger:response updateMovingExpensePreconditionFailed
*/
type UpdateMovingExpensePreconditionFailed struct {

	/*
	  In: Body
	*/
	Payload *internalmessages.ClientError `json:"body,omitempty"`
}

// NewUpdateMovingExpensePreconditionFailed creates UpdateMovingExpensePreconditionFailed with default headers values
func NewUpdateMovingExpensePreconditionFailed() *UpdateMovingExpensePreconditionFailed {

	return &UpdateMovingExpensePreconditionFailed{}
}

// WithPayload adds the payload to the update moving expense precondition failed response
func (o *UpdateMovingExpensePreconditionFailed) WithPayload(payload *internalmessages.ClientError) *UpdateMovingExpensePreconditionFailed {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update moving expense precondition failed response
func (o *UpdateMovingExpensePreconditionFailed) SetPayload(payload *internalmessages.ClientError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateMovingExpensePreconditionFailed) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(412)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateMovingExpenseUnprocessableEntityCode is the HTTP code returned for type UpdateMovingExpenseUnprocessableEntity
const UpdateMovingExpenseUnprocessableEntityCode int = 422

/*
UpdateMovingExpenseUnprocessableEntity The payload was unprocessable.

swagger:response updateMovingExpenseUnprocessableEntity
*/
type UpdateMovingExpenseUnprocessableEntity struct {

	/*
	  In: Body
	*/
	Payload *internalmessages.ValidationError `json:"body,omitempty"`
}

// NewUpdateMovingExpenseUnprocessableEntity creates UpdateMovingExpenseUnprocessableEntity with default headers values
func NewUpdateMovingExpenseUnprocessableEntity() *UpdateMovingExpenseUnprocessableEntity {

	return &UpdateMovingExpenseUnprocessableEntity{}
}

// WithPayload adds the payload to the update moving expense unprocessable entity response
func (o *UpdateMovingExpenseUnprocessableEntity) WithPayload(payload *internalmessages.ValidationError) *UpdateMovingExpenseUnprocessableEntity {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update moving expense unprocessable entity response
func (o *UpdateMovingExpenseUnprocessableEntity) SetPayload(payload *internalmessages.ValidationError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateMovingExpenseUnprocessableEntity) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(422)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateMovingExpenseInternalServerErrorCode is the HTTP code returned for type UpdateMovingExpenseInternalServerError
const UpdateMovingExpenseInternalServerErrorCode int = 500

/*
UpdateMovingExpenseInternalServerError A server error occurred.

swagger:response updateMovingExpenseInternalServerError
*/
type UpdateMovingExpenseInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *internalmessages.Error `json:"body,omitempty"`
}

// NewUpdateMovingExpenseInternalServerError creates UpdateMovingExpenseInternalServerError with default headers values
func NewUpdateMovingExpenseInternalServerError() *UpdateMovingExpenseInternalServerError {

	return &UpdateMovingExpenseInternalServerError{}
}

// WithPayload adds the payload to the update moving expense internal server error response
func (o *UpdateMovingExpenseInternalServerError) WithPayload(payload *internalmessages.Error) *UpdateMovingExpenseInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update moving expense internal server error response
func (o *UpdateMovingExpenseInternalServerError) SetPayload(payload *internalmessages.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateMovingExpenseInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
