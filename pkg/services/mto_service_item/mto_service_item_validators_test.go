package mtoserviceitem

import (
	"fmt"
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MTOServiceItemServiceSuite) TestUpdateMTOServiceItemData() {

	// Set up the data needed for updateMTOServiceItemData obj
	checker := movetaskorder.NewMoveTaskOrderChecker()
	now := time.Now()
	setupTestData := func() (models.MTOServiceItem, models.MTOServiceItem) {
		// Create a service item to serve as the old object
		oldServiceItem := testdatagen.MakeDefaultMTOServiceItem(suite.DB())
		oldServiceItem.CustomerContacts = models.MTOServiceItemCustomerContacts{
			models.MTOServiceItemCustomerContact{
				Type:                       models.CustomerContactTypeFirst,
				TimeMilitary:               "1300Z",
				FirstAvailableDeliveryDate: time.Now().AddDate(0, 0, 3),
			},
		}
		// Shallow copy service item to create the "updated" object
		updatedServiceItem := oldServiceItem
		return oldServiceItem, updatedServiceItem
	}

	// Test successful check for linked IDs
	suite.Run("checkLinkedIDs - success", func() {
		// Under test:  checkLinkedIDs function, which checks that two linked
		//              service items have the same move, shipment, or reService IDs
		// Set up:      Create a service item and compare against another
		// Expected outcome: PreconditionFailedError
		oldServiceItem, newServiceItem := setupTestData()
		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem: newServiceItem, // as-is, should succeed
			oldServiceItem:     oldServiceItem,
			verrs:              validate.NewErrors(),
		}
		err := serviceItemData.checkLinkedIDs(suite.AppContextForTest())

		suite.NoError(err)
		suite.NoVerrs(serviceItemData.verrs)
	})

	// Test unsuccessful check for linked IDs
	suite.Run("checkLinkedIDs - failure", func() {
		oldServiceItem, newServiceItem := setupTestData()
		fakeUUID := uuid.FromStringOrNil("00010001-0001-0001-0001-000100010001")
		newServiceItem.MoveTaskOrderID = fakeUUID
		newServiceItem.MTOShipmentID = &fakeUUID
		newServiceItem.ReServiceID = fakeUUID

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem: newServiceItem,
			oldServiceItem:     oldServiceItem,
			verrs:              validate.NewErrors(),
		}
		err := serviceItemData.checkLinkedIDs(suite.AppContextForTest())

		suite.NoError(err)
		suite.True(serviceItemData.verrs.HasAny())
		suite.Contains(serviceItemData.verrs.Keys(), "moveTaskOrderID")
		suite.Contains(serviceItemData.verrs.Keys(), "mtoShipmentID")
		suite.Contains(serviceItemData.verrs.Keys(), "reServiceID")
	})

	// Test successful check for Prime availability
	suite.Run("checkPrimeAvailability - success", func() {
		oldServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil),
				LinkOnly: true,
			},
		}, nil)
		newServiceItemPrime := oldServiceItemPrime // Shallow copy model

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem:  newServiceItemPrime,
			oldServiceItem:      oldServiceItemPrime,
			availabilityChecker: checker,
			verrs:               validate.NewErrors(),
		}
		err := serviceItemData.checkPrimeAvailability(suite.AppContextForTest())

		suite.NoError(err)
		suite.NoVerrs(serviceItemData.verrs)
	})

	// Test unsuccessful check for Prime availability
	suite.Run("checkPrimeAvailability - failure", func() {
		oldServiceItem, newServiceItem := setupTestData()

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem:  newServiceItem, // the default errorServiceItem should not be Prime-available
			oldServiceItem:      oldServiceItem,
			availabilityChecker: checker,
			verrs:               validate.NewErrors(),
		}
		err := serviceItemData.checkPrimeAvailability(suite.AppContextForTest())

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.NoVerrs(serviceItemData.verrs) // this check doesn't add a validation error
	})

	// Test successful check for non-Prime fields
	suite.Run("checkNonPrimeFields - success", func() {
		oldServiceItem, newServiceItem := setupTestData() // These

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem: newServiceItem, // as-is, should succeed because all the values are the same
			oldServiceItem:     oldServiceItem,
			verrs:              validate.NewErrors(),
		}
		err := serviceItemData.checkNonPrimeFields(suite.AppContextForTest())

		suite.NoError(err)
		suite.NoVerrs(serviceItemData.verrs)
	})

	// Test unsuccessful check for non-Prime fields
	suite.Run("checkNonPrimeFields - failure", func() {
		// Update the non-updateable fields:
		oldServiceItem, newServiceItem := setupTestData() // These

		newServiceItem.Status = models.MTOServiceItemStatusApproved
		newServiceItem.RejectionReason = handlers.FmtString("reason")
		newServiceItem.ApprovedAt = &now
		newServiceItem.RejectedAt = &now

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem: newServiceItem,
			oldServiceItem:     oldServiceItem,
			verrs:              validate.NewErrors(),
		}
		err := serviceItemData.checkNonPrimeFields(suite.AppContextForTest())

		suite.NoError(err)
		suite.True(serviceItemData.verrs.HasAny())
		suite.Contains(serviceItemData.verrs.Keys(), "status")
		suite.Contains(serviceItemData.verrs.Keys(), "rejectionReason")
		suite.Contains(serviceItemData.verrs.Keys(), "approvedAt")
		suite.Contains(serviceItemData.verrs.Keys(), "rejectedAt")
	})

	// Test successful check for SIT departure service item - not updating SITDepartureDate
	suite.Run("checkSITDeparture w/ no SITDepartureDate update - success", func() {
		oldServiceItem, newServiceItem := setupTestData() // These

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem: newServiceItem, // default is not DDDSIT/DOPSIT
			oldServiceItem:     oldServiceItem,
			verrs:              validate.NewErrors(),
		}
		err := serviceItemData.checkSITDeparture(suite.AppContextForTest())

		suite.NoError(err)
		suite.NoVerrs(serviceItemData.verrs)
	})

	// Test successful check for SIT departure service item - DDDSIT
	suite.Run("checkSITDeparture w/ DDDSIT - success", func() {
		// Under test:  checkSITDeparture checks that the service item is a
		//			    DDDSIT or DOPSIT if the user is trying to update the
		// 			    SITDepartureDate
		// Set up:      Create an old and new DDDSIT, with a new date and try to update.
		// Expected outcome: Success if both are DDDSIT
		oldDDDSIT := factory.BuildMTOServiceItem(nil, []factory.Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDDSIT,
				},
			},
		}, nil)
		newDDDSIT := oldDDDSIT
		newDDDSIT.SITDepartureDate = &now

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem: newDDDSIT,
			oldServiceItem:     oldDDDSIT,
			verrs:              validate.NewErrors(),
		}
		err := serviceItemData.checkSITDeparture(suite.AppContextForTest())

		suite.NoError(err)
		suite.NoVerrs(serviceItemData.verrs)
	})

	// Test unsuccessful check for SIT departure service item - not a departure SIT item
	suite.Run("checkSITDeparture w/ non-departure SIT - failure", func() {
		// Under test:  checkSITDeparture checks that the service item is a
		//			    DDDSIT or DOPSIT if the user is trying to update the
		// 			    SITDepartureDate
		// Set up:      Create any non DDDSIT service item
		// Expected outcome: Conflict Error
		oldDOFSIT := factory.BuildMTOServiceItem(nil, []factory.Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
				},
			},
		}, nil)
		newDOFSIT := oldDOFSIT
		newDOFSIT.SITDepartureDate = &now
		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem: newDOFSIT, // default is not DDDSIT/DOPSIT
			oldServiceItem:     oldDOFSIT,
			verrs:              validate.NewErrors(),
		}
		err := serviceItemData.checkSITDeparture(suite.AppContextForTest())

		suite.Error(err)
		suite.IsType(apperror.ConflictError{}, err)
		suite.NoVerrs(serviceItemData.verrs) // this check doesn't add a validation error
		suite.Contains(err.Error(), fmt.Sprintf("SIT Departure Date may only be manually updated for %s and %s service items", models.ReServiceCodeDDDSIT, models.ReServiceCodeDOPSIT))
	})

	// Test successful check for service item w/out payment request
	suite.Run("checkPaymentRequests - success", func() {
		// Under test:  checkPaymentRequests checks if there are payment requests
		//			    associated with this service item and returns a conflict error if so
		// Set up:      Create any service item with no payment requests
		// Expected outcome: No error
		oldServiceItem, newServiceItem := setupTestData() // These

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem: newServiceItem, // as-is, should succeed
			oldServiceItem:     oldServiceItem,
			verrs:              validate.NewErrors(),
		}
		err := serviceItemData.checkPaymentRequests(suite.AppContextForTest())

		suite.NoError(err)
		suite.NoVerrs(serviceItemData.verrs)
	})

	// Test unsuccessful check service item with an existing payment request
	suite.Run("checkPaymentRequests - failure", func() {
		// Under test:  checkPaymentRequests checks if there are payment requests
		//			    associated with this service item and returns a conflict error if so
		// Set up:      Create any service item with associated payment requests
		// Expected outcome: ConflictError
		oldServiceItem, newServiceItem := setupTestData() // These
		newServiceItem.Description = models.StringPointer("1234")

		paymentRequest := factory.BuildPaymentRequest(suite.DB(), nil, nil)
		factory.BuildPaymentServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    paymentRequest,
				LinkOnly: true,
			}, {
				Model:    oldServiceItem,
				LinkOnly: true,
			},
		}, nil)

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem: newServiceItem,
			oldServiceItem:     oldServiceItem,
			verrs:              validate.NewErrors(),
		}
		err := serviceItemData.checkPaymentRequests(suite.AppContextForTest())

		suite.Error(err)
		suite.IsType(apperror.ConflictError{}, err)
		suite.NoVerrs(serviceItemData.verrs) // this check doesn't add a validation error
		suite.Contains(err.Error(), "this service item has an existing payment request and can no longer be updated")
	})

	// Test getVerrs for successful example
	suite.Run("getVerrs - success", func() {
		// Under test:  getVerrs returns a list of validation errors
		// Set up:      Create a service item, run 2 validations that should pass
		// Expected outcome: No errors
		oldServiceItem, newServiceItem := setupTestData()

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem: newServiceItem,
			oldServiceItem:     oldServiceItem,
			verrs:              validate.NewErrors(),
		}
		_ = serviceItemData.checkLinkedIDs(suite.AppContextForTest()) // this test should pass regardless of potential errors here
		_ = serviceItemData.checkNonPrimeFields(suite.AppContextForTest())
		err := serviceItemData.getVerrs()

		suite.NoError(err)
		suite.NoVerrs(serviceItemData.verrs)
	})

	// Test getVerrs for unsuccessful example
	suite.Run("getVerrs - failure", func() {
		// Under test:  getVerrs returns a list of validation errors
		// Set up:      Create a service item, edit the non-prime fields and linked ids
		//              Run 2 validations that should fail
		// Expected outcome: InvalidInput error

		oldServiceItem, newServiceItem := setupTestData()

		// Change non prime fields
		newServiceItem.Status = models.MTOServiceItemStatusApproved
		newServiceItem.RejectionReason = handlers.FmtString("reason")
		newServiceItem.ApprovedAt = &now
		newServiceItem.RejectedAt = &now

		// Change linked ids
		fakeUUID := uuid.FromStringOrNil("00010001-0001-0001-0001-000100010001")
		newServiceItem.MoveTaskOrderID = fakeUUID
		newServiceItem.MTOShipmentID = &fakeUUID
		newServiceItem.ReServiceID = fakeUUID

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem: newServiceItem, // as-is, should fail
			oldServiceItem:     oldServiceItem,
			verrs:              validate.NewErrors(),
		}
		_ = serviceItemData.checkLinkedIDs(suite.AppContextForTest())
		_ = serviceItemData.checkNonPrimeFields(suite.AppContextForTest())
		err := serviceItemData.getVerrs()

		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.True(serviceItemData.verrs.HasAny())
		suite.Equal(7, serviceItemData.verrs.Count())
	})

	// Test setNewMTOServiceItem for successful example
	suite.Run("setNewMTOServiceItem - success", func() {
		oldServiceItem, editServiceItem := setupTestData() // These

		editServiceItem.Description = handlers.FmtString("testing update service item validators")
		editServiceItem.Reason = handlers.FmtString("")
		editServiceItem.SITEntryDate = &now
		editServiceItem.ApprovedAt = new(time.Time) // this is the zero time, what we need to nullify the field
		actualWeight := int64(4000)
		estimatedWeight := int64(4200)
		editServiceItem.ActualWeight = handlers.PoundPtrFromInt64Ptr(&actualWeight)
		editServiceItem.EstimatedWeight = handlers.PoundPtrFromInt64Ptr(&estimatedWeight)
		editServiceItem.CustomerContacts = models.MTOServiceItemCustomerContacts{
			models.MTOServiceItemCustomerContact{
				Type:                       models.CustomerContactTypeFirst,
				TimeMilitary:               "1400Z",
				FirstAvailableDeliveryDate: time.Now().AddDate(0, 0, 5),
			},
		}
		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem: editServiceItem,
			oldServiceItem:     oldServiceItem,
			verrs:              validate.NewErrors(),
		}
		newServiceItem := serviceItemData.setNewMTOServiceItem()

		suite.NoVerrs(serviceItemData.verrs)
		suite.Nil(newServiceItem.Reason)
		suite.Nil(newServiceItem.ApprovedAt)
		suite.Equal(newServiceItem.SITEntryDate, editServiceItem.SITEntryDate)
		suite.Equal(newServiceItem.Description, editServiceItem.Description)
		suite.NotEqual(newServiceItem.Description, oldServiceItem.Description)
		suite.NotEqual(newServiceItem.Description, serviceItemData.oldServiceItem.Description)
		suite.NotEqual(newServiceItem.CustomerContacts[0].TimeMilitary, serviceItemData.oldServiceItem.CustomerContacts[0].TimeMilitary)
		suite.NotEqual(newServiceItem.CustomerContacts[0].FirstAvailableDeliveryDate, serviceItemData.oldServiceItem.CustomerContacts[0].FirstAvailableDeliveryDate)
	})

	suite.Run("setNewMTOServiceItem - success with updating a service item that already has a sit destination final address", func() {
		oldServiceItem, editServiceItem := setupTestData()

		// Create the old address that has been saved to the db
		oldSitDestinationFinalAddress := factory.BuildAddress(suite.DB(), nil, nil)
		// Create an address that has not yet been saved to the db
		newSitDestinationFinalAddress := models.Address{
			StreetAddress1: "123 Any Street",
			StreetAddress2: models.StringPointer("P.O. Box 12345"),
			StreetAddress3: models.StringPointer("c/o Some Person"),
			City:           "Beverly Hills",
			State:          "CA",
			PostalCode:     "90210",
			Country:        models.StringPointer("US"),
		}

		// Set the old address and id to the old service item
		oldServiceItem.SITDestinationFinalAddress = &oldSitDestinationFinalAddress
		oldServiceItem.SITDestinationFinalAddressID = &oldSitDestinationFinalAddress.ID

		// Set the address to the new service item. We don't need to set the ID here because this replicates when
		// we are updating a sitDestinationFinalAddress for a service item that already has one.
		editServiceItem.SITDestinationFinalAddress = &newSitDestinationFinalAddress

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem: editServiceItem,
			oldServiceItem:     oldServiceItem,
			verrs:              validate.NewErrors(),
		}
		newServiceItem := serviceItemData.setNewMTOServiceItem()

		// Check that the IDs do not match the old address since we want to replace the record.
		suite.NotEqual(newServiceItem.SITDestinationFinalAddressID, &oldSitDestinationFinalAddress.ID)
		suite.NotEqual(newServiceItem.SITDestinationFinalAddress.ID, oldSitDestinationFinalAddress.ID)

		// Check that the address information matches the new address.
		suite.Equal(newServiceItem.SITDestinationFinalAddress.PostalCode, newSitDestinationFinalAddress.PostalCode)
		suite.Equal(newServiceItem.SITDestinationFinalAddress.StreetAddress1, newSitDestinationFinalAddress.StreetAddress1)
		suite.Equal(newServiceItem.SITDestinationFinalAddress.City, newSitDestinationFinalAddress.City)
	})

	suite.Run("setNewMTOServiceItem - success with updating a service item that does not have a sit destination final address", func() {
		oldServiceItem, editServiceItem := setupTestData()

		// Create an address that has not yet been saved to the db
		newSitDestinationFinalAddress := models.Address{
			StreetAddress1: "123 Any Street",
			StreetAddress2: models.StringPointer("P.O. Box 12345"),
			StreetAddress3: models.StringPointer("c/o Some Person"),
			City:           "Beverly Hills",
			State:          "CA",
			PostalCode:     "90210",
			Country:        models.StringPointer("US"),
		}

		// Set the address to the new service item. We don't need to set the ID here because this replicates when
		// we are updating a sitDestinationFinalAddress for a service item that already has one.
		editServiceItem.SITDestinationFinalAddress = &newSitDestinationFinalAddress

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem: editServiceItem,
			oldServiceItem:     oldServiceItem,
			verrs:              validate.NewErrors(),
		}
		newServiceItem := serviceItemData.setNewMTOServiceItem()
		nilUUID := uuid.Nil

		// Check that the IDs match the new address and that both are nil.
		suite.Equal(newServiceItem.SITDestinationFinalAddress.ID, newSitDestinationFinalAddress.ID)
		suite.Equal(nilUUID, newServiceItem.SITDestinationFinalAddress.ID)

		// Check that the address information matches the new address.
		suite.Equal(newServiceItem.SITDestinationFinalAddress.PostalCode, newSitDestinationFinalAddress.PostalCode)
		suite.Equal(newServiceItem.SITDestinationFinalAddress.StreetAddress1, newSitDestinationFinalAddress.StreetAddress1)
		suite.Equal(newServiceItem.SITDestinationFinalAddress.City, newSitDestinationFinalAddress.City)
	})

	suite.Run("setNewCustomerContacts - success with one old and one updated", func() {
		oldServiceItem, editServiceItem := setupTestData()

		editServiceItem.CustomerContacts = models.MTOServiceItemCustomerContacts{
			models.MTOServiceItemCustomerContact{
				Type:                       models.CustomerContactTypeFirst,
				TimeMilitary:               "1400Z",
				FirstAvailableDeliveryDate: time.Now().AddDate(0, 0, 5),
			},
		}
		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem: editServiceItem,
			oldServiceItem:     oldServiceItem,
			verrs:              validate.NewErrors(),
		}
		newCustomerContacts := serviceItemData.setNewCustomerContacts()

		suite.Equal(1, len(newCustomerContacts))
		suite.NotEqual(newCustomerContacts[0].TimeMilitary, serviceItemData.oldServiceItem.CustomerContacts[0].TimeMilitary)
		suite.NotEqual(newCustomerContacts[0].FirstAvailableDeliveryDate, serviceItemData.oldServiceItem.CustomerContacts[0].FirstAvailableDeliveryDate)

		suite.Equal(newCustomerContacts[0].TimeMilitary, serviceItemData.updatedServiceItem.CustomerContacts[0].TimeMilitary)
		suite.Equal(newCustomerContacts[0].FirstAvailableDeliveryDate, serviceItemData.updatedServiceItem.CustomerContacts[0].FirstAvailableDeliveryDate)
	})

	suite.Run("setNewCustomerContacts - success with one old and zero updated", func() {
		oldServiceItem, editServiceItem := setupTestData()

		editServiceItem.CustomerContacts = models.MTOServiceItemCustomerContacts{}
		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem: editServiceItem,
			oldServiceItem:     oldServiceItem,
			verrs:              validate.NewErrors(),
		}
		newCustomerContacts := serviceItemData.setNewCustomerContacts()

		suite.Equal(1, len(newCustomerContacts))
		suite.Equal(newCustomerContacts, serviceItemData.oldServiceItem.CustomerContacts)
	})

	suite.Run("setNewCustomerContacts - success with zero old and one updated", func() {
		oldServiceItem, editServiceItem := setupTestData()
		oldServiceItem.CustomerContacts = models.MTOServiceItemCustomerContacts{}

		editServiceItem.CustomerContacts = models.MTOServiceItemCustomerContacts{
			models.MTOServiceItemCustomerContact{
				Type:                       models.CustomerContactTypeFirst,
				TimeMilitary:               "1400Z",
				FirstAvailableDeliveryDate: time.Now().AddDate(0, 0, 5),
			},
		}
		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem: editServiceItem,
			oldServiceItem:     oldServiceItem,
			verrs:              validate.NewErrors(),
		}
		newCustomerContacts := serviceItemData.setNewCustomerContacts()

		suite.Equal(1, len(newCustomerContacts))
		suite.Equal(newCustomerContacts, serviceItemData.updatedServiceItem.CustomerContacts)
	})
	suite.Run("setNewCustomerContacts - success with updated having different type than old", func() {
		oldServiceItem, editServiceItem := setupTestData()

		editServiceItem.CustomerContacts = models.MTOServiceItemCustomerContacts{
			models.MTOServiceItemCustomerContact{
				Type:                       models.CustomerContactTypeSecond,
				TimeMilitary:               "1400Z",
				FirstAvailableDeliveryDate: time.Now().AddDate(0, 0, 5),
			},
		}
		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem: editServiceItem,
			oldServiceItem:     oldServiceItem,
			verrs:              validate.NewErrors(),
		}
		newCustomerContacts := serviceItemData.setNewCustomerContacts()

		// There should be two customer contacts
		suite.Equal(2, len(newCustomerContacts))
		for _, newContact := range newCustomerContacts {
			if newContact.Type == models.CustomerContactTypeFirst {
				suite.Equal(newContact.TimeMilitary, serviceItemData.oldServiceItem.CustomerContacts[0].TimeMilitary)
				suite.Equal(newContact.FirstAvailableDeliveryDate, serviceItemData.oldServiceItem.CustomerContacts[0].FirstAvailableDeliveryDate)
				suite.NotEqual(newContact.TimeMilitary, serviceItemData.updatedServiceItem.CustomerContacts[0].TimeMilitary)
				suite.NotEqual(newContact.FirstAvailableDeliveryDate, serviceItemData.updatedServiceItem.CustomerContacts[0].FirstAvailableDeliveryDate)
			}
			if newContact.Type == models.CustomerContactTypeSecond {
				suite.NotEqual(newContact.TimeMilitary, serviceItemData.oldServiceItem.CustomerContacts[0].TimeMilitary)
				suite.NotEqual(newContact.FirstAvailableDeliveryDate, serviceItemData.oldServiceItem.CustomerContacts[0].FirstAvailableDeliveryDate)
				suite.Equal(newContact.TimeMilitary, serviceItemData.updatedServiceItem.CustomerContacts[0].TimeMilitary)
				suite.Equal(newContact.FirstAvailableDeliveryDate, serviceItemData.updatedServiceItem.CustomerContacts[0].FirstAvailableDeliveryDate)
			}
		}
	})

	suite.Run("checkSITDestinationFinalAddress - adding SITDestinationFinalAddress", func() {
		oldServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil),
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDDSIT,
				},
			},
		}, nil)
		newServiceItemPrime := oldServiceItemPrime

		// Try to update SITDestinationFinalAddress
		newAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})
		newServiceItemPrime.SITDestinationFinalAddress = &newAddress

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem:  newServiceItemPrime,
			oldServiceItem:      oldServiceItemPrime,
			verrs:               validate.NewErrors(),
			availabilityChecker: checker,
		}
		err := serviceItemData.checkSITDestinationFinalAddress(suite.AppContextForTest())

		suite.NoError(err)
	})

	suite.Run("checkSITDestinationFinalAddress - invalid input failure: SITDestinationFinalAddress", func() {
		oldServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil),
				LinkOnly: true,
			},
			{
				Model: models.Address{},
				Type:  &factory.Addresses.SITDestinationFinalAddress,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDASIT,
				},
			},
		}, nil)
		newServiceItemPrime := oldServiceItemPrime

		// Try to update SITDestinationFinalAddress
		newAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})
		newServiceItemPrime.SITDestinationFinalAddress = &newAddress

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem:  newServiceItemPrime,
			oldServiceItem:      oldServiceItemPrime,
			verrs:               validate.NewErrors(),
			availabilityChecker: checker,
		}
		err := serviceItemData.checkSITDestinationFinalAddress(suite.AppContextForTest())

		suite.Error(err)
		conflictError := apperror.ConflictError{}
		suite.IsType(conflictError, err)
	})

	suite.Run("checkSITDestinationFinalAddress - invalid input failure: SITDestinationFinalAddress", func() {
		oldServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil),
				LinkOnly: true,
			},
			{
				Model: models.Address{},
				Type:  &factory.Addresses.SITDestinationFinalAddress,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDDSIT,
				},
			},
		}, nil)
		newServiceItemPrime := oldServiceItemPrime

		// Try to update SITDestinationFinalAddress
		newAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})
		newServiceItemPrime.SITDestinationFinalAddress = &newAddress

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem:  newServiceItemPrime,
			oldServiceItem:      oldServiceItemPrime,
			verrs:               validate.NewErrors(),
			availabilityChecker: checker,
		}
		err := serviceItemData.checkSITDestinationFinalAddress(suite.AppContextForTest())

		suite.NoError(err)
		suite.True(serviceItemData.verrs.HasAny())
		suite.Contains(serviceItemData.verrs.Keys(), "SITDestinationFinalAddress")
	})

	suite.Run("checkSITDestinationOriginalAddress - invalid input failure: adding SITDestinationOriginalAddress", func() {
		oldServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil),
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDDSIT,
				},
			},
		}, nil)
		newServiceItemPrime := oldServiceItemPrime

		// Try to add SITDestinationOriginalAddress
		newAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})
		newServiceItemPrime.SITDestinationOriginalAddress = &newAddress

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem:  newServiceItemPrime,
			oldServiceItem:      oldServiceItemPrime,
			verrs:               validate.NewErrors(),
			availabilityChecker: checker,
		}
		err := serviceItemData.checkSITDestinationOriginalAddress(suite.AppContextForTest())

		suite.NoError(err)
		suite.True(serviceItemData.verrs.HasAny())
		suite.Contains(serviceItemData.verrs.Keys(), "SITDestinationOriginalAddress")
	})

	suite.Run("checkSITDestinationOriginalAddress - invalid input failure: updating SITDestinationOriginalAddress", func() {
		oldServiceItemPrime := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil),
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDDSIT,
				},
			},
			{
				Model: models.Address{},
				Type:  &factory.Addresses.SITDestinationOriginalAddress,
			},
		}, nil)
		newServiceItemPrime := oldServiceItemPrime

		// Try to update SITDestinationOriginalAddress
		newAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})
		newServiceItemPrime.SITDestinationOriginalAddress = &newAddress

		serviceItemData := updateMTOServiceItemData{
			updatedServiceItem:  newServiceItemPrime,
			oldServiceItem:      oldServiceItemPrime,
			verrs:               validate.NewErrors(),
			availabilityChecker: checker,
		}
		err := serviceItemData.checkSITDestinationOriginalAddress(suite.AppContextForTest())

		suite.NoError(err)
		suite.True(serviceItemData.verrs.HasAny())
		suite.Contains(serviceItemData.verrs.Keys(), "SITDestinationOriginalAddress")
	})
}
