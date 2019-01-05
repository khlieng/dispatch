import React, { memo } from 'react';
import Confirm from 'components/modals/Confirm';

const Modals = () => (
  <>
    <Confirm />
  </>
);

export default memo(Modals, () => true);
