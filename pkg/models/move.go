package models

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/db/dberr"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/random"
	"github.com/transcom/mymove/pkg/unit"
)

// MoveStatus represents the status of an order record's lifecycle
type MoveStatus string

const (
	// MoveStatusDRAFT captures enum value "DRAFT"
	MoveStatusDRAFT MoveStatus = "DRAFT"
	// MoveStatusSUBMITTED captures enum value "SUBMITTED"
	MoveStatusSUBMITTED MoveStatus = "SUBMITTED"
	// MoveStatusAPPROVED captures enum value "APPROVED"
	MoveStatusAPPROVED MoveStatus = "APPROVED"
	// MoveStatusCANCELED captures enum value "CANCELED"
	MoveStatusCANCELED MoveStatus = "CANCELED"
	// MoveStatusAPPROVALSREQUESTED captures enum value "APPROVALS REQUESTED"
	MoveStatusAPPROVALSREQUESTED MoveStatus = "APPROVALS REQUESTED"
	// MoveStatusNeedsServiceCounseling captures enum value "NEEDS SERVICE COUNSELING"
	MoveStatusNeedsServiceCounseling MoveStatus = "NEEDS SERVICE COUNSELING"
	// MoveStatusServiceCounselingCompleted captures enum value "SERVICE COUNSELING COMPLETED"
	MoveStatusServiceCounselingCompleted MoveStatus = "SERVICE COUNSELING COMPLETED"
)

const maxLocatorAttempts = 3
const locatorLength = 6

// This set of letters should produce 'non-word' type strings
var locatorLetters = []rune("346789BCDFGHJKMPQRTVWXY")

// Move is an object representing a move
type Move struct {
	ID                           uuid.UUID               `json:"id" db:"id"`
	Locator                      string                  `json:"locator" db:"locator"`
	CreatedAt                    time.Time               `json:"created_at" db:"created_at"`
	UpdatedAt                    time.Time               `json:"updated_at" db:"updated_at"`
	SubmittedAt                  *time.Time              `json:"submitted_at" db:"submitted_at"`
	OrdersID                     uuid.UUID               `json:"orders_id" db:"orders_id"`
	Orders                       Order                   `belongs_to:"orders" fk_id:"orders_id"`
	PersonallyProcuredMoves      PersonallyProcuredMoves `has_many:"personally_procured_moves" fk_id:"move_id" order_by:"created_at desc"`
	Status                       MoveStatus              `json:"status" db:"status"`
	SignedCertifications         SignedCertifications    `has_many:"signed_certifications" fk_id:"move_id" order_by:"created_at desc"`
	CancelReason                 *string                 `json:"cancel_reason" db:"cancel_reason"`
	Show                         *bool                   `json:"show" db:"show"`
	TIORemarks                   *string                 `db:"tio_remarks"`
	AvailableToPrimeAt           *time.Time              `db:"available_to_prime_at"`
	ContractorID                 *uuid.UUID              `db:"contractor_id"`
	Contractor                   *Contractor             `belongs_to:"contractors" fk_id:"contractor_id"`
	PPMEstimatedWeight           *unit.Pound             `db:"ppm_estimated_weight"`
	PPMType                      *string                 `db:"ppm_type"`
	MTOServiceItems              MTOServiceItems         `has_many:"mto_service_items" fk_id:"move_id"`
	PaymentRequests              PaymentRequests         `has_many:"payment_requests" fk_id:"move_id"`
	MTOShipments                 MTOShipments            `has_many:"mto_shipments" fk_id:"move_id"`
	ReferenceID                  *string                 `db:"reference_id"`
	ServiceCounselingCompletedAt *time.Time              `db:"service_counseling_completed_at"`
	PrimeCounselingCompletedAt   *time.Time              `db:"prime_counseling_completed_at"`
	ExcessWeightQualifiedAt      *time.Time              `db:"excess_weight_qualified_at"`
	ExcessWeightUploadID         *uuid.UUID              `db:"excess_weight_upload_id"`
	ExcessWeightUpload           *Upload                 `belongs_to:"uploads" fk_id:"excess_weight_upload_id"`
	ExcessWeightAcknowledgedAt   *time.Time              `db:"excess_weight_acknowledged_at"`
	BillableWeightsReviewedAt    *time.Time              `db:"billable_weights_reviewed_at"`
	FinancialReviewFlag          bool                    `db:"financial_review_flag"`
	FinancialReviewFlagSetAt     *time.Time              `db:"financial_review_flag_set_at"`
	FinancialReviewRemarks       *string                 `db:"financial_review_remarks"`
	ShipmentGBLOC                MoveToGBLOCs            `has_many:"move_to_gbloc" fk_id:"move_id"`
	CloseoutOfficeID             *uuid.UUID              `db:"closeout_office_id"`
	CloseoutOffice               *TransportationOffice   `belongs_to:"transportation_offices" fk_id:"closeout_office_id"`
	ApprovalsRequestedAt         *time.Time              `db:"approvals_requested_at"`
}

