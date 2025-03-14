// Code generated by go-swagger; DO NOT EDIT.

package webhook_subscriptions

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// UpdateWebhookSubscriptionHandlerFunc turns a function with the right signature into a update webhook subscription handler
type UpdateWebhookSubscriptionHandlerFunc func(UpdateWebhookSubscriptionParams) middleware.Responder

// Handle executing the request and returning a response
func (fn UpdateWebhookSubscriptionHandlerFunc) Handle(params UpdateWebhookSubscriptionParams) middleware.Responder {
	return fn(params)
}

// UpdateWebhookSubscriptionHandler interface for that can handle valid update webhook subscription params
type UpdateWebhookSubscriptionHandler interface {
	Handle(UpdateWebhookSubscriptionParams) middleware.Responder
}

// NewUpdateWebhookSubscription creates a new http.Handler for the update webhook subscription operation
func NewUpdateWebhookSubscription(ctx *middleware.Context, handler UpdateWebhookSubscriptionHandler) *UpdateWebhookSubscription {
	return &UpdateWebhookSubscription{Context: ctx, Handler: handler}
}

/*
	UpdateWebhookSubscription swagger:route PATCH /webhook-subscriptions/{webhookSubscriptionId} Webhook subscriptions updateWebhookSubscription

# Update a Webhook Subscription

This endpoint updates a single Webhook Subscription by ID. Do not use this
endpoint directly as it is meant to be used with the Admin UI exclusively.
*/
type UpdateWebhookSubscription struct {
	Context *middleware.Context
	Handler UpdateWebhookSubscriptionHandler
}

func (o *UpdateWebhookSubscription) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewUpdateWebhookSubscriptionParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
