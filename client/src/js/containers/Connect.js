import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';
import Connect from 'components/pages/Connect';
import { getConnectDefaults, getApp } from 'state/app';
import { join } from 'state/channels';
import { connect as connectServer } from 'state/servers';
import { select } from 'state/tab';

const mapState = createStructuredSelector({
  defaults: getConnectDefaults,
  hexIP: state => getApp(state).hexIP
});

const mapDispatch = {
  join,
  connect: connectServer,
  select
};

export default connect(
  mapState,
  mapDispatch
)(Connect);
