import { sswIsDisabled } from './PaymentsPanel';

import { SIGNED_CERT_OPTIONS } from 'shared/constants';

describe('Download Shipment Summary button', () => {
  describe('PPM only move', () => {
    it('is disabled when no ppm', () => {
      const ppm = null;
      const signedCertification = { certification_type: SIGNED_CERT_OPTIONS.PPM_PAYMENT };

      expect(sswIsDisabled(ppm, signedCertification)).toEqual(true);
    });

    it('is disabled when missing net weight', () => {
      const ppm = { actual_move_date: '2018-11-11' };
      const signedCertification = { certification_type: SIGNED_CERT_OPTIONS.PPM_PAYMENT };

      expect(sswIsDisabled(ppm, signedCertification)).toEqual(true);
    });

    it('is disabled when missing actual move date', () => {
      const ppm = { net_weight: 8000 };
      const signedCertification = { certification_type: SIGNED_CERT_OPTIONS.PPM_PAYMENT };

      expect(sswIsDisabled(ppm, signedCertification)).toEqual(true);
    });

    it('is disabled when missing signature', () => {
      const ppm = { net_weight: 8000, actual_move_date: '2018-11-11' };
      const signedCertification = null;

      expect(sswIsDisabled(ppm, signedCertification)).toEqual(true);
    });

    it('is enabled when has signature, actual move date, net weight', () => {
      const ppm = { net_weight: 8000, actual_move_date: '2018-11-11' };
      const signedCertification = { certification_type: SIGNED_CERT_OPTIONS.PPM_PAYMENT };

      expect(sswIsDisabled(ppm, signedCertification)).toEqual(false);
    });
  });
});
