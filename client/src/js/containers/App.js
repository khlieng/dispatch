import React, { PureComponent } from 'react';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';
import { push } from '../util/router';
import Route from './Route';
import Chat from './Chat';
import Connect from './Connect';
import Settings from './Settings';
import TabList from '../components/TabList';
import { getChannels } from '../state/channels';
import { getPrivateChats } from '../state/privateChats';
import { getServers } from '../state/servers';
import { getSelectedTab, select } from '../state/tab';
import { getShowTabList, hideMenu } from '../state/ui';

class App extends PureComponent {
  handleClick = () => {
    if (this.props.showTabList) {
      this.props.hideMenu();
    }
  };

  render() {
    const { showTabList } = this.props;
    const mainClass = showTabList ? 'main-container off-canvas' : 'main-container';

    return (
      <div onClick={this.handleClick}>
        <TabList {...this.props} />
        <div className={mainClass}>
          <Route name="chat"><Chat /></Route>
          <Route name="connect"><Connect /></Route>
          <Route name="settings"><Settings /></Route>
        </div>
      </div>
    );
  }
}

const mapState = createStructuredSelector({
  channels: getChannels,
  privateChats: getPrivateChats,
  servers: getServers,
  showTabList: getShowTabList,
  tab: getSelectedTab
});

export default connect(mapState, { pushPath: push, select, hideMenu })(App);