// TableName overrides the table name used by Pop.
func (m Move) TableName() string {
	return "moves"
}

// MoveOptions is used when creating new moves based on parameters
type MoveOptions struct {
	Show *bool
}

type Moves []Move

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (m *Move) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: m.Locator, Name: "Locator"},
		&validators.UUIDIsPresent{Field: m.OrdersID, Name: "OrdersID"},
		&validators.StringIsPresent{Field: string(m.Status), Name: "Status"},
		&OptionalTimeIsPresent{Field: m.ExcessWeightQualifiedAt, Name: "ExcessWeightQualifiedAt"},
		&OptionalUUIDIsPresent{Field: m.ExcessWeightUploadID, Name: "ExcessWeightUploadID"},
		&OptionalUUIDIsPresent{Field: m.CloseoutOfficeID, Name: "CloseoutOfficeID"},
	), nil
}

// FetchMove fetches and validates a Move for this User
func FetchMove(db *pop.Connection, session *auth.Session, id uuid.UUID) (*Move, error) {
	var move Move

	err := db.Q().Eager(
		"SignedCertifications",
		"Orders.ServiceMember",
		"Orders.UploadedAmendedOrders",
		"CloseoutOffice",
	).Where("show = TRUE").Find(&move, id)

	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return nil, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return nil, err
	}

	var shipments MTOShipments
	err = db.Q().Scope(utilities.ExcludeDeletedScope()).Eager("MTOAgents",
		"PickupAddress",
		"SecondaryPickupAddress",
		"DestinationAddress",
		"SecondaryDeliveryAddress",
		"PPMShipment").Where("mto_shipments.move_id = ?", move.ID).All(&shipments)

	if err != nil {
		return nil, err
	}

	move.MTOShipments = shipments

	// Ensure that the logged-in user is authorized to access this move
	if session.IsMilApp() && move.Orders.ServiceMember.ID != session.ServiceMemberID {
		return nil, ErrFetchForbidden
	}

	return &move, nil
}

// CreatePPM creates a new PPM associated with this move
func (m Move) CreatePPM(db *pop.Connection,
	weightEstimate *unit.Pound,
	originalMoveDate *time.Time,
	pickupPostalCode *string,
	hasAdditionalPostalCode *bool,
	additionalPickupPostalCode *string,
	destinationPostalCode *string,
	hasSit *bool,
	daysInStorage *int64,
	estimatedStorageReimbursement *string,
	hasRequestedAdvance bool,
	advance *Reimbursement) (*PersonallyProcuredMove, *validate.Errors, error) {

	newPPM := PersonallyProcuredMove{
		MoveID:                        m.ID,
		Move:                          m,
		WeightEstimate:                weightEstimate,
		OriginalMoveDate:              originalMoveDate,
		PickupPostalCode:              pickupPostalCode,
		HasAdditionalPostalCode:       hasAdditionalPostalCode,
		AdditionalPickupPostalCode:    additionalPickupPostalCode,
		DestinationPostalCode:         destinationPostalCode,
		HasSit:                        hasSit,
		DaysInStorage:                 daysInStorage,
		Status:                        PPMStatusDRAFT,
		HasRequestedAdvance:           hasRequestedAdvance,
		Advance:                       advance,
		EstimatedStorageReimbursement: estimatedStorageReimbursement,
	}

	verrs, err := SavePersonallyProcuredMove(db, &newPPM)
	if err != nil || verrs.HasAny() {
		return nil, verrs, err
	}

	return &newPPM, verrs, nil
}

