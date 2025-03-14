package ghcapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/pop/v6"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/queues"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi/internal/payloads"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
)

// GetMovesQueueHandler returns the moves for the TOO queue user via GET /queues/moves
type GetMovesQueueHandler struct {
	handlers.HandlerConfig
	services.OrderFetcher
}

// FilterOption defines the type for the functional arguments used for private functions in OrderFetcher
type FilterOption func(*pop.Query)

// Handle returns the paginated list of moves for the TOO user
func (h GetMovesQueueHandler) Handle(params queues.GetMovesQueueParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			if !appCtx.Session().IsOfficeUser() ||
				!appCtx.Session().Roles.HasRole(roles.RoleTypeTOO) {
				forbiddenErr := apperror.NewForbiddenError(
					"user is not authenticated with TOO office role",
				)
				appCtx.Logger().Error(forbiddenErr.Error())
				return queues.NewGetMovesQueueForbidden(), forbiddenErr
			}

			ListOrderParams := services.ListOrderParams{
				Branch:                  params.Branch,
				Locator:                 params.Locator,
				DodID:                   params.DodID,
				LastName:                params.LastName,
				DestinationDutyLocation: params.DestinationDutyLocation,
				OriginDutyLocation:      params.OriginDutyLocation,
				AppearedInTOOAt:         handlers.FmtDateTimePtrToPopPtr(params.AppearedInTooAt),
				RequestedMoveDate:       params.RequestedMoveDate,
				Status:                  params.Status,
				Page:                    params.Page,
				PerPage:                 params.PerPage,
				Sort:                    params.Sort,
				Order:                   params.Order,
			}

			// Let's set default values for page and perPage if we don't get arguments for them. We'll use 1 for page and 20
			// for perPage.
			if params.Page == nil {
				ListOrderParams.Page = models.Int64Pointer(1)
			}
			// Same for perPage
			if params.PerPage == nil {
				ListOrderParams.PerPage = models.Int64Pointer(20)
			}

			moves, count, err := h.OrderFetcher.ListOrders(
				appCtx,
				appCtx.Session().OfficeUserID,
				&ListOrderParams,
			)
			if err != nil {
				appCtx.Logger().
					Error("error fetching list of moves for office user", zap.Error(err))
				return queues.NewGetMovesQueueInternalServerError(), err
			}

			queueMoves := payloads.QueueMoves(moves)

			result := &ghcmessages.QueueMovesResult{
				Page:       *ListOrderParams.Page,
				PerPage:    *ListOrderParams.PerPage,
				TotalCount: int64(count),
				QueueMoves: *queueMoves,
			}

			return queues.NewGetMovesQueueOK().WithPayload(result), nil
		})
}

// GetPaymentRequestsQueueHandler returns the payment requests for the TIO queue user via GET /queues/payment-requests
type GetPaymentRequestsQueueHandler struct {
	handlers.HandlerConfig
	services.PaymentRequestListFetcher
}

// Handle returns the paginated list of payment requests for the TIO user
func (h GetPaymentRequestsQueueHandler) Handle(
	params queues.GetPaymentRequestsQueueParams,
) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			if !appCtx.Session().Roles.HasRole(roles.RoleTypeTIO) {
				forbiddenErr := apperror.NewForbiddenError(
					"user is not authenticated with TIO office role",
				)
				appCtx.Logger().Error(forbiddenErr.Error())
				return queues.NewGetPaymentRequestsQueueForbidden(), forbiddenErr
			}

			listPaymentRequestParams := services.FetchPaymentRequestListParams{
				Branch:                  params.Branch,
				Locator:                 params.Locator,
				DodID:                   params.DodID,
				LastName:                params.LastName,
				DestinationDutyLocation: params.DestinationDutyLocation,
				Status:                  params.Status,
				Page:                    params.Page,
				PerPage:                 params.PerPage,
				SubmittedAt:             handlers.FmtDateTimePtrToPopPtr(params.SubmittedAt),
				Sort:                    params.Sort,
				Order:                   params.Order,
				OriginDutyLocation:      params.OriginDutyLocation,
			}

			// Let's set default values for page and perPage if we don't get arguments for them. We'll use 1 for page and 20
			// for perPage.
			if params.Page == nil {
				listPaymentRequestParams.Page = models.Int64Pointer(1)
			}
			// Same for perPage
			if params.PerPage == nil {
				listPaymentRequestParams.PerPage = models.Int64Pointer(20)
			}

			paymentRequests, count, err := h.FetchPaymentRequestList(
				appCtx,
				appCtx.Session().OfficeUserID,
				&listPaymentRequestParams,
			)
			if err != nil {
				appCtx.Logger().
					Error("payment requests queue", zap.String("office_user_id", appCtx.Session().OfficeUserID.String()), zap.Error(err))
				return queues.NewGetPaymentRequestsQueueInternalServerError(), err
			}

			queuePaymentRequests := payloads.QueuePaymentRequests(paymentRequests)

			result := &ghcmessages.QueuePaymentRequestsResult{
				TotalCount:           int64(count),
				Page:                 int64(*listPaymentRequestParams.Page),
				PerPage:              int64(*listPaymentRequestParams.PerPage),
				QueuePaymentRequests: *queuePaymentRequests,
			}

			return queues.NewGetPaymentRequestsQueueOK().WithPayload(result), nil
		})
}

