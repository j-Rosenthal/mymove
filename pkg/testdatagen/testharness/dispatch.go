package testharness

import (
	"errors"
	"sort"
	"sync"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/testdatagen"
)

type testHarnessResponse interface{}

type actionFunc func(appCtx appcontext.AppContext) testHarnessResponse

var actionDispatcher = map[string]actionFunc{
	"DefaultAdminUser": func(appCtx appcontext.AppContext) testHarnessResponse {
		return factory.BuildDefaultAdminUser(appCtx.DB())
	},
	"DefaultMove": func(appCtx appcontext.AppContext) testHarnessResponse {
		return factory.BuildMove(appCtx.DB(), nil, nil)
	},
	"MoveWithOrders": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeMoveWithOrders(appCtx.DB())
	},
	"SpouseProGearMove": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeSpouseProGearMove(appCtx.DB())
	},
	"WithShipmentMove": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeWithShipmentMove(appCtx)
	},
	"HHGMoveWithNTSAndNeedsSC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveWithNTSAndNeedsSC(appCtx)
	},
	"MoveWithMinimalNTSRNeedsSC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeMoveWithMinimalNTSRNeedsSC(appCtx)
	},
	"HHGMoveNeedsSC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveNeedsSC(appCtx)
	},
	"HHGMoveForSeparationNeedsSC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveForSeparationNeedsSC(appCtx)
	},
	"HHGMoveForRetireeNeedsSC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveForRetireeNeedsSC(appCtx)
	},
	"HHGMoveIn200DaysSIT": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveIn200DaysSIT(appCtx)
	},
	"HHGMoveIn200DaysSITWithPendingExtension": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveIn200DaysSITWithPendingExtension(appCtx)
	},
	"HHGMoveWithServiceItemsAndPaymentRequestsAndFilesForTOO": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveWithServiceItemsAndPaymentRequestsAndFilesForTOO(appCtx)
	},
	"HHGMoveWithRetireeForTOO": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveWithRetireeForTOO(appCtx)
	},
	"HHGMoveWithNTSShipmentsForTOO": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveWithNTSShipmentsForTOO(appCtx)
	},
	"MoveWithNTSShipmentsForTOO": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeMoveWithNTSShipmentsForTOO(appCtx)
	},
	"HHGMoveWithExternalNTSShipmentsForTOO": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveWithExternalNTSShipmentsForTOO(appCtx)
	},
	"HHGMoveWithApprovedNTSShipmentsForTOO": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveWithApprovedNTSShipmentsForTOO(appCtx)
	},
	"HHGMoveWithNTSRShipmentsForTOO": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveWithNTSRShipmentsForTOO(appCtx)
	},
	"HHGMoveWithApprovedNTSRShipmentsForTOO": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveWithApprovedNTSRShipmentsForTOO(appCtx)
	},
	"HHGMoveWithExternalNTSRShipmentsForTOO": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveWithExternalNTSRShipmentsForTOO(appCtx)
	},
	"HHGMoveWithServiceItemsandPaymentRequestsForTIO": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveWithServiceItemsandPaymentRequestsForTIO(appCtx)
	},
	"HHGMoveInSITWithAddressChangeRequestOver50Miles": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveInSITWithAddressChangeRequestOver50Miles(appCtx)
	},
	"HHGMoveInSITWithAddressChangeRequestUnder50Miles": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeHHGMoveInSITWithAddressChangeRequestUnder50Miles(appCtx)
	},
	"NTSRMoveWithPaymentRequest": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeNTSRMoveWithPaymentRequest(appCtx)
	},
	"NTSRMoveWithServiceItemsAndPaymentRequest": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeNTSRMoveWithServiceItemsAndPaymentRequest(appCtx)
	},
	"PrimeSimulatorMoveNeedsShipmentUpdate": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakePrimeSimulatorMoveNeedsShipmentUpdate(appCtx)
	},
	"NeedsOrdersUser": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeNeedsOrdersUser(appCtx.DB())
	},
	"PPMInProgressMove": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakePPMInProgressMove(appCtx)
	},
	"MoveWithPPMShipmentReadyForFinalCloseout": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeMoveWithPPMShipmentReadyForFinalCloseout(appCtx)
	},
	"PPMMoveWithCloseout": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakePPMMoveWithCloseout(appCtx)
	},
	"PPMMoveWithCloseoutOffice": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakePPMMoveWithCloseoutOffice(appCtx)
	},
	"ApprovedMoveWithPPM": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeApprovedMoveWithPPM(appCtx)
	},
	"SubmittedMoveWithPPMShipmentForSC": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeSubmittedMoveWithPPMShipmentForSC(appCtx)
	},
	"UnSubmittedMoveWithPPMShipmentThroughEstimatedWeights": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeUnSubmittedMoveWithPPMShipmentThroughEstimatedWeights(appCtx)
	},
	"ApprovedMoveWithPPMWithAboutFormComplete": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeApprovedMoveWithPPMWithAboutFormComplete(appCtx)
	},
	"UnsubmittedMoveWithMultipleFullPPMShipmentComplete": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeUnsubmittedMoveWithMultipleFullPPMShipmentComplete(appCtx)
	},
	"ApprovedMoveWithPPMProgearWeightTicket": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeApprovedMoveWithPPMProgearWeightTicket(appCtx)
	},
	"ApprovedMoveWithPPMProgearWeightTicketOffice": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeApprovedMoveWithPPMProgearWeightTicketOffice(appCtx)
	},
	"ApprovedMoveWithPPMWeightTicketOffice": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeApprovedMoveWithPPMWeightTicketOffice(appCtx)
	},
	"ApprovedMoveWithPPMMovingExpense": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeApprovedMoveWithPPMMovingExpense(appCtx)
	},
	"ApprovedMoveWithPPMMovingExpenseOffice": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeApprovedMoveWithPPMMovingExpenseOffice(appCtx)
	},
	"ApprovedMoveWithPPMAllDocTypesOffice": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeApprovedMoveWithPPMAllDocTypesOffice(appCtx)
	},
	"DraftMoveWithPPMWithDepartureDate": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeDraftMoveWithPPMWithDepartureDate(appCtx)
	},
	"OfficeUserWithTOOAndTIO": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeOfficeUserWithTOOAndTIO(appCtx)
	},
	"WebhookSubscription": func(appCtx appcontext.AppContext) testHarnessResponse {
		return testdatagen.MakeWebhookSubscription(appCtx.DB(), testdatagen.Assertions{})
	},
	"ApprovedMoveWithPPMShipmentAndExcessWeight": func(appCtx appcontext.AppContext) testHarnessResponse {
		return MakeApprovedMoveWithPPMShipmentAndExcessWeight(appCtx)
	},
}

func Actions() []string {
	actions := make([]string, 0, len(actionDispatcher))
	for k := range actionDispatcher {
		actions = append(actions, k)
	}
	sort.Strings(actions)
	return actions
}

var mutex sync.Mutex

func Dispatch(appCtx appcontext.AppContext, action string) (testHarnessResponse, error) {

	// ensure only one dispatch is running at a time in a heavy handed
	// way to prevent multiple setup functions from stomping on each
	// other when creating shared data (like duty locations)
	mutex.Lock()
	defer mutex.Unlock()

	dispatcher, ok := actionDispatcher[action]
	if !ok {
		appCtx.Logger().Error("Cannot find testharness dispatcher", zap.Any("action", action))
		return nil, errors.New("Cannot find testharness dispatcher for action: `" + action + "`")
	}

	appCtx.Logger().Info("Found testharness dispatcher", zap.Any("action", action))
	return dispatcher(appCtx), nil

}
