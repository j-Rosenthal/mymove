package factory

import (
	"log"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

type mtoServiceItemBuildType byte

const (
	mtoServiceItemBuildBasic mtoServiceItemBuildType = iota
	mtoServiceItemBuildExtended
)

// buildMTOServiceItemWithBuildType creates a single MTOServiceItem.
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
func buildMTOServiceItemWithBuildType(db *pop.Connection, customs []Customization, traits []Trait, buildType mtoServiceItemBuildType) models.MTOServiceItem {
	customs = setupCustomizations(customs, traits)

	// Find address customization and extract the custom address
	var cMTOServiceItem models.MTOServiceItem
	if result := findValidCustomization(customs, MTOServiceItem); result != nil {
		cMTOServiceItem = result.Model.(models.MTOServiceItem)
		if result.LinkOnly {
			return cMTOServiceItem
		}
	}

	var mtoShipmentID *uuid.UUID
	var mtoShipment models.MTOShipment
	var move models.Move
	if buildType == mtoServiceItemBuildExtended {
		// BuildMTOShipment creates a move as necessary
		mtoShipment = BuildMTOShipment(db, customs, traits)
		mtoShipmentID = &mtoShipment.ID
		move = mtoShipment.MoveTaskOrder
	} else {
		move = BuildMove(db, customs, traits)
	}

	var reService models.ReService
	if result := findValidCustomization(customs, ReService); result != nil {
		reService = FetchOrBuildReService(db, customs, nil)
	} else {
		reService = FetchOrBuildReServiceByCode(db, models.ReServiceCode("STEST"))
	}

	// Create default MTOServiceItem
	mtoServiceItem := models.MTOServiceItem{
		MoveTaskOrder:   move,
		MoveTaskOrderID: move.ID,
		MTOShipment:     mtoShipment,
		MTOShipmentID:   mtoShipmentID,
		ReService:       reService,
		ReServiceID:     reService.ID,
		Status:          models.MTOServiceItemStatusSubmitted,
	}

	// only set SITOriginHHGOriginalAddress if a customization is provided
	if result := findValidCustomization(customs, Addresses.SITOriginHHGOriginalAddress); result != nil {
		addressCustoms := convertCustomizationInList(customs, Addresses.SITOriginHHGOriginalAddress, Address)
		address := BuildAddress(db, addressCustoms, traits)
		mtoServiceItem.SITOriginHHGOriginalAddress = &address
		mtoServiceItem.SITOriginHHGOriginalAddressID = &address.ID
	}

	// only set SITOriginHHGActualAddress if a customization is provided
	if result := findValidCustomization(customs, Addresses.SITOriginHHGActualAddress); result != nil {
		addressCustoms := convertCustomizationInList(customs, Addresses.SITOriginHHGActualAddress, Address)
		address := BuildAddress(db, addressCustoms, traits)
		mtoServiceItem.SITOriginHHGActualAddress = &address
		mtoServiceItem.SITOriginHHGActualAddressID = &address.ID
	}

	// only set SITDestinationFinalAddress if a customization is provided
	if result := findValidCustomization(customs, Addresses.SITDestinationFinalAddress); result != nil {
		addressCustoms := convertCustomizationInList(customs, Addresses.SITDestinationFinalAddress, Address)
		address := BuildAddress(db, addressCustoms, traits)
		mtoServiceItem.SITDestinationFinalAddress = &address
		mtoServiceItem.SITDestinationFinalAddressID = &address.ID
	}

	// only set SITDestinationOriginalAddress if a customization is provided
	if result := findValidCustomization(customs, Addresses.SITDestinationOriginalAddress); result != nil {
		addressCustoms := convertCustomizationInList(customs, Addresses.SITDestinationOriginalAddress, Address)
		address := BuildAddress(db, addressCustoms, traits)
		mtoServiceItem.SITDestinationOriginalAddress = &address
		mtoServiceItem.SITDestinationOriginalAddressID = &address.ID
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&mtoServiceItem, cMTOServiceItem)

	// If db is false, it's a stub. No need to create in database.
	if db != nil {
		mustCreate(db, &mtoServiceItem)
	}

	return mtoServiceItem
}

// BuildMTOServiceItem creates a single extended MTOServiceItem
func BuildMTOServiceItem(db *pop.Connection, customs []Customization, traits []Trait) models.MTOServiceItem {
	return buildMTOServiceItemWithBuildType(db, customs, traits, mtoServiceItemBuildExtended)
}

// BuildMTOServiceItemBasic creates a single basic MTOServiceItem
func BuildMTOServiceItemBasic(db *pop.Connection, customs []Customization, traits []Trait) models.MTOServiceItem {
	return buildMTOServiceItemWithBuildType(db, customs, traits, mtoServiceItemBuildBasic)
}

// Needed by BuildRealMTOServiceItemWithAllDeps

var (
	paramActualPickupDate = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameActualPickupDate,
		Description: "actual pickup date",
		Type:        models.ServiceItemParamTypeDate,
		Origin:      models.ServiceItemParamOriginPrime,
	}
	paramContractCode = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameContractCode,
		Description: "contract code",
		Type:        models.ServiceItemParamTypeString,
		Origin:      models.ServiceItemParamOriginSystem,
	}
	paramContractYearName = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameContractYearName,
		Description: "contract year name",
		Type:        models.ServiceItemParamTypeString,
		Origin:      models.ServiceItemParamOriginPricer,
	}
	paramDistanceZip = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameDistanceZip,
		Description: "distance zip",
		Type:        models.ServiceItemParamTypeInteger,
		Origin:      models.ServiceItemParamOriginSystem,
	}
	paramEIAFuelPrice = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameEIAFuelPrice,
		Description: "eia fuel price",
		Type:        models.ServiceItemParamTypeInteger,
		Origin:      models.ServiceItemParamOriginSystem,
	}
	paramEscalationCompounded = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameEscalationCompounded,
		Description: "escalation compounded",
		Type:        models.ServiceItemParamTypeDecimal,
		Origin:      models.ServiceItemParamOriginPricer,
	}
	paramFSCMultiplier = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameFSCMultiplier,
		Description: "fsc multiplier",
		Type:        models.ServiceItemParamTypeDecimal,
		Origin:      models.ServiceItemParamOriginPricer,
	}
	paramFSCPriceDifferenceInCents = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameFSCPriceDifferenceInCents,
		Description: "fsc price difference in cents",
		Type:        models.ServiceItemParamTypeDecimal,
		Origin:      models.ServiceItemParamOriginPricer,
	}
	paramFSCWeightBasedDistanceMultiplier = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier,
		Description: "fsc weight based multiplier",
		Type:        models.ServiceItemParamTypeDecimal,
		Origin:      models.ServiceItemParamOriginSystem,
	}
	paramIsPeak = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameIsPeak,
		Description: "is peak",
		Type:        models.ServiceItemParamTypeBoolean,
		Origin:      models.ServiceItemParamOriginPricer,
	}
	paramMTOAvailableAToPrimeAt = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameMTOAvailableToPrimeAt,
		Description: "mto available to prime at",
		Type:        models.ServiceItemParamTypeTimestamp,
		Origin:      models.ServiceItemParamOriginSystem,
	}
	paramNumberDaysSIT = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameNumberDaysSIT,
		Description: "number days SIT",
		Type:        models.ServiceItemParamTypeInteger,
		Origin:      models.ServiceItemParamOriginPrime,
	}
	paramPriceRateOrFactor = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNamePriceRateOrFactor,
		Description: "price, rate, or factor",
		Type:        models.ServiceItemParamTypeDecimal,
		Origin:      models.ServiceItemParamOriginPricer,
	}
	paramReferenceDate = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameReferenceDate,
		Description: "reference date",
		Type:        models.ServiceItemParamTypeDate,
		Origin:      models.ServiceItemParamOriginSystem,
	}
	paramRequestedPickupDate = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameRequestedPickupDate,
		Description: "requested pickup date",
		Type:        models.ServiceItemParamTypeDate,
		Origin:      models.ServiceItemParamOriginPrime,
	}
	paramServiceAreaOrigin = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameServiceAreaOrigin,
		Description: "service area origin",
		Type:        models.ServiceItemParamTypeString,
		Origin:      models.ServiceItemParamOriginSystem,
	}
	paramServicesScheduleOrigin = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameServicesScheduleOrigin,
		Description: "services schedule origin",
		Type:        models.ServiceItemParamTypeInteger,
		Origin:      models.ServiceItemParamOriginSystem,
	}
	paramSITPaymentRequestEnd = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameSITPaymentRequestEnd,
		Description: "SIT payment request end",
		Type:        models.ServiceItemParamTypeDate,
		Origin:      models.ServiceItemParamOriginPaymentRequest,
	}
	paramSITPaymentRequestStart = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameSITPaymentRequestStart,
		Description: "SIT payment request start",
		Type:        models.ServiceItemParamTypeDate,
		Origin:      models.ServiceItemParamOriginPaymentRequest,
	}
	paramWeightAdjusted = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameWeightAdjusted,
		Description: "weight adjusted",
		Type:        models.ServiceItemParamTypeInteger,
		Origin:      models.ServiceItemParamOriginSystem,
	}
	paramWeightBilled = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameWeightBilled,
		Description: "weight billed",
		Type:        models.ServiceItemParamTypeInteger,
		Origin:      models.ServiceItemParamOriginSystem,
	}
	paramWeightEstimated = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameWeightEstimated,
		Description: "weight estimated",
		Type:        models.ServiceItemParamTypeInteger,
		Origin:      models.ServiceItemParamOriginPrime,
	}
	paramWeightOriginal = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameWeightOriginal,
		Description: "weight original",
		Type:        models.ServiceItemParamTypeInteger,
		Origin:      models.ServiceItemParamOriginPrime,
	}
	paramWeightReweigh = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameWeightReweigh,
		Description: "weight reweigh",
		Type:        models.ServiceItemParamTypeInteger,
		Origin:      models.ServiceItemParamOriginPrime,
	}
	paramZipDestAddress = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameZipDestAddress,
		Description: "zip dest address",
		Type:        models.ServiceItemParamTypeString,
		Origin:      models.ServiceItemParamOriginPrime,
	}
	paramZipPickupAddress = models.ServiceItemParamKey{
		Key:         models.ServiceItemParamNameZipPickupAddress,
		Description: "zip pickup address",
		Type:        models.ServiceItemParamTypeString,
		Origin:      models.ServiceItemParamOriginPrime,
	}
	fixtureServiceItemParamsMap = map[models.ReServiceCode]models.ServiceItemParamKeys{
		models.ReServiceCodeCS: {
			paramContractCode,
			paramMTOAvailableAToPrimeAt,
			paramPriceRateOrFactor,
		},
		models.ReServiceCodeMS: {
			paramContractCode,
			paramMTOAvailableAToPrimeAt,
			paramPriceRateOrFactor,
		},
		models.ReServiceCodeDLH: {
			paramActualPickupDate,
			paramContractCode,
			paramContractYearName,
			paramDistanceZip,
			paramEscalationCompounded,
			paramIsPeak,
			paramPriceRateOrFactor,
			paramReferenceDate,
			paramRequestedPickupDate,
			paramServiceAreaOrigin,
			paramWeightAdjusted,
			paramWeightBilled,
			paramWeightEstimated,
			paramWeightOriginal,
			paramWeightReweigh,
			paramZipDestAddress,
			paramZipPickupAddress,
		},
		models.ReServiceCodeFSC: {
			paramActualPickupDate,
			paramContractCode,
			paramDistanceZip,
			paramEIAFuelPrice,
			paramFSCMultiplier,
			paramFSCPriceDifferenceInCents,
			paramFSCWeightBasedDistanceMultiplier,
			paramWeightAdjusted,
			paramWeightBilled,
			paramWeightEstimated,
			paramWeightOriginal,
			paramWeightReweigh,
			paramZipDestAddress,
			paramZipPickupAddress,
		},
		models.ReServiceCodeDPK: {
			paramActualPickupDate,
			paramContractCode,
			paramContractYearName,
			paramEscalationCompounded,
			paramIsPeak,
			paramPriceRateOrFactor,
			paramReferenceDate,
			paramRequestedPickupDate,
			paramServiceAreaOrigin,
			paramServicesScheduleOrigin,
			paramWeightAdjusted,
			paramWeightBilled,
			paramWeightEstimated,
			paramWeightOriginal,
			paramWeightReweigh,
			paramZipPickupAddress,
		},
		models.ReServiceCodeDOP: {
			paramActualPickupDate,
			paramContractCode,
			paramContractYearName,
			paramEscalationCompounded,
			paramIsPeak,
			paramPriceRateOrFactor,
			paramReferenceDate,
			paramRequestedPickupDate,
			paramServiceAreaOrigin,
			paramWeightAdjusted,
			paramWeightBilled,
			paramWeightEstimated,
			paramWeightOriginal,
			paramWeightReweigh,
			paramZipPickupAddress,
		},
		models.ReServiceCodeDOASIT: {
			paramActualPickupDate,
			paramContractCode,
			paramContractYearName,
			paramEscalationCompounded,
			paramIsPeak,
			paramNumberDaysSIT,
			paramPriceRateOrFactor,
			paramReferenceDate,
			paramRequestedPickupDate,
			paramServiceAreaOrigin,
			paramSITPaymentRequestEnd,
			paramSITPaymentRequestStart,
			paramWeightAdjusted,
			paramWeightBilled,
			paramWeightEstimated,
			paramWeightOriginal,
			paramWeightReweigh,
			paramZipPickupAddress,
		},
	}
)

