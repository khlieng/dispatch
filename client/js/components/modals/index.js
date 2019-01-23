import React, { memo } from 'react';
import AddChannel from 'components/modals/AddChannel';
import Confirm from 'components/modals/Confirm';

const Modals = () => (
  <>
    <AddChannel />
    <Confirm />
  </>
);

export default memo(Modals, () => true);
