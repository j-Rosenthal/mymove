package paymentrequest

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

// verify that the MoveTaskOrderID on the payment request is not a nil uuid
func checkMTOIDField() paymentRequestValidator {
	return paymentRequestValidatorFunc(func(_ appcontext.AppContext, paymentRequest models.PaymentRequest, oldPaymentRequest *models.PaymentRequest) error {
		// Verify that the MTO ID exists
		if paymentRequest.MoveTaskOrderID == uuid.Nil {
			return apperror.NewInvalidCreateInputError(nil, "Invalid Create Input Error: MoveTaskOrderID is required on PaymentRequest create")
		}

		return nil
	})
}

func checkMTOIDMatchesServiceItemMTOID() paymentRequestValidator {
	return paymentRequestValidatorFunc(func(_ appcontext.AppContext, paymentRequest models.PaymentRequest, oldPaymentRequest *models.PaymentRequest) error {
		var paymentRequestServiceItems = paymentRequest.PaymentServiceItems
		for _, paymentRequestServiceItem := range paymentRequestServiceItems {
			if paymentRequest.MoveTaskOrderID != paymentRequestServiceItem.MTOServiceItem.MoveTaskOrderID && paymentRequestServiceItem.MTOServiceItemID != uuid.Nil {
				return apperror.NewConflictError(paymentRequestServiceItem.MTOServiceItem.MoveTaskOrderID, "Conflict Error: Payment Request MoveTaskOrderID does not match Service Item MoveTaskOrderID")
			}
		}
		return nil
	})
}
