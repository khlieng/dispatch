import React, { PureComponent } from 'react';
import { connect } from 'react-redux';
import { push } from '../util/router';
import Route from './Route';
import Chat from './Chat';
import Connect from './Connect';
import Settings from './Settings';
import TabList from '../components/TabList';
import { select } from '../actions/tab';
import { hideMenu } from '../actions/ui';

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

function mapStateToProps(state) {
  return {
    servers: state.servers,
    channels: state.channels,
    privateChats: state.privateChats,
    showTabList: state.ui.showTabList,
    tab: state.tab.selected
  };
}

export default connect(mapStateToProps, { pushPath: push, select, hideMenu })(App);
