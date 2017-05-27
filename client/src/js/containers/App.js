import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';
import App from '../components/App';
import { getChannels } from '../state/channels';
import { getPrivateChats } from '../state/privateChats';
import { getServers } from '../state/servers';
import { getSelectedTab, select } from '../state/tab';
import { getShowTabList, hideMenu } from '../state/ui';
import { push } from '../util/router';

const mapState = createStructuredSelector({
  channels: getChannels,
  privateChats: getPrivateChats,
  servers: getServers,
  showTabList: getShowTabList,
  tab: getSelectedTab
});

const mapDispatch = (dispatch, props) => ({
  onClick: () => {
    if (props.showTabList) {
      dispatch(hideMenu());
    }
  },
  ...bindActionCreators({ push, select }, dispatch)
});

export default connect(mapState, mapDispatch)(App);