// CreateSignedCertification creates a new SignedCertification associated with this move
func (m Move) CreateSignedCertification(db *pop.Connection,
	submittingUserID uuid.UUID,
	certificationText string,
	signature string,
	date time.Time,
	ppmID *uuid.UUID,
	certificationType *SignedCertificationType) (*SignedCertification, *validate.Errors, error) {

	newSignedCertification := SignedCertification{
		MoveID:                   m.ID,
		PersonallyProcuredMoveID: ppmID,
		CertificationType:        certificationType,
		SubmittingUserID:         submittingUserID,
		CertificationText:        certificationText,
		Signature:                signature,
		Date:                     date,
	}

	verrs, err := db.ValidateAndCreate(&newSignedCertification)
	if err != nil || verrs.HasAny() {
		return nil, verrs, err
	}

	return &newSignedCertification, verrs, nil
}

// GetMovesForUserID gets all move models for a given user ID
func GetMovesForUserID(db *pop.Connection, userID uuid.UUID) (Moves, error) {
	var moves Moves
	query := db.Where("user_id = $1", userID)
	err := query.All(&moves)
	return moves, err
}

// GenerateLocator constructs a record locator - a unique 6 character alphanumeric string
func GenerateLocator() string {
	// Get a UUID as a source of (almost certainly) unique bytes
	seed, err := uuid.NewV4()
	if err != nil {
		return ""
	}
	// Scramble them via SHA256 in case UUID has structure
	scrambledBytes := sha256.Sum256(seed.Bytes())
	// Now convert bytes to letters
	locatorRunes := make([]rune, locatorLength)
	for idx := 0; idx < locatorLength; idx++ {
		j := int(scrambledBytes[idx]) % len(locatorLetters)
		locatorRunes[idx] = locatorLetters[j]
	}
	return string(locatorRunes)
}

// createNewMove adds a new Move record into the DB. In the (unlikely) event that we have a clash on Locators we
// retry with a new record locator.
func createNewMove(db *pop.Connection,
	orders Order,
	moveOptions MoveOptions) (*Move, *validate.Errors, error) {

	show := BoolPointer(true)
	if moveOptions.Show != nil {
		show = moveOptions.Show
	}

	var contractor Contractor
	err := db.Where("type='Prime'").First(&contractor)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not find contractor: %w", err)
	}

	referenceID, err := GenerateReferenceID(db)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not generate a unique ReferenceID: %w", err)
	}

	for i := 0; i < maxLocatorAttempts; i++ {
		move := Move{
			Orders:       orders,
			OrdersID:     orders.ID,
			Locator:      GenerateLocator(),
			Status:       MoveStatusDRAFT,
			Show:         show,
			ContractorID: &contractor.ID,
			ReferenceID:  &referenceID,
		}
		verrs, err := db.ValidateAndCreate(&move)
		if verrs.HasAny() {
			return nil, verrs, nil
		}
		if err != nil {
			if dberr.IsDBErrorForConstraint(err, pgerrcode.UniqueViolation, "moves_locator_idx") {
				// If we have a collision, try again for maxLocatorAttempts
				continue
			}
			return nil, verrs, err
		}

		return &move, verrs, nil
	}
	// the only way we get here is if we got a unique constraint error maxLocatorAttempts times.
	verrs := validate.NewErrors()
	return nil, verrs, ErrLocatorGeneration
}

// GenerateReferenceID generates a reference ID for the MTO
func GenerateReferenceID(db *pop.Connection) (string, error) {
	const maxAttempts = 10
	var referenceID string
	var err error
	for i := 0; i < maxAttempts; i++ {
		referenceID, err = generateReferenceIDHelper(db)
		if err == nil {
			return referenceID, nil
		}
	}
	return "", fmt.Errorf("move: failed to generate reference id; %w", err)
}

