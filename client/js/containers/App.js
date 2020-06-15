import { createStructuredSelector } from 'reselect';
import App from 'components/App';
import { getConnected } from 'state/app';
import { getSortedChannels } from 'state/channels';
import { openModal, getHasOpenModals } from 'state/modals';
import { getPrivateChats } from 'state/privateChats';
import { getNetworks } from 'state/networks';
import { getSelectedTab, select } from 'state/tab';
import { getShowTabList, hideMenu } from 'state/ui';
import connect from 'utils/connect';
import { push } from 'utils/router';

const mapState = createStructuredSelector({
  channels: getSortedChannels,
  connected: getConnected,
  privateChats: getPrivateChats,
  networks: getNetworks,
  showTabList: getShowTabList,
  tab: getSelectedTab,
  newVersionAvailable: state => state.app.newVersionAvailable,
  hasOpenModals: getHasOpenModals
});

const mapDispatch = { push, select, hideMenu, openModal };

export default connect(mapState, mapDispatch)(App);
