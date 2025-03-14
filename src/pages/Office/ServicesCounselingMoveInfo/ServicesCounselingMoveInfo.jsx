import React, { Suspense, lazy, useState, useEffect } from 'react';
import { Routes, useParams, Route, Navigate, useLocation, matchPath } from 'react-router-dom';
import { useSelector } from 'react-redux';

import 'styles/office.scss';

import ServicesCounselorTabNav from 'components/Office/ServicesCounselingTabNav/ServicesCounselingTabNav';
import CustomerHeader from 'components/CustomerHeader';
import SystemError from 'components/SystemError';
import { servicesCounselingRoutes } from 'constants/routes';
import { useTXOMoveInfoQueries } from 'hooks/queries';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';

const ServicesCounselingMoveDocumentWrapper = lazy(() =>
  import('pages/Office/ServicesCounselingMoveDocumentWrapper/ServicesCounselingMoveDocumentWrapper'),
);
const ServicesCounselingMoveDetails = lazy(() =>
  import('pages/Office/ServicesCounselingMoveDetails/ServicesCounselingMoveDetails'),
);
const ServicesCounselingAddShipment = lazy(() =>
  import('pages/Office/ServicesCounselingAddShipment/ServicesCounselingAddShipment'),
);
const ServicesCounselingEditShipmentDetails = lazy(() =>
  import('pages/Office/ServicesCounselingEditShipmentDetails/ServicesCounselingEditShipmentDetails'),
);
const CustomerInfo = lazy(() => import('pages/Office/CustomerInfo/CustomerInfo'));
const ServicesCounselorCustomerSupportRemarks = lazy(() =>
  import('pages/Office/ServicesCounselorCustomerSupportRemarks/ServicesCounselorCustomerSupportRemarks'),
);
const MoveHistory = lazy(() => import('pages/Office/MoveHistory/MoveHistory'));
const ReviewDocuments = lazy(() => import('pages/Office/PPM/ReviewDocuments/ReviewDocuments'));
const ServicesCounselingReviewShipmentWeights = lazy(() =>
  import('pages/Office/ServicesCounselingReviewShipmentWeights/ServicesCounselingReviewShipmentWeights'),
);

