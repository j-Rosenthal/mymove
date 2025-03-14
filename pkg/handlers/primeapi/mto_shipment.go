package primeapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapi/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
)

// CreateMTOShipmentHandler is the handler to create MTO shipments
type CreateMTOShipmentHandler struct {
	handlers.HandlerConfig
	services.ShipmentCreator
	mtoAvailabilityChecker services.MoveTaskOrderChecker
}

// Handle creates the mto shipment
func (h CreateMTOShipmentHandler) Handle(params mtoshipmentops.CreateMTOShipmentParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			payload := params.Body
			if payload == nil {
				err := apperror.NewBadDataError("the MTO Shipment request body cannot be empty")
				appCtx.Logger().Error(err.Error())
				return mtoshipmentops.NewCreateMTOShipmentBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage,
					err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
			}

			for _, mtoServiceItem := range params.Body.MtoServiceItems() {
				// restrict creation to a list
				if _, ok := CreateableServiceItemMap[mtoServiceItem.ModelType()]; !ok {
					// throw error if modelType() not on the list
					mapKeys := GetMapKeys(CreateableServiceItemMap)
					detailErr := fmt.Sprintf("MTOServiceItem modelType() not allowed: %s ", mtoServiceItem.ModelType())
					verrs := validate.NewErrors()
					verrs.Add("modelType", fmt.Sprintf("allowed modelType() %v", mapKeys))

					appCtx.Logger().Error("primeapi.CreateMTOShipmentHandler error", zap.Error(verrs))
					return mtoshipmentops.NewCreateMTOShipmentUnprocessableEntity().WithPayload(payloads.ValidationError(
						detailErr, h.GetTraceIDFromRequest(params.HTTPRequest), verrs)), verrs
				}
			}

			mtoShipment := payloads.MTOShipmentModelFromCreate(payload)
			mtoShipment.Status = models.MTOShipmentStatusSubmitted
			mtoServiceItemsList, verrs := payloads.MTOServiceItemModelListFromCreate(payload)

			if verrs != nil && verrs.HasAny() {
				appCtx.Logger().Error("Error validating mto service item list: ", zap.Error(verrs))

				return mtoshipmentops.NewCreateMTOShipmentUnprocessableEntity().WithPayload(payloads.ValidationError(
					"The MTO service item list is invalid.", h.GetTraceIDFromRequest(params.HTTPRequest), nil)), verrs
			}

			mtoShipment.MTOServiceItems = mtoServiceItemsList

			moveTaskOrderID := uuid.FromStringOrNil(payload.MoveTaskOrderID.String())
			mtoAvailableToPrime, err := h.mtoAvailabilityChecker.MTOAvailableToPrime(appCtx, moveTaskOrderID)

			if mtoAvailableToPrime {
				mtoShipment, err = h.ShipmentCreator.CreateShipment(appCtx, mtoShipment)
			} else if err == nil {
				appCtx.Logger().Error("primeapi.CreateMTOShipmentHandler error - MTO is not available to Prime")
				return mtoshipmentops.NewCreateMTOShipmentNotFound().WithPayload(payloads.ClientError(
					handlers.NotFoundMessage, fmt.Sprintf("id: %s not found for moveTaskOrder", moveTaskOrderID), h.GetTraceIDFromRequest(params.HTTPRequest))), err
			}

			// Could be the error from MTOAvailableToPrime or CreateMTOShipment:
			if err != nil {
				appCtx.Logger().Error("primeapi.CreateMTOShipmentHandler error", zap.Error(err))
				switch e := err.(type) {
				case apperror.NotFoundError:
					return mtoshipmentops.NewCreateMTOShipmentNotFound().WithPayload(
						payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.InvalidInputError:
					return mtoshipmentops.NewCreateMTOShipmentUnprocessableEntity().WithPayload(
						payloads.ValidationError(err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)), err
				case apperror.QueryError:
					if e.Unwrap() != nil {
						// If you can unwrap, log the internal error (usually a pq error) for better debugging
						appCtx.Logger().Error("primeapi.CreateMTOShipmentHandler query error", zap.Error(e.Unwrap()))
					}
					return mtoshipmentops.NewCreateMTOShipmentInternalServerError().WithPayload(
						payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
				default:
					return mtoshipmentops.NewCreateMTOShipmentInternalServerError().WithPayload(
						payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
				}
			}
			returnPayload := payloads.MTOShipment(mtoShipment)
			return mtoshipmentops.NewCreateMTOShipmentOK().WithPayload(returnPayload), nil
		})
}

// UpdateMTOShipmentHandler is the handler to update MTO shipments
type UpdateMTOShipmentHandler struct {
	handlers.HandlerConfig
	services.ShipmentUpdater
}

// Handle handler that updates a mto shipment
func (h UpdateMTOShipmentHandler) Handle(params mtoshipmentops.UpdateMTOShipmentParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			mtoShipment := payloads.MTOShipmentModelFromUpdate(params.Body, params.MtoShipmentID)

			dbShipment, err := mtoshipment.FindShipment(appCtx, mtoShipment.ID, "DestinationAddress",
				"SecondaryPickupAddress",
				"SecondaryDeliveryAddress",
				"MTOAgents",
				"StorageFacility",
				"PPMShipment")
			if err != nil {
				return mtoshipmentops.NewUpdateMTOShipmentNotFound().WithPayload(
					payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
			}

			// Validate further prime restrictions on model
			mtoShipment.ShipmentType = dbShipment.ShipmentType

			appCtx.Logger().Info("primeapi.UpdateMTOShipmentHandler info", zap.String("pointOfContact", params.Body.PointOfContact))
			mtoShipment, err = h.ShipmentUpdater.UpdateShipment(appCtx, mtoShipment, params.IfMatch)
			if err != nil {
				appCtx.Logger().Error("primeapi.UpdateMTOShipmentHandler error", zap.Error(err))
				switch e := err.(type) {
				case apperror.NotFoundError:
					return mtoshipmentops.NewUpdateMTOShipmentNotFound().WithPayload(
						payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.InvalidInputError:
					payload := payloads.ValidationError(err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)
					return mtoshipmentops.NewUpdateMTOShipmentUnprocessableEntity().WithPayload(payload), err
				case apperror.PreconditionFailedError:
					return mtoshipmentops.NewUpdateMTOShipmentPreconditionFailed().WithPayload(
						payloads.ClientError(handlers.PreconditionErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				default:
					return mtoshipmentops.NewUpdateMTOShipmentInternalServerError().WithPayload(
						payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
				}
			}
			mtoShipmentPayload := payloads.MTOShipment(mtoShipment)
			return mtoshipmentops.NewUpdateMTOShipmentOK().WithPayload(mtoShipmentPayload), nil
		})
}

// DeleteMTOShipmentHandler is the handler to soft delete MTO shipments
type DeleteMTOShipmentHandler struct {
	handlers.HandlerConfig
	services.ShipmentDeleter
}

// Handle handler that updates a mto shipment
func (h DeleteMTOShipmentHandler) Handle(params mtoshipmentops.DeleteMTOShipmentParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			shipmentID := uuid.FromStringOrNil(params.MtoShipmentID.String())
			_, err := h.DeleteShipment(appCtx, shipmentID)
			if err != nil {
				appCtx.Logger().Error("primeapi.DeleteMTOShipmentHandler", zap.Error(err))

				switch err.(type) {
				case apperror.NotFoundError:
					return mtoshipmentops.NewDeleteMTOShipmentNotFound().WithPayload(
						payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.ConflictError:
					return mtoshipmentops.NewDeleteMTOShipmentConflict(), err
				case apperror.ForbiddenError:
					return mtoshipmentops.NewDeleteMTOShipmentForbidden().WithPayload(
						payloads.ClientError(handlers.ForbiddenErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.UnprocessableEntityError:
					return mtoshipmentops.NewDeleteMTOShipmentUnprocessableEntity(), err
				default:
					return mtoshipmentops.NewDeleteMTOShipmentInternalServerError().WithPayload(
						payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
				}
			}

			return mtoshipmentops.NewDeleteMTOShipmentNoContent(), nil
		})
}

// UpdateMTOShipmentStatusHandler is the handler to update MTO Shipments' status
type UpdateMTOShipmentStatusHandler struct {
	handlers.HandlerConfig
	checker services.MTOShipmentUpdater
	updater services.MTOShipmentStatusUpdater
}

// Handle handler that updates a mto shipment's status
func (h UpdateMTOShipmentStatusHandler) Handle(params mtoshipmentops.UpdateMTOShipmentStatusParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			shipmentID := uuid.FromStringOrNil(params.MtoShipmentID.String())

			availableToPrime, err := h.checker.MTOShipmentsMTOAvailableToPrime(appCtx, shipmentID)
			if err != nil {
				appCtx.Logger().Error("primeapi.UpdateMTOShipmentHandler error - MTO is not available to prime", zap.Error(err))
				switch e := err.(type) {
				case apperror.NotFoundError:
					return mtoshipmentops.NewUpdateMTOShipmentStatusNotFound().WithPayload(
						payloads.ClientError(handlers.NotFoundMessage, e.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				default:
					return mtoshipmentops.NewUpdateMTOShipmentStatusInternalServerError().WithPayload(
						payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
				}
			}
			if !availableToPrime {
				return mtoshipmentops.NewUpdateMTOShipmentStatusInternalServerError().WithPayload(
					payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest))), err
			}

			status := models.MTOShipmentStatus(params.Body.Status)
			eTag := params.IfMatch

			shipment, err := h.updater.UpdateMTOShipmentStatus(appCtx, shipmentID, status, nil, eTag)
			if err != nil {
				appCtx.Logger().Error("UpdateMTOShipmentStatusStatus error: ", zap.Error(err))

				switch e := err.(type) {
				case apperror.NotFoundError:
					return mtoshipmentops.NewUpdateMTOShipmentStatusNotFound().WithPayload(
						payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case apperror.InvalidInputError:
					return mtoshipmentops.NewUpdateMTOShipmentStatusUnprocessableEntity().WithPayload(
						payloads.ValidationError("The input provided did not pass validation.", h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)), err
				case apperror.PreconditionFailedError:
					return mtoshipmentops.NewUpdateMTOShipmentStatusPreconditionFailed().WithPayload(
						payloads.ClientError(handlers.PreconditionErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				case mtoshipment.ConflictStatusError:
					return mtoshipmentops.NewUpdateMTOShipmentStatusConflict().WithPayload(
						payloads.ClientError(handlers.ConflictErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				default:
					return mtoshipmentops.NewUpdateMTOShipmentStatusInternalServerError().WithPayload(
						payloads.InternalServerError(handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest))), err
				}
			}

			return mtoshipmentops.NewUpdateMTOShipmentStatusOK().WithPayload(payloads.MTOShipment(shipment)), nil
		})
}