func BuildRealMTOServiceItemWithAllDeps(db *pop.Connection, serviceCode models.ReServiceCode, mto models.Move, mtoShipment models.MTOShipment) models.MTOServiceItem {
	// look up the service item param keys we need
	if serviceItemParamKeys, ok := fixtureServiceItemParamsMap[serviceCode]; ok {
		// get or create the ReService
		reService := FetchOrBuildReServiceByCode(db, serviceCode)

		// create all params defined for this particular service
		for _, serviceParamKeyToCreate := range serviceItemParamKeys {
			serviceItemParamKey := FetchOrBuildServiceItemParamKey(db, []Customization{
				{
					Model: serviceParamKeyToCreate,
				},
			}, nil)
			BuildServiceParam(db, []Customization{
				{
					Model:    reService,
					LinkOnly: true,
				},
				{
					Model:    serviceItemParamKey,
					LinkOnly: true,
				},
			}, nil)
		}

		// create a service item and return it
		mtoServiceItem := BuildMTOServiceItem(db, []Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model:    mtoShipment,
				LinkOnly: true,
			},
			{
				Model:    reService,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					Status: models.MTOServiceItemStatusApproved,
				},
			},
		}, nil)

		return mtoServiceItem
	}

	log.Panicf("couldn't create service item service code %s not defined", serviceCode)
	return models.MTOServiceItem{}

}

