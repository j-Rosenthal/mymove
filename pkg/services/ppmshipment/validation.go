package ppmshipment

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type ppmShipmentValidator interface {
	// Validate checks the newPPMShipment for adherence to business rules. The
	// oldPPMShipment parameter is expected to be nil in creation use cases.
	// It is safe to return a *validate.Errors with zero added errors as
	// a success case.
	Validate(a appcontext.AppContext, newPPMShipment models.PPMShipment, oldPPMShipment *models.PPMShipment, shipment *models.MTOShipment) error
}

// validatePPMShipment checks a PPM shipment against a passed-in set of business rule checks
func validatePPMShipment(
	appCtx appcontext.AppContext,
	newPPMShipment models.PPMShipment,
	oldPPMShipment *models.PPMShipment,
	shipment *models.MTOShipment,
	checks ...ppmShipmentValidator,
) (result error) {
	verrs := validate.NewErrors()
	for _, checker := range checks {
		if err := checker.Validate(appCtx, newPPMShipment, oldPPMShipment, shipment); err != nil {
			switch e := err.(type) {
			case *validate.Errors:
				// accumulate validation errors
				verrs.Append(e)
			default:
				// non-validation errors have priority,
				// and short-circuit doing any further checks
				return err
			}
		}
	}
	if verrs.HasAny() {
		result = apperror.NewInvalidInputError(newPPMShipment.ID, nil, verrs, "Invalid input found while validating the PPM shipment.")
	}

	return result
}

// ppmShipmentValidatorFunc is an adapter type for converting a function into an implementation of ppmShipmentValidator
type ppmShipmentValidatorFunc func(appcontext.AppContext, models.PPMShipment, *models.PPMShipment, *models.MTOShipment) error

// Validate fulfills the ppmShipmentValidator interface
func (fn ppmShipmentValidatorFunc) Validate(appCtx appcontext.AppContext, newer models.PPMShipment, older *models.PPMShipment, ship *models.MTOShipment) error {
	return fn(appCtx, newer, older, ship)
}

