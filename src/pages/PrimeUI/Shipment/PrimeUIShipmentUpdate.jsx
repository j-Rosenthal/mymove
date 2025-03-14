import React, { useState } from 'react';
import { Formik } from 'formik';
import * as Yup from 'yup';
import { useNavigate, useParams, generatePath } from 'react-router-dom';
import { useMutation } from '@tanstack/react-query';
import { Alert, Button, Grid, GridContainer } from '@trussworks/react-uswds';
import { connect } from 'react-redux';
import { func } from 'prop-types';

import primeStyles from 'pages/PrimeUI/Prime.module.scss';
import { usePrimeSimulatorGetMove } from 'hooks/queries';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { primeSimulatorRoutes } from 'constants/routes';
import scrollToTop from 'shared/scrollToTop';
import { updatePrimeMTOShipment, updatePrimeMTOShipmentStatus } from 'services/primeApi';
import styles from 'components/Office/CustomerContactInfoForm/CustomerContactInfoForm.module.scss';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import { addressSchema, InvalidZIPTypeError, ZIP5_CODE_REGEX } from 'utils/validation';
import { isEmpty, isValidWeight } from 'shared/utils';
import { formatAddressForPrimeAPI, formatSwaggerDate, fromPrimeAPIAddressFormat } from 'utils/formatters';
import PrimeUIShipmentUpdateForm from 'pages/PrimeUI/Shipment/PrimeUIShipmentUpdateForm';
import PrimeUIShipmentUpdatePPMForm from 'pages/PrimeUI/Shipment/PrimeUIShipmentUpdatePPMForm';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';
import { SHIPMENT_OPTIONS } from 'shared/constants';

