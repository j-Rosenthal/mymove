import React from 'react';
import { get } from 'lodash';
import PropTypes from 'prop-types';

import ReviewSection from '../../components/Customer/ReviewSection';
import SectionWrapper from '../../components/Customer/SectionWrapper';

import Address from './Address';

import { formatDateSM } from 'utils/formatters';
import { getFullSMName } from 'utils/moveSetupFlow';

import './Review.css';

function ServiceMemberSummary(props) {
  const {
    orders,
    serviceMember,
    schemaRank,
    schemaAffiliation,
    schemaOrdersType,
    moveIsApproved,
    editOrdersPath,
    uploads,
  } = props;

  const editProfilePath = 'edit-profile';

  const yesNoMap = { true: 'Yes', false: 'No' };

  const currentResidentialAddress = <Address address={get(serviceMember, 'residential_address')} />;

  const serviceMemberData = [
    { label: 'Name', value: getFullSMName(serviceMember) },
    { label: 'Branch', value: get(schemaAffiliation['x-display-value'], get(serviceMember, 'affiliation')) },
    { label: 'Rank', value: get(schemaRank['x-display-value'], get(serviceMember, 'rank')) },
    { label: 'DoD ID#', value: get(serviceMember, 'edipi') },
    { label: 'Current duty location', value: get(serviceMember, 'current_location.name') },
    { label: 'Contact info' },
    { label: 'Best contact phone', value: get(serviceMember, 'telephone') },
    { label: 'Personal email', value: get(serviceMember, 'personal_email') },
    { label: 'Current pickup address', value: currentResidentialAddress },
  ];

  const ordersData = [
    { label: 'Orders type', value: get(schemaOrdersType['x-display-value'], get(orders, 'orders_type')) },
    { label: 'Orders date', value: formatDateSM(get(orders, 'issue_date')) },
    { label: 'Report by date', value: formatDateSM(get(orders, 'report_by_date')) },
    { label: 'New duty location', value: get(orders, 'new_duty_location.name') },
    { label: 'Dependents', value: orders && yesNoMap[get(orders, 'has_dependents', '').toString()] },
    { label: 'Orders', value: uploads && uploads.length },
  ];

  return (
    <div className="service-member-summary">
      <SectionWrapper>
        <ReviewSection fieldData={serviceMemberData} title="Profile" editLink={editProfilePath} />
      </SectionWrapper>
      <div>
        {moveIsApproved && '*'}
        {!moveIsApproved && (
          <SectionWrapper>
            <ReviewSection fieldData={ordersData} title="Orders" editLink={editOrdersPath} />
          </SectionWrapper>
        )}
      </div>
    </div>
  );
}

ServiceMemberSummary.propTypes = {
  serviceMember: PropTypes.object,
  schemaRank: PropTypes.object.isRequired,
  schemaAffiliation: PropTypes.object.isRequired,
  schemaOrdersType: PropTypes.object.isRequired,
  orders: PropTypes.object,
  moveIsApproved: PropTypes.bool,
  editOrdersPath: PropTypes.string,
};

export default ServiceMemberSummary;
