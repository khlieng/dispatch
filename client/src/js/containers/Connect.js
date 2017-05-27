import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';
import Connect from '../components/pages/Connect';
import { join } from '../state/channels';
import { getConnectDefaults } from '../state/environment';
import { connect as connectServer } from '../state/servers';
import { select } from '../state/tab';

const mapState = createStructuredSelector({
  defaults: getConnectDefaults
});

const mapDispatch = {
  join,
  connect: connectServer,
  select
};

export default connect(mapState, mapDispatch)(Connect);