const ServicesCounselingMoveInfo = () => {
  const [unapprovedShipmentCount, setUnapprovedShipmentCount] = useState(0);

  const [infoSavedAlert, setInfoSavedAlert] = useState(null);
  const { hasRecentError, traceId } = useSelector((state) => state.interceptor);
  const onInfoSavedUpdate = (alertType) => {
    if (alertType === 'error') {
      setInfoSavedAlert({
        alertType,
        message: 'Something went wrong, and your changes were not saved. Please try again later.',
      });
    } else {
      setInfoSavedAlert({
        alertType,
        message: 'Your changes were saved.',
      });
    }
  };

  // Clear the alert when route changes
  const location = useLocation();
  useEffect(() => {
    if (
      infoSavedAlert &&
      !matchPath(
        {
          path: servicesCounselingRoutes.BASE_MOVE_VIEW_PATH,
          end: true,
        },
        location.pathname,
      )
    ) {
      setInfoSavedAlert(null);
    }
  }, [infoSavedAlert, location]);

  const { moveCode } = useParams();
  const { order, customerData, isLoading, isError } = useTXOMoveInfoQueries(moveCode);

  const { pathname } = useLocation();
  const hideNav =
    matchPath(
      {
        path: servicesCounselingRoutes.BASE_SHIPMENT_ADD_PATH,
        end: true,
      },
      pathname,
    ) ||
    matchPath(
      {
        path: servicesCounselingRoutes.BASE_SHIPMENT_EDIT_PATH,
        end: true,
      },
      pathname,
    ) ||
    matchPath(
      {
        path: servicesCounselingRoutes.BASE_ORDERS_EDIT_PATH,
        end: true,
      },
      pathname,
    ) ||
    matchPath(
      {
        path: servicesCounselingRoutes.BASE_ALLOWANCES_EDIT_PATH,
        end: true,
      },
      pathname,
    ) ||
    matchPath(
      {
        path: servicesCounselingRoutes.BASE_SHIPMENT_REVIEW_PATH,
        end: true,
      },
      pathname,
    ) ||
    matchPath(
      {
        path: servicesCounselingRoutes.BASE_CUSTOMER_INFO_EDIT_PATH,
        end: true,
      },
      pathname,
    );

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  return (
    <>
      <CustomerHeader order={order} customer={customerData} moveCode={moveCode} />
      {hasRecentError && (
        <SystemError>
          Something isn&apos;t working, but we&apos;re not sure what. Wait a minute and try again.
          <br />
          If that doesn&apos;t fix it, contact the{' '}
          <a href="https://move.mil/customer-service#technical-help-desk">Technical Help Desk</a> and give them this
          code: <strong>{traceId}</strong>
        </SystemError>
      )}

      {!hideNav && <ServicesCounselorTabNav unapprovedShipmentCount={unapprovedShipmentCount} moveCode={moveCode} />}

      <Suspense fallback={<LoadingPlaceholder />}>
        <Routes>
          {/* TODO - Routes not finalized, revisit */}
          <Route path={servicesCounselingRoutes.SHIPMENT_REVIEW_PATH} end element={<ReviewDocuments />} />

          <Route
            path={servicesCounselingRoutes.MOVE_VIEW_PATH}
            end
            element={
              <ServicesCounselingMoveDetails
                infoSavedAlert={infoSavedAlert}
                setUnapprovedShipmentCount={setUnapprovedShipmentCount}
              />
            }
          />

          <Route
            key="servicesCounselingAddShipment"
            end
            path={servicesCounselingRoutes.SHIPMENT_ADD_PATH}
            element={<ServicesCounselingAddShipment />}
          />

          <Route
            path={servicesCounselingRoutes.CUSTOMER_SUPPORT_REMARKS_PATH}
            end
            element={<ServicesCounselorCustomerSupportRemarks />}
          />

          <Route path={servicesCounselingRoutes.MOVE_HISTORY_PATH} end element={<MoveHistory moveCode={moveCode} />} />

          <Route
            path={servicesCounselingRoutes.ALLOWANCES_EDIT_PATH}
            end
            element={<ServicesCounselingMoveDocumentWrapper />}
          />
          <Route
            path={servicesCounselingRoutes.ORDERS_EDIT_PATH}
            end
            element={<ServicesCounselingMoveDocumentWrapper />}
          />

          <Route
            path={servicesCounselingRoutes.CUSTOMER_INFO_EDIT_PATH}
            end
            element={
              <CustomerInfo
                ordersId={order.id}
                customer={customerData}
                isLoading={isLoading}
                isError={isError}
                onUpdate={onInfoSavedUpdate}
              />
            }
          />

          <Route
            path={servicesCounselingRoutes.SHIPMENT_EDIT_PATH}
            end
            element={<ServicesCounselingEditShipmentDetails onUpdate={onInfoSavedUpdate} />}
          />

          <Route
            path={servicesCounselingRoutes.SHIPMENT_ADVANCE_PATH}
            end
            element={<ServicesCounselingEditShipmentDetails onUpdate={onInfoSavedUpdate} isAdvancePage />}
          />

          <Route
            path={servicesCounselingRoutes.SHIPMENT_REVIEW_PATH}
            exact
            element={<ReviewDocuments onUpdate={onInfoSavedUpdate} />}
          />

          <Route
            path={servicesCounselingRoutes.REVIEW_SHIPMENT_WEIGHTS_PATH}
            exact
            element={<ServicesCounselingReviewShipmentWeights moveCode={moveCode} />}
          />

          {/* TODO - clarify role/tab access */}
          <Route path="/" element={<Navigate to={servicesCounselingRoutes.MOVE_VIEW_PATH} replace />} />
        </Routes>
      </Suspense>
    </>
  );
};

export default ServicesCounselingMoveInfo;
