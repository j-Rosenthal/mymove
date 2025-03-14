package models_test

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ModelSuite) TestPPMShipmentValidation() {
	validPPMShipmentStatuses := strings.Join(models.AllowedPPMShipmentStatuses, ", ")
	validPPMShipmentAdvanceStatuses := strings.Join(models.AllowedPPMAdvanceStatuses, ", ")
	validSITLocations := strings.Join(models.AllowedSITLocationTypes, ", ")

	blankAdvanceStatus := models.PPMAdvanceStatus("")
	blankSITLocation := models.SITLocationType("")

	testCases := map[string]struct {
		ppmShipment  models.PPMShipment
		expectedErrs map[string][]string
	}{
		"Successful Minimal Validation": {
			ppmShipment: models.PPMShipment{
				ShipmentID:            uuid.Must(uuid.NewV4()),
				ExpectedDepartureDate: testdatagen.PeakRateCycleStart,
				Status:                models.PPMShipmentStatusDraft,
				PickupPostalCode:      "90210",
				DestinationPostalCode: "94535",
			},
			expectedErrs: nil,
		},
		"Missing Required Fields": {
			ppmShipment: models.PPMShipment{},
			expectedErrs: map[string][]string{
				"shipment_id":             {"ShipmentID can not be blank."},
				"expected_departure_date": {"ExpectedDepartureDate can not be blank."},
				"status":                  {fmt.Sprintf("Status is not in the list [%s].", validPPMShipmentStatuses)},
				"pickup_postal_code":      {"PickupPostalCode can not be blank."},
				"destination_postal_code": {"DestinationPostalCode can not be blank."},
			},
		},
		"Optional fields raise errors with invalid values": {
			ppmShipment: models.PPMShipment{
				// Setting up min required fields here so that we don't get these in our errors.
				ShipmentID:            uuid.Must(uuid.NewV4()),
				ExpectedDepartureDate: testdatagen.PeakRateCycleStart,
				Status:                models.PPMShipmentStatusDraft,
				PickupPostalCode:      "90210",
				DestinationPostalCode: "94535",

				// Now setting optional fields with invalid values.
				DeletedAt:                      models.TimePointer(time.Time{}),
				ActualMoveDate:                 models.TimePointer(time.Time{}),
				SubmittedAt:                    models.TimePointer(time.Time{}),
				ReviewedAt:                     models.TimePointer(time.Time{}),
				ApprovedAt:                     models.TimePointer(time.Time{}),
				W2AddressID:                    models.UUIDPointer(uuid.Nil),
				SecondaryPickupPostalCode:      models.StringPointer(""),
				ActualPickupPostalCode:         models.StringPointer(""),
				SecondaryDestinationPostalCode: models.StringPointer(""),
				ActualDestinationPostalCode:    models.StringPointer(""),
				EstimatedWeight:                models.PoundPointer(unit.Pound(-1)),
				ProGearWeight:                  models.PoundPointer(unit.Pound(-1)),
				SpouseProGearWeight:            models.PoundPointer(unit.Pound(-1)),
				EstimatedIncentive:             models.CentPointer(unit.Cents(0)),
				FinalIncentive:                 models.CentPointer(unit.Cents(0)),
				AdvanceAmountRequested:         models.CentPointer(unit.Cents(0)),
				AdvanceStatus:                  &blankAdvanceStatus,
				AdvanceAmountReceived:          models.CentPointer(unit.Cents(0)),
				SITLocation:                    &blankSITLocation,
				SITEstimatedWeight:             models.PoundPointer(unit.Pound(-1)),
				SITEstimatedEntryDate:          models.TimePointer(time.Time{}),
				SITEstimatedDepartureDate:      models.TimePointer(time.Time{}),
				SITEstimatedCost:               models.CentPointer(unit.Cents(0)),
				AOAPacketID:                    models.UUIDPointer(uuid.Nil),
				PaymentPacketID:                models.UUIDPointer(uuid.Nil),
			},
			expectedErrs: map[string][]string{
				"deleted_at":                        {"DeletedAt can not be blank."},
				"actual_move_date":                  {"ActualMoveDate can not be blank."},
				"submitted_at":                      {"SubmittedAt can not be blank."},
				"reviewed_at":                       {"ReviewedAt can not be blank."},
				"approved_at":                       {"ApprovedAt can not be blank."},
				"w2_address_id":                     {"W2AddressID can not be blank."},
				"secondary_pickup_postal_code":      {"SecondaryPickupPostalCode can not be blank."},
				"actual_pickup_postal_code":         {"ActualPickupPostalCode can not be blank."},
				"secondary_destination_postal_code": {"SecondaryDestinationPostalCode can not be blank."},
				"actual_destination_postal_code":    {"ActualDestinationPostalCode can not be blank."},
				"estimated_weight":                  {"-1 is less than zero."},
				"pro_gear_weight":                   {"-1 is less than zero."},
				"spouse_pro_gear_weight":            {"-1 is less than zero."},
				"estimated_incentive":               {"EstimatedIncentive must be greater than zero, got: 0."},
				"final_incentive":                   {"FinalIncentive must be greater than zero, got: 0."},
				"advance_amount_requested":          {"AdvanceAmountRequested must be greater than zero, got: 0."},
				"advance_status":                    {fmt.Sprintf("AdvanceStatus is not in the list [%s].", validPPMShipmentAdvanceStatuses)},
				"advance_amount_received":           {"AdvanceAmountReceived must be greater than zero, got: 0."},
				"sitlocation":                       {fmt.Sprintf("SITLocation is not in the list [%s].", validSITLocations)},
				"sitestimated_weight":               {"-1 is less than zero."},
				"sitestimated_entry_date":           {"SITEstimatedEntryDate can not be blank."},
				"sitestimated_departure_date":       {"SITEstimatedDepartureDate can not be blank."},
				"sitestimated_cost":                 {"SITEstimatedCost must be greater than zero, got: 0."},
				"aoapacket_id":                      {"AOAPacketID can not be blank."},
				"payment_packet_id":                 {"PaymentPacketID can not be blank."},
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		suite.Run(name, func() {
			suite.verifyValidationErrors(testCase.ppmShipment, testCase.expectedErrs)
		})
	}
}
