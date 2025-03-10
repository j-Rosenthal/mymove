import React, { Component } from 'react';
import { connect } from 'react-redux';
import { debounce, get } from 'lodash';
import { reduxForm } from 'redux-form';

import SaveCancelButtons from './SaveCancelButtons';
import profileImage from './images/profile.png';

import Alert from 'shared/Alert';
import { formatCentsRange, formatCents } from 'utils/formatters';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { fetchLatestOrders } from 'shared/Entities/modules/orders';
import { loadEntitlementsFromState } from 'shared/entitlements';
import scrollToTop from 'shared/scrollToTop';
import {
  selectCurrentOrders,
  selectCurrentPPM,
  selectPPMEstimateRange,
  selectServiceMemberFromLoggedInUser,
} from 'store/entities/selectors';
import { calculatePPMEstimate, getPPMsForMove, patchPPM, persistPPMEstimate } from 'services/internalApi';
import { updatePPM, updatePPMEstimate, updatePPMs } from 'store/entities/actions';
import { setPPMEstimateError } from 'store/onboarding/actions';
import { selectPPMEstimateError } from 'store/onboarding/selectors';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';
import EntitlementBar from 'scenes/EntitlementBar';
import './Review.css';
import './EditWeight.css';
import withRouter from 'utils/routing';

const editWeightFormName = 'edit_weight';
const weightEstimateDebounce = 300;
const examples = [
  {
    weight: 100,
    incentive: '$0 - 100',
    description: 'A few items in your car',
  },
  { weight: 400, incentive: '$300 - 400' },
  { weight: 600, incentive: '$500 - 600' },
  {
    weight: 1000,
    incentive: '$800 - 1,000',
    description: 'A trailer full of household goods',
  },
  { weight: 2000, incentive: '$1,500 - 1,800' },
  {
    weight: 5000,
    incentive: '$3,100 - 3,700',
    description: 'A moving truck',
  },
  { weight: 10000, incentive: '$5,900 - 6,800' },
];

const validateWeight = (value, formValues, props, fieldName) => {
  if (value && props.entitlement && value > props.entitlement.sum) {
    return 'Cannot be more than your full entitlement';
  }
  return undefined;
};

