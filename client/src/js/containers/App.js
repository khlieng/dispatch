import React, { Component } from 'react';
import { connect } from 'react-redux';
import { pushPath } from 'redux-simple-router';
import pure from 'pure-render-decorator';
import TabList from '../components/TabList';
import * as actions from '../actions/tab';

@pure
class App extends Component {
  render() {
    const { showMenu, children } = this.props;
    const mainClass = showMenu ? 'main-container off-canvas' : 'main-container';
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
    showMenu: state.showMenu,
    selected: state.tab.selected
  };
}

export default connect(mapStateToProps, { pushPath, ...actions })(App);
