import React, { PureComponent } from 'react';
import { connect } from 'react-redux';
import { push } from 'react-router-redux';
import TabList from '../components/TabList';
import { select } from '../actions/tab';
import { hideMenu } from '../actions/ui';

class App extends PureComponent {
  render() {
    const { showTabList, children } = this.props;
    const mainClass = showTabList ? 'main-container off-canvas' : 'main-container';
    return (
      <div>
        <TabList {...this.props} />
        <div className={mainClass}>
          {children}
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
    selected: state.tab.selected
  };
}

export default connect(mapStateToProps, { pushPath: push, select, hideMenu })(App);