const PrimeUIShipmentUpdate = ({ setFlashMessage }) => {
  const [errorMessage, setErrorMessage] = useState();
  const { moveCodeOrID, shipmentId } = useParams();
  const { moveTaskOrder, isLoading, isError } = usePrimeSimulatorGetMove(moveCodeOrID);
  const mtoShipments = moveTaskOrder?.mtoShipments;
  const shipment = mtoShipments?.find((mtoShipment) => mtoShipment?.id === shipmentId);
  const navigate = useNavigate();

  const handleClose = () => {
    navigate(generatePath(primeSimulatorRoutes.VIEW_MOVE_PATH, { moveCodeOrID }));
  };

  const { mutateAsync: mutateMTOShipmentStatus } = useMutation(updatePrimeMTOShipmentStatus, {
    onSuccess: (updatedMTOShipment) => {
      mtoShipments[mtoShipments.findIndex((mtoShipment) => mtoShipment.id === updatedMTOShipment.id)] =
        updatedMTOShipment;
      setFlashMessage(`MSG_CANCELATION_SUCCESS${shipmentId}`, 'success', `Successfully canceled shipment`, '', true);
      handleClose();
    },
    // TODO: This method is duplicated for now. Refactor if neccessary.
    onError: (error) => {
      const { response: { body } = {} } = error;

      if (body) {
        /*
        {
          "detail": "Invalid data found in input",
          "instance":"00000000-0000-0000-0000-000000000000",
          "title":"Validation Error",
          "invalidFields": {
            "primeEstimatedWeight":["the time period for updating the estimated weight for a shipment has expired, please contact the TOO directly to request updates to this shipment’s estimated weight","Invalid Input."]
          }
        }
         */
        let invalidFieldsStr = '';
        if (body.invalidFields) {
          Object.keys(body.invalidFields).forEach((key) => {
            const value = body.invalidFields[key];
            invalidFieldsStr += `\n${key} - ${value && value.length > 0 ? value[0] : ''} ;`;
          });
        }
        setErrorMessage({
          title: `Prime API: ${body.title} `,
          detail: `${body.detail}${invalidFieldsStr}\n\nPlease cancel and Update Shipment again`,
        });
      } else {
        setErrorMessage({
          title: 'Unexpected error',
          detail: 'An unknown error has occurred, please check the state of the shipment and values',
        });
      }
      scrollToTop();
    },
  });

  const { mutateAsync: mutateMTOShipment } = useMutation(updatePrimeMTOShipment, {
    onSuccess: (updatedMTOShipment) => {
      mtoShipments[mtoShipments.findIndex((mtoShipment) => mtoShipment.id === updatedMTOShipment.id)] =
        updatedMTOShipment;
      setFlashMessage(`MSG_CREATE_PAYMENT_SUCCESS${shipmentId}`, 'success', `Successfully updated shipment`, '', true);
      handleClose();
    },
    onError: (error) => {
      const { response: { body } = {} } = error;

      if (body) {
        /*
        {
          "detail": "Invalid data found in input",
          "instance":"00000000-0000-0000-0000-000000000000",
          "title":"Validation Error",
          "invalidFields": {
            "primeEstimatedWeight":["the time period for updating the estimated weight for a shipment has expired, please contact the TOO directly to request updates to this shipment’s estimated weight","Invalid Input."]
          }
        }
         */
        let invalidFieldsStr = '';
        if (body.invalidFields) {
          Object.keys(body.invalidFields).forEach((key) => {
            const value = body.invalidFields[key];
            invalidFieldsStr += `\n${key} - ${value && value.length > 0 ? value[0] : ''} ;`;
          });
        }
        setErrorMessage({
          title: `Prime API: ${body.title} `,
          detail: `${body.detail}${invalidFieldsStr}\n\nPlease cancel and Update Shipment again`,
        });
      } else {
        setErrorMessage({
          title: 'Unexpected error',
          detail: 'An unknown error has occurred, please check the state of the shipment and values',
        });
      }
      scrollToTop();
    },
  });

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const isPPM = shipment.shipmentType === SHIPMENT_OPTIONS.PPM;

  const emptyAddress = {
    streetAddress1: '',
    streetAddress2: '',
    streetAddress3: '',
    city: '',
    state: '',
    postalCode: '',
  };

  const editableWeightEstimateField = !isValidWeight(shipment.primeEstimatedWeight);
  const editableWeightActualField = true;
  const reformatPrimeApiPickupAddress = fromPrimeAPIAddressFormat(shipment.pickupAddress);
  const reformatPrimeApiDestinationAddress = fromPrimeAPIAddressFormat(shipment.destinationAddress);
  const editablePickupAddress = isEmpty(reformatPrimeApiPickupAddress);
  const editableDestinationAddress = isEmpty(reformatPrimeApiDestinationAddress);

  const onCancelShipmentClick = () => {
    mutateMTOShipmentStatus({ mtoShipmentID: shipmentId, ifMatchETag: shipment.eTag }).then(() => {
      /* console.info("It's done and canceled."); */
    });
  };

  const onSubmit = (values, { setSubmitting }) => {
    let body;
    if (isPPM) {
      const {
        ppmShipment: {
          expectedDepartureDate,
          pickupPostalCode,
          secondaryPickupPostalCode,
          destinationPostalCode,
          secondaryDestinationPostalCode,
          sitExpected,
          sitLocation,
          sitEstimatedWeight,
          sitEstimatedEntryDate,
          sitEstimatedDepartureDate,
          estimatedWeight,
          hasProGear,
          proGearWeight,
          spouseProGearWeight,
        },
        counselorRemarks,
      } = values;
      body = {
        ppmShipment: {
          expectedDepartureDate: expectedDepartureDate ? formatSwaggerDate(expectedDepartureDate) : null,
          pickupPostalCode,
          secondaryPickupPostalCode: secondaryPickupPostalCode || null,
          destinationPostalCode,
          secondaryDestinationPostalCode: secondaryDestinationPostalCode || null,
          sitExpected,
          ...(sitExpected && {
            sitLocation: sitLocation || null,
            sitEstimatedWeight: sitEstimatedWeight ? parseInt(sitEstimatedWeight, 10) : null,
            sitEstimatedEntryDate: sitEstimatedEntryDate ? formatSwaggerDate(sitEstimatedEntryDate) : null,
            sitEstimatedDepartureDate: sitEstimatedDepartureDate ? formatSwaggerDate(sitEstimatedDepartureDate) : null,
          }),
          estimatedWeight: estimatedWeight ? parseInt(estimatedWeight, 10) : null,
          hasProGear,
          ...(hasProGear && {
            proGearWeight: proGearWeight ? parseInt(proGearWeight, 10) : null,
            spouseProGearWeight: spouseProGearWeight ? parseInt(spouseProGearWeight, 10) : null,
          }),
        },
        counselorRemarks: counselorRemarks || null,
      };
    } else {
      const {
        estimatedWeight,
        actualWeight,
        actualPickupDate,
        scheduledPickupDate,
        actualDeliveryDate,
        scheduledDeliveryDate,
        pickupAddress,
        destinationAddress,
        destinationType,
        diversion,
      } = values;

      body = {
        primeEstimatedWeight: editableWeightEstimateField ? parseInt(estimatedWeight, 10) : null,
        primeActualWeight: parseInt(actualWeight, 10),
        scheduledPickupDate: scheduledPickupDate ? formatSwaggerDate(scheduledPickupDate) : null,
        actualPickupDate: actualPickupDate ? formatSwaggerDate(actualPickupDate) : null,
        scheduledDeliveryDate: scheduledDeliveryDate ? formatSwaggerDate(scheduledDeliveryDate) : null,
        actualDeliveryDate: actualDeliveryDate ? formatSwaggerDate(actualDeliveryDate) : null,
        pickupAddress: editablePickupAddress ? formatAddressForPrimeAPI(pickupAddress) : null,
        destinationAddress: editableDestinationAddress ? formatAddressForPrimeAPI(destinationAddress) : null,
        destinationType,
        diversion,
      };
    }

    mutateMTOShipment({ mtoShipmentID: shipmentId, ifMatchETag: shipment.eTag, body }).then(() => {
      setSubmitting(false);
    });
  };

  let initialValues;
  let validationSchema;
  if (isPPM) {
    initialValues = {
      ppmShipment: {
        expectedDepartureDate: shipment.ppmShipment.expectedDepartureDate,
        pickupPostalCode: shipment.ppmShipment.pickupPostalCode || '',
        secondaryPickupPostalCode: shipment.ppmShipment.secondaryPickupPostalCode || '',
        destinationPostalCode: shipment.ppmShipment.destinationPostalCode || '',
        secondaryDestinationPostalCode: shipment.ppmShipment.secondaryDestinationPostalCode || '',
        sitExpected: shipment.ppmShipment.sitExpected,
        sitLocation: shipment.ppmShipment.sitLocation,
        sitEstimatedWeight: shipment.ppmShipment.sitEstimatedWeight?.toString(),
        sitEstimatedEntryDate: shipment.ppmShipment.sitEstimatedEntryDate,
        sitEstimatedDepartureDate: shipment.ppmShipment.sitEstimatedDepartureDate,
        estimatedWeight: shipment.ppmShipment.estimatedWeight?.toString(),
        hasProGear: shipment.ppmShipment.hasProGear,
        proGearWeight: shipment.ppmShipment.proGearWeight?.toString(),
        spouseProGearWeight: shipment.ppmShipment.spouseProGearWeight?.toString(),
      },
      counselorRemarks: shipment.counselorRemarks || '',
    };
    validationSchema = Yup.object().shape({
      ppmShipment: Yup.object().shape({
        expectedDepartureDate: Yup.date()
          .required('Required')
          .typeError('Invalid date. Must be in the format: DD MMM YYYY'),
        pickupPostalCode: Yup.string().matches(ZIP5_CODE_REGEX, InvalidZIPTypeError).required('Required'),
        secondaryPickupPostalCode: Yup.string().matches(ZIP5_CODE_REGEX, InvalidZIPTypeError).nullable(),
        destinationPostalCode: Yup.string().matches(ZIP5_CODE_REGEX, InvalidZIPTypeError).required('Required'),
        secondaryDestinationPostalCode: Yup.string().matches(ZIP5_CODE_REGEX, InvalidZIPTypeError).nullable(),
        sitExpected: Yup.boolean().required('Required'),
        sitLocation: Yup.string().when('sitExpected', {
          is: true,
          then: (schema) => schema.required('Required'),
        }),
        sitEstimatedWeight: Yup.number().when('sitExpected', {
          is: true,
          then: (schema) => schema.required('Required'),
        }),
        // TODO: Figure out how to validate this but be optional.  Right now, when you first
        //   go to the page with sitEnabled of false, the "Save" button remains disabled.
        // sitEstimatedEntryDate: Yup.date().when('sitExpected', {
        //   is: true,
        //   then: (schema) =>
        //     schema.typeError('Enter a complete date in DD MMM YYYY format (day, month, year).').required('Required'),
        // }),
        // sitEstimatedDepartureDate: Yup.date().when('sitExpected', {
        //   is: true,
        //   then: (schema) =>
        //     schema.typeError('Enter a complete date in DD MMM YYYY format (day, month, year).').required('Required'),
        // }),
        estimatedWeight: Yup.number().required('Required'),
        hasProGear: Yup.boolean().required('Required'),
        proGearWeight: Yup.number().when(['hasProGear', 'spouseProGearWeight'], {
          is: (hasProGear, spouseProGearWeight) => hasProGear && !spouseProGearWeight,
          then: (schema) =>
            schema.required(
              `Enter a weight into at least one pro-gear field. If you won't have pro-gear, uncheck above.`,
            ),
        }),
        spouseProGearWeight: Yup.number(),
      }),
      // counselorRemarks is an optional string
    });
  } else {
    initialValues = {
      estimatedWeight: shipment.primeEstimatedWeight?.toLocaleString(),
      actualWeight: shipment.primeActualWeight?.toLocaleString(),
      requestedPickupDate: shipment.requestedPickupDate,
      scheduledPickupDate: shipment.scheduledPickupDate,
      actualPickupDate: shipment.actualPickupDate,
      scheduledDeliveryDate: shipment.scheduledDeliveryDate,
      actualDeliveryDate: shipment.actualDeliveryDate,
      pickupAddress: editablePickupAddress ? emptyAddress : reformatPrimeApiPickupAddress,
      destinationAddress: editableDestinationAddress ? emptyAddress : reformatPrimeApiDestinationAddress,
      destinationType: shipment.destinationType,
      diversion: shipment.diversion,
    };

    validationSchema = Yup.object().shape({
      pickupAddress: addressSchema,
      destinationAddress: addressSchema,
      scheduledPickupDate: Yup.date().typeError('Invalid date. Must be in the format: DD MMM YYYY'),
      actualPickupDate: Yup.date().typeError('Invalid date. Must be in the format: DD MMM YYYY'),
    });
  }

  return (
    <div className={styles.tabContent}>
      <div className={styles.container}>
        <GridContainer className={styles.gridContainer}>
          <Grid row>
            <Grid col desktop={{ col: 8, offset: 2 }}>
              {errorMessage?.detail && (
                <div className={primeStyles.errorContainer}>
                  <Alert headingLevel="h4" type="error">
                    <span className={primeStyles.errorTitle}>{errorMessage.title}</span>
                    <span className={primeStyles.errorDetail}>{errorMessage.detail}</span>
                  </Alert>
                </div>
              )}
              <Button type="button" onClick={onCancelShipmentClick} className="usa-button usa-button-secondary">
                Cancel Shipment
              </Button>
              <Formik
                initialValues={initialValues}
                onSubmit={onSubmit}
                validationSchema={validationSchema}
                validateOnMount
              >
                {({ isValid, isSubmitting, handleSubmit }) => {
                  return (
                    <Form className={formStyles.form}>
                      {isPPM ? (
                        <PrimeUIShipmentUpdatePPMForm />
                      ) : (
                        <PrimeUIShipmentUpdateForm
                          editableWeightEstimateField={editableWeightEstimateField}
                          editableWeightActualField={editableWeightActualField}
                          editablePickupAddress={editablePickupAddress}
                          editableDestinationAddress={editableDestinationAddress}
                          estimatedWeight={initialValues.estimatedWeight}
                          actualWeight={initialValues.actualWeight}
                          requestedPickupDate={initialValues.requestedPickupDate}
                          pickupAddress={initialValues.pickupAddress}
                          destinationAddress={initialValues.destinationAddress}
                          diversion={initialValues.diversion}
                        />
                      )}
                      <div className={formStyles.formActions}>
                        <WizardNavigation
                          editMode
                          disableNext={!isValid || isSubmitting}
                          onCancelClick={handleClose}
                          onNextClick={handleSubmit}
                        />
                      </div>
                    </Form>
                  );
                }}
              </Formik>
            </Grid>
          </Grid>
        </GridContainer>
      </div>
    </div>
  );
};

PrimeUIShipmentUpdate.propTypes = {
  setFlashMessage: func,
};

PrimeUIShipmentUpdate.defaultProps = {
  setFlashMessage: () => {},
};

const mapDispatchToProps = {
  setFlashMessage: setFlashMessageAction,
};

export default connect(() => ({}), mapDispatchToProps)(PrimeUIShipmentUpdate);