let EditWeightForm = (props) => {
  const {
    schema,
    handleSubmit,
    submitting,
    valid,
    entitlement,
    dirty,
    incentiveEstimateMin,
    incentiveEstimateMax,
    onWeightChange,
    initialValues,
  } = props;
  // Error class if below advance amount, otherwise warn class if incentive has changed
  let incentiveClass = '';
  let fieldClass = dirty ? 'warn' : '';
  let advanceError = false;
  const advanceAmt = get(initialValues, 'advance.requested_amount', 0);
  if (incentiveEstimateMax && advanceAmt && incentiveEstimateMax < formatCents(advanceAmt)) {
    advanceError = true;
    incentiveClass = 'error';
    fieldClass = 'error';
  } else if (get(initialValues, 'incentive_estimate_min') !== incentiveEstimateMin) {
    // Min and max are linked, so we only need to check one
    incentiveClass = 'warn';
  }

  const fullFieldClass = `weight-estimate-input ${fieldClass}`;
  return (
    <div className="grid-container usa-prose">
      <div className="grid-row">
        <div className="grid-col-12">
          <form onSubmit={handleSubmit}>
            <img src={profileImage} alt="" />
            <h1
              style={{
                display: 'inline-block',
                marginLeft: 10,
                marginBottom: 0,
                marginTop: 20,
              }}
            >
              Profile
            </h1>
            <hr />
            <h3>Edit PPM Weight:</h3>
            <p>Changes could impact your move, including the estimated PPM incentive.</p>
            <EntitlementBar entitlement={entitlement} />
            <div className="edit-weight-container">
              <div className="usa-width-one-half">
                <h4>Move estimate</h4>
                <div>
                  <SwaggerField
                    className={fullFieldClass}
                    fieldName="weight_estimate"
                    swagger={schema}
                    onChange={onWeightChange}
                    validate={validateWeight}
                    required
                  />
                  <span> lbs</span>
                </div>
                <div>
                  {!advanceError && initialValues && initialValues.incentive_estimate_min && dirty && (
                    <div className="usa-alert usa-alert--warning">
                      <div className="usa-alert__body">
                        <p className="usa-alert__text">This update will change your incentive.</p>
                      </div>
                    </div>
                  )}
                  {advanceError && (
                    <p className="advance-error">Weight is too low and will require paying back the advance.</p>
                  )}
                </div>

                <div className="display-value todo">
                  <p>Estimated Incentive</p>
                  <p className={incentiveClass}>
                    <strong>
                      {formatCentsRange(incentiveEstimateMin, incentiveEstimateMax) || 'Unable to Calculate'}
                    </strong>
                  </p>
                  {initialValues &&
                    initialValues.incentive_estimate_min &&
                    initialValues.incentive_estimate_min !== incentiveEstimateMin && (
                      <p className="subtext">
                        Originally{' '}
                        {formatCentsRange(initialValues.incentive_estimate_min, initialValues.incentive_estimate_max)}
                      </p>
                    )}
                </div>

                {get(initialValues, 'has_requested_advance') && (
                  <div className="display-value">
                    <p>Advance</p>
                    <p>
                      <strong>${formatCents(advanceAmt)}</strong>
                    </p>
                  </div>
                )}
              </div>

              <div className="usa-width-one-half">
                <h4>Examples</h4>
                <table className="examples-table">
                  <thead>
                    <tr>
                      <th>Weight</th>
                      <th>Incentive</th>
                      <th />
                    </tr>
                  </thead>
                  <tbody>
                    {examples.map((ex) => (
                      <tr key={ex.weight}>
                        <td>{ex.weight.toLocaleString()}</td>
                        <td>{ex.incentive}</td>
                        <td>{ex.description || ''}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
            <SaveCancelButtons valid={valid} submitting={submitting} />
          </form>
        </div>
      </div>
    </div>
  );
};
EditWeightForm = reduxForm({
  form: editWeightFormName,
  enableReinitialize: true,
})(EditWeightForm);

class EditWeight extends Component {
  constructor(props) {
    super(props);
    this.state = { currentPPM: {} };
  }

  componentDidMount() {
    const {
      router: { params: moveId },
    } = this.props;
    getPPMsForMove(moveId).then((response) => this.props.updatePPMs(response));
    this.props.fetchLatestOrders(this.props.serviceMemberId);
    const { currentPPM, originDutyLocationZip, orders } = this.props;
    this.handleWeightChange(
      currentPPM.original_move_date,
      currentPPM.pickup_postal_code,
      originDutyLocationZip,
      orders.id,
      currentPPM.weight_estimate,
    );
    scrollToTop();
  }

  handleWeightChange = (moveDate, originZip, originDutyLocationZip, ordersId, weightEstimate) => {
    calculatePPMEstimate(moveDate, originZip, originDutyLocationZip, ordersId, weightEstimate)
      .then((response) => {
        this.props.updatePPMEstimate(response);
        this.props.setPPMEstimateError(null);
      })
      .catch((error) => {
        this.props.setPPMEstimateError(error);
      });
  };

  debouncedHandleWeightChange = debounce(this.handleWeightChange, weightEstimateDebounce);

  onWeightChange = (e, newValue) => {
    const { currentPPM, entitlement, originDutyLocationZip, orders } = this.props;
    if (newValue > 0 && newValue <= entitlement.sum) {
      this.debouncedHandleWeightChange(
        currentPPM.original_move_date,
        currentPPM.pickup_postal_code,
        originDutyLocationZip,
        orders.id,
        newValue,
      );
    } else {
      this.debouncedHandleWeightChange.cancel();
    }
  };

  updatePpm = (values, dispatch, props) => {
    const { setFlashMessage, router } = this.props;
    const { params: moveId, navigate } = router;
    return patchPPM(moveId, {
      id: this.props.currentPPM.id,
      weight_estimate: values.weight_estimate,
    })
      .then((response) => {
        this.props.updatePPM(response);
        return response;
      })
      .then((response) => persistPPMEstimate(moveId, response.id))
      .then((response) => this.props.updatePPM(response))
      .then(() => {
        setFlashMessage('EDIT_PPM_WEIGHT_SUCCESS', 'success', '', 'Your changes have been saved.');

        navigate(-1);
      })
      .catch(() => {
        scrollToTop();
      });
  };

  chooseEstimateErrorText(hasEstimateError, rateEngineError) {
    if (rateEngineError) {
      return (
        <div className="grid-row">
          <div className="grid-col-12 error-message">
            <Alert type="warning" heading="Could not retrieve estimate">
              MilMove does not presently support short-haul PPM moves. Please contact your PPPO.
            </Alert>
          </div>
        </div>
      );
    }

    if (hasEstimateError) {
      return (
        <div className="grid-row">
          <div className="grid-col-12 error-message">
            <Alert type="warning" heading="Could not retrieve estimate">
              There was an issue retrieving an estimate for your incentive. You still qualify but may need to talk with
              your local PPPO.
            </Alert>
          </div>
        </div>
      );
    }
    // don't show any text if there is no error
    return undefined;
  }

  render() {
    const {
      error,
      schema,
      entitlement,
      hasEstimateError,
      rateEngineError,
      currentPPM,
      incentiveEstimateMin,
      incentiveEstimateMax,
    } = this.props;

    return (
      <div className="grid-container usa-prose">
        {error && (
          <div className="grid-row">
            <div className="grid-col-12 error-message">
              <Alert type="error" heading="An error occurred">
                {error.message}
              </Alert>
            </div>
          </div>
        )}

        <div className="grid-container usa-prose">
          <div className="grid-row">
            <div className="grid-col-12">{this.chooseEstimateErrorText(hasEstimateError, rateEngineError)}</div>
          </div>
        </div>

        <div className="grid-row">
          <div className="grid-col-12">
            <EditWeightForm
              initialValues={currentPPM}
              incentiveEstimateMin={incentiveEstimateMin}
              incentiveEstimateMax={incentiveEstimateMax}
              onSubmit={this.updatePpm}
              onWeightChange={this.onWeightChange}
              entitlement={entitlement}
              schema={schema}
            />
          </div>
        </div>
      </div>
    );
  }
}

function mapStateToProps(state) {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const serviceMemberId = serviceMember?.id;

  return {
    serviceMemberId,
    currentPPM: selectCurrentPPM(state) || {},
    incentiveEstimateMin: selectPPMEstimateRange(state)?.range_min,
    incentiveEstimateMax: selectPPMEstimateRange(state)?.range_max,
    entitlement: loadEntitlementsFromState(state),
    schema: get(state, 'swaggerInternal.spec.definitions.UpdatePersonallyProcuredMovePayload', {}),
    originDutyLocationZip: serviceMember?.current_location?.address?.postalCode,
    orders: selectCurrentOrders(state) || {},
    rateEngineError: selectPPMEstimateError(state),
  };
}

const mapDispatchToProps = {
  fetchLatestOrders,
  updatePPM,
  updatePPMs,
  updatePPMEstimate,
  setPPMEstimateError,
  setFlashMessage: setFlashMessageAction,
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(EditWeight));