// BuildFullDLHMTOServiceItems makes a DLH type service item along with
// all its expected parameters returns the created move and all
// service items
//
// NOTE: the original did an override of the MTOShipment.Status to
// ensure it was Approved, but that is now the responsibility of the
// caller
func BuildFullDLHMTOServiceItems(db *pop.Connection, customs []Customization, traits []Trait) (models.Move, models.MTOServiceItems) {

	mtoShipment := BuildMTOShipment(db, customs, traits)

	move := mtoShipment.MoveTaskOrder
	move.MTOShipments = models.MTOShipments{mtoShipment}

	var mtoServiceItems models.MTOServiceItems
	// Service Item MS
	mtoServiceItemMS := BuildRealMTOServiceItemWithAllDeps(db,
		models.ReServiceCodeMS, move, mtoShipment)
	mtoServiceItems = append(mtoServiceItems, mtoServiceItemMS)
	// Service Item CS
	mtoServiceItemCS := BuildRealMTOServiceItemWithAllDeps(db,
		models.ReServiceCodeCS, move, mtoShipment)
	mtoServiceItems = append(mtoServiceItems, mtoServiceItemCS)
	// Service Item DLH
	mtoServiceItemDLH := BuildRealMTOServiceItemWithAllDeps(db,
		models.ReServiceCodeDLH, move, mtoShipment)
	mtoServiceItems = append(mtoServiceItems, mtoServiceItemDLH)
	// Service Item FSC
	mtoServiceItemFSC := BuildRealMTOServiceItemWithAllDeps(db,
		models.ReServiceCodeFSC, move, mtoShipment)
	mtoServiceItems = append(mtoServiceItems, mtoServiceItemFSC)

	return move, mtoServiceItems
}

// BuildFullOriginMTOServiceItems (follow-on to
// BuildFullDLHMTOServiceItem) makes a DLH type service item along
// with all its expected parameters returns the created move and all
// service items
//
// NOTE: the original did an override of the MTOShipment.Status to
// ensure it was Approved, but that is now the responsibility of the
// caller
func BuildFullOriginMTOServiceItems(db *pop.Connection, customs []Customization, traits []Trait) (models.Move, models.MTOServiceItems) {
	mtoShipment := BuildMTOShipment(db, customs, traits)

	move := mtoShipment.MoveTaskOrder
	move.MTOShipments = models.MTOShipments{mtoShipment}

	var mtoServiceItems models.MTOServiceItems
	// Service Item DPK
	mtoServiceItemDPK := BuildRealMTOServiceItemWithAllDeps(db,
		models.ReServiceCodeDPK, move, mtoShipment)
	mtoServiceItems = append(mtoServiceItems, mtoServiceItemDPK)
	// Service Item DOP
	mtoServiceItemDOP := BuildRealMTOServiceItemWithAllDeps(db,
		models.ReServiceCodeDOP, move, mtoShipment)
	mtoServiceItems = append(mtoServiceItems, mtoServiceItemDOP)

	return move, mtoServiceItems
}
