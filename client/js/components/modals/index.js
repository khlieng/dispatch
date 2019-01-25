import React, { memo } from 'react';
import AddChannel from 'components/modals/AddChannel';
import Confirm from 'components/modals/Confirm';
import Topic from 'components/modals/Topic';

const Modals = () => (
  <>
    <AddChannel />
    <Confirm />
    <Topic />
  </>
);

export default memo(Modals, () => true);
