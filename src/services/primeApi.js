import Swagger from 'swagger-client';

import { makeSwaggerRequest, requestInterceptor, responseInterceptor } from './swaggerRequest';

let primeSimulatorClient = null;

// setting up the same config from Swagger/api.js
export async function getPrimeSimulatorClient() {
  if (!primeSimulatorClient) {
    primeSimulatorClient = await Swagger({
      url: '/prime/v1/swagger.yaml',
      requestInterceptor,
      responseInterceptor,
    });
  }
  return primeSimulatorClient;
}

export async function makePrimeSimulatorRequest(operationPath, params = {}, options = {}) {
  const client = await getPrimeSimulatorClient();
  return makeSwaggerRequest(client, operationPath, params, options);
}

export async function getPrimeSimulatorAvailableMoves() {
  const operationPath = 'moveTaskOrder.listMoves';
  return makePrimeSimulatorRequest(operationPath, {}, { schemaKey: 'listMoves', normalize: false });
}

export async function getPrimeSimulatorMove(key, locator) {
  return makePrimeSimulatorRequest('moveTaskOrder.getMoveTaskOrder', { moveID: locator }, { normalize: false });
}

export async function createPaymentRequest({ moveTaskOrderID, serviceItems }) {
  return makePrimeSimulatorRequest(
    'paymentRequest.createPaymentRequest',
    { body: { moveTaskOrderID, serviceItems } },
    { normalize: false },
  );
}

export async function completeCounseling({ moveTaskOrderID, ifMatchETag }) {
  return makePrimeSimulatorRequest(
    'moveTaskOrder.updateMTOPostCounselingInformation',
    { moveTaskOrderID, 'If-Match': ifMatchETag },
    { normalize: false },
  );
}

export async function deleteShipment({ mtoShipmentID }) {
  return makePrimeSimulatorRequest('mtoShipment.deleteMTOShipment', { mtoShipmentID }, { normalize: false });
}

export async function createUpload({ paymentRequestID, file }) {
  return makePrimeSimulatorRequest('paymentRequest.createUpload', { paymentRequestID, file }, { normalize: false });
}

export async function createServiceRequestDocumentUpload({ mtoServiceItemID, file }) {
  return makePrimeSimulatorRequest(
    'mtoServiceItem.createServiceRequestDocumentUpload',
    { mtoServiceItemID, file },
    { normalize: false },
  );
}

export function createPrimeMTOShipment({ normalize = false, schemaKey = 'mtoShipment', body }) {
  const operationPath = 'mtoShipment.createMTOShipment';
  return makePrimeSimulatorRequest(
    operationPath,
    {
      body,
    },
    { schemaKey, normalize },
  );
}

export function updatePrimeMTOShipment({
  mtoShipmentID,
  ifMatchETag,
  normalize = true,
  schemaKey = 'mtoShipment',
  body,
}) {
  const operationPath = 'mtoShipment.updateMTOShipment';
  return makePrimeSimulatorRequest(
    operationPath,
    {
      mtoShipmentID,
      'If-Match': ifMatchETag,
      body,
    },
    { schemaKey, normalize },
  );
}

export function createServiceItem({ body }) {
  return makePrimeSimulatorRequest('mtoServiceItem.createMTOServiceItem', { body: { ...body } }, { normalize: false });
}

export function updatePrimeMTOShipmentAddress({
  mtoShipmentID,
  ifMatchETag,
  addressID,
  normalize = false,
  schemaKey = 'mtoShipment',
  body,
}) {
  const operationPath = 'mtoShipment.updateMTOShipmentAddress';
  return makePrimeSimulatorRequest(
    operationPath,
    {
      mtoShipmentID,
      addressID,
      'If-Match': ifMatchETag,
      body,
    },
    { schemaKey, normalize },
  );
}

export function updatePrimeMTOShipmentReweigh({
  mtoShipmentID,
  reweighID,
  ifMatchETag,
  normalize = false,
  schemaKey = 'mtoShipment',
  body,
}) {
  const operationPath = 'mtoShipment.updateReweigh';
  return makePrimeSimulatorRequest(
    operationPath,
    {
      mtoShipmentID,
      reweighID,
      'If-Match': ifMatchETag,
      body,
    },
    { schemaKey, normalize },
  );
}

// updatePrimeMTOShipmentStatus This function is used to by the Prime Simulator
// to send a cancellation request to the Prime API.
export function updatePrimeMTOShipmentStatus({
  mtoShipmentID,
  ifMatchETag,
  normalize = false,
  schemaKey = 'mtoShipment',
}) {
  const operationPath = 'mtoShipment.updateMTOShipmentStatus';
  // Default body is defined here as we can only send a status of CANCELED at
  // this time. See documentation here:
  // https://transcom.github.io/mymove-docs/api/prime#tag/mtoShipment/operation/updateMTOShipmentStatus
  const body = {
    status: 'CANCELED',
  };
  return makePrimeSimulatorRequest(
    operationPath,
    {
      mtoShipmentID,
      'If-Match': ifMatchETag,
      body,
    },
    { schemaKey, normalize },
  );
}