// GetServicesCounselingQueueHandler returns the moves for the Service Counselor queue user via GET /queues/counselor
type GetServicesCounselingQueueHandler struct {
	handlers.HandlerConfig
	services.OrderFetcher
}

// Handle returns the paginated list of moves for the services counselor
func (h GetServicesCounselingQueueHandler) Handle(
	params queues.GetServicesCounselingQueueParams,
) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			if !appCtx.Session().IsOfficeUser() ||
				!appCtx.Session().Roles.HasRole(roles.RoleTypeServicesCounselor) {
				forbiddenErr := apperror.NewForbiddenError(
					"user is not authenticated with an office role",
				)
				appCtx.Logger().Error(forbiddenErr.Error())
				return queues.NewGetServicesCounselingQueueForbidden(), forbiddenErr
			}

			ListOrderParams := services.ListOrderParams{
				Branch:                  params.Branch,
				Locator:                 params.Locator,
				DodID:                   params.DodID,
				LastName:                params.LastName,
				OriginDutyLocation:      params.OriginDutyLocation,
				DestinationDutyLocation: params.DestinationDutyLocation,
				OriginGBLOC:             params.OriginGBLOC,
				SubmittedAt:             handlers.FmtDateTimePtrToPopPtr(params.SubmittedAt),
				RequestedMoveDate:       params.RequestedMoveDate,
				Page:                    params.Page,
				PerPage:                 params.PerPage,
				Sort:                    params.Sort,
				Order:                   params.Order,
				NeedsPPMCloseout:        params.NeedsPPMCloseout,
				PPMType:                 params.PpmType,
				CloseoutInitiated:       handlers.FmtDateTimePtrToPopPtr(params.CloseoutInitiated),
				CloseoutLocation:        params.CloseoutLocation,
			}

			if params.NeedsPPMCloseout != nil && *params.NeedsPPMCloseout {
				ListOrderParams.Status = []string{string(models.MoveStatusAPPROVED)}
			} else if len(params.Status) == 0 {
				ListOrderParams.Status = []string{string(models.MoveStatusNeedsServiceCounseling)}
			} else {
				ListOrderParams.Status = params.Status
			}

			// Let's set default values for page and perPage if we don't get arguments for them. We'll use 1 for page and 20
			// for perPage.
			if params.Page == nil {
				ListOrderParams.Page = models.Int64Pointer(1)
			}
			// Same for perPage
			if params.PerPage == nil {
				ListOrderParams.PerPage = models.Int64Pointer(20)
			}

			moves, count, err := h.OrderFetcher.ListOrders(
				appCtx,
				appCtx.Session().OfficeUserID,
				&ListOrderParams,
			)
			if err != nil {
				appCtx.Logger().
					Error("error fetching list of moves for office user", zap.Error(err))
				return queues.NewGetServicesCounselingQueueInternalServerError(), err
			}

			queueMoves := payloads.QueueMoves(moves)

			result := &ghcmessages.QueueMovesResult{
				Page:       *ListOrderParams.Page,
				PerPage:    *ListOrderParams.PerPage,
				TotalCount: int64(count),
				QueueMoves: *queueMoves,
			}

			return queues.NewGetServicesCounselingQueueOK().WithPayload(result), nil
		})
}