// GenerateReferenceID creates a random ID for an MTO. Format (xxxx-xxxx) with X being a number 0-9 (ex. 0009-1234. 4321-4444)
func generateReferenceIDHelper(db *pop.Connection) (string, error) {
	min := 0
	max := 10000
	firstNum, err := random.GetRandomIntAddend(min, max)
	if err != nil {
		return "", err
	}

	secondNum, err := random.GetRandomIntAddend(min, max)
	if err != nil {
		return "", err
	}

	newReferenceID := fmt.Sprintf("%04d-%04d", firstNum, secondNum)

	exists, err := db.Where(`reference_id= $1`, newReferenceID).Exists(&Move{})

	if err != nil {
		return "", err
	} else if exists {
		return "", errors.New("move: reference_id already exists")
	}

	return newReferenceID, nil
}

// SaveMoveDependencies safely saves a Move status, ppms' advances' statuses, orders statuses,
// and shipment GBLOCs.
func SaveMoveDependencies(db *pop.Connection, move *Move) (*validate.Errors, error) {
	responseVErrors := validate.NewErrors()
	var responseError error

	for _, ppm := range move.PersonallyProcuredMoves {
		copyOfPpm := ppm // Make copy to avoid implicit memory aliasing of items from a range statement.
		if copyOfPpm.Advance != nil {
			if verrs, err := db.ValidateAndSave(copyOfPpm.Advance); verrs.HasAny() || err != nil {
				responseVErrors.Append(verrs)
				responseError = errors.Wrap(err, "Error Saving Advance")
			}
		}

		if verrs, err := db.ValidateAndSave(&copyOfPpm); verrs.HasAny() || err != nil {
			responseVErrors.Append(verrs)
			responseError = errors.Wrap(err, "Error Saving PPM")
		}
	}

	if verrs, err := db.ValidateAndSave(&move.Orders); verrs.HasAny() || err != nil {
		responseVErrors.Append(verrs)
		responseError = errors.Wrap(err, "Error Saving Orders")
	}

	if verrs, err := db.ValidateAndSave(move); verrs.HasAny() || err != nil {
		responseVErrors.Append(verrs)
		responseError = errors.Wrap(err, "Error Saving Move")
	}

	return responseVErrors, responseError
}

// FetchMoveForMoveDates returns a Move along with all the associations needed to determine
// the move dates summary information.
func FetchMoveForMoveDates(db *pop.Connection, moveID uuid.UUID) (Move, error) {
	var move Move
	err := db.
		Eager(
			"Orders.ServiceMember.DutyLocation.Address",
			"Orders.NewDutyLocation.Address",
			"Orders.ServiceMember",
		).
		Find(&move, moveID)

	return move, err
}

// FetchMoveByOrderID returns a Move for a given id
func FetchMoveByOrderID(db *pop.Connection, orderID uuid.UUID) (Move, error) {
	var move Move
	err := db.Where("orders_id = ?", orderID).First(&move)
	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return Move{}, ErrFetchNotFound
		}
		return Move{}, err
	}
	return move, nil
}

// FetchMoveByMoveID returns a Move for a given id
func FetchMoveByMoveID(db *pop.Connection, moveID uuid.UUID) (Move, error) {
	var move Move
	err := db.Q().Find(&move, moveID)

	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return Move{}, ErrFetchNotFound
		}
		return Move{}, err
	}
	return move, nil
}

// IsCanceled returns true if the Move's status is `CANCELED`, false otherwise
func (m Move) IsCanceled() *bool {
	if m.Status == MoveStatusCANCELED {
		return BoolPointer(true)
	}
	return BoolPointer(false)
}

// IsPPMOnly returns true of the only type of shipment associate with the move is "PPM", false otherwise
func (m Move) IsPPMOnly() bool {
	if len(m.MTOShipments) == 0 {
		return false
	}
	ppmOnlyMove := true
	for _, s := range m.MTOShipments {
		if s.ShipmentType != MTOShipmentTypePPM {
			ppmOnlyMove = false
			break
		}
	}
	return ppmOnlyMove
}