func mergePPMShipment(newPPMShipment models.PPMShipment, oldPPMShipment *models.PPMShipment) *models.PPMShipment {
	if oldPPMShipment == nil {
		return &newPPMShipment
	}

	ppmShipment := *oldPPMShipment

	ppmShipment.ActualMoveDate = services.SetOptionalDateTimeField(newPPMShipment.ActualMoveDate, ppmShipment.ActualMoveDate)

	ppmShipment.SecondaryPickupPostalCode = services.SetOptionalStringField(newPPMShipment.SecondaryPickupPostalCode, ppmShipment.SecondaryPickupPostalCode)
	ppmShipment.ActualPickupPostalCode = services.SetOptionalStringField(newPPMShipment.ActualPickupPostalCode, ppmShipment.ActualPickupPostalCode)
	ppmShipment.SecondaryDestinationPostalCode = services.SetOptionalStringField(newPPMShipment.SecondaryDestinationPostalCode, ppmShipment.SecondaryDestinationPostalCode)
	ppmShipment.ActualDestinationPostalCode = services.SetOptionalStringField(newPPMShipment.ActualDestinationPostalCode, ppmShipment.ActualDestinationPostalCode)
	ppmShipment.HasProGear = services.SetNoNilOptionalBoolField(newPPMShipment.HasProGear, ppmShipment.HasProGear)
	ppmShipment.EstimatedWeight = services.SetNoNilOptionalPoundField(newPPMShipment.EstimatedWeight, ppmShipment.EstimatedWeight)
	ppmShipment.ProGearWeight = services.SetNoNilOptionalPoundField(newPPMShipment.ProGearWeight, ppmShipment.ProGearWeight)
	ppmShipment.SpouseProGearWeight = services.SetNoNilOptionalPoundField(newPPMShipment.SpouseProGearWeight, ppmShipment.SpouseProGearWeight)
	ppmShipment.EstimatedIncentive = services.SetNoNilOptionalCentField(newPPMShipment.EstimatedIncentive, ppmShipment.EstimatedIncentive)
	ppmShipment.HasRequestedAdvance = services.SetNoNilOptionalBoolField(newPPMShipment.HasRequestedAdvance, ppmShipment.HasRequestedAdvance)
	ppmShipment.AdvanceAmountRequested = services.SetNoNilOptionalCentField(newPPMShipment.AdvanceAmountRequested, ppmShipment.AdvanceAmountRequested)
	ppmShipment.FinalIncentive = services.SetNoNilOptionalCentField(newPPMShipment.FinalIncentive, ppmShipment.FinalIncentive)
	ppmShipment.HasReceivedAdvance = services.SetNoNilOptionalBoolField(newPPMShipment.HasReceivedAdvance, ppmShipment.HasReceivedAdvance)
	ppmShipment.AdvanceAmountReceived = services.SetNoNilOptionalCentField(newPPMShipment.AdvanceAmountReceived, ppmShipment.AdvanceAmountReceived)

	ppmShipment.SITExpected = services.SetNoNilOptionalBoolField(newPPMShipment.SITExpected, ppmShipment.SITExpected)
	ppmShipment.SITEstimatedWeight = services.SetNoNilOptionalPoundField(newPPMShipment.SITEstimatedWeight, ppmShipment.SITEstimatedWeight)
	ppmShipment.SITEstimatedEntryDate = services.SetOptionalDateTimeField(newPPMShipment.SITEstimatedEntryDate, ppmShipment.SITEstimatedEntryDate)
	ppmShipment.SITEstimatedDepartureDate = services.SetOptionalDateTimeField(newPPMShipment.SITEstimatedDepartureDate, ppmShipment.SITEstimatedDepartureDate)

	if newPPMShipment.SITLocation != nil {
		ppmShipment.SITLocation = newPPMShipment.SITLocation
	}

	if newPPMShipment.AdvanceStatus != nil {
		ppmShipment.AdvanceStatus = newPPMShipment.AdvanceStatus
	}

	if newPPMShipment.W2Address != nil {
		ppmShipment.W2Address = newPPMShipment.W2Address
		if ppmShipment.W2AddressID != nil {
			ppmShipment.W2Address.ID = *ppmShipment.W2AddressID
		} else {
			ppmShipment.W2Address.ID = uuid.Nil
		}
	}

	if ppmShipment.SITExpected != nil && !*ppmShipment.SITExpected {
		ppmShipment.SITLocation = nil
		ppmShipment.SITEstimatedWeight = nil
		ppmShipment.SITEstimatedEntryDate = nil
		ppmShipment.SITEstimatedDepartureDate = nil
		ppmShipment.SITEstimatedCost = nil
	}

	if ppmShipment.HasProGear != nil && !*ppmShipment.HasProGear {
		ppmShipment.ProGearWeight = nil
		ppmShipment.SpouseProGearWeight = nil
	}

	if ppmShipment.HasRequestedAdvance != nil && !*ppmShipment.HasRequestedAdvance {
		ppmShipment.AdvanceAmountRequested = nil
	}

	if ppmShipment.HasReceivedAdvance != nil && !*ppmShipment.HasReceivedAdvance {
		ppmShipment.AdvanceAmountReceived = nil
	}

	if !newPPMShipment.ExpectedDepartureDate.IsZero() {
		ppmShipment.ExpectedDepartureDate = newPPMShipment.ExpectedDepartureDate
	}

	if newPPMShipment.PickupPostalCode != "" {
		ppmShipment.PickupPostalCode = newPPMShipment.PickupPostalCode
	}
	if newPPMShipment.DestinationPostalCode != "" {
		ppmShipment.DestinationPostalCode = newPPMShipment.DestinationPostalCode
	}

	if newPPMShipment.WeightTickets != nil && len(newPPMShipment.WeightTickets) >= 1 {
		ppmShipment.WeightTickets = newPPMShipment.WeightTickets
	}

	return &ppmShipment
}
