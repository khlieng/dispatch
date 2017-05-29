import React, { Component } from 'react';
import Route from '../containers/Route';
import Chat from '../containers/Chat';
import Connect from '../containers/Connect';
import Settings from '../containers/Settings';
import TabList from '../components/TabList';

export default class App extends Component {
  handleClick = () => {
    const { showTabList, hideMenu } = this.props;
    if (showTabList) {
      hideMenu();
    }
  };

  render() {
    const mainClass = this.props.showTabList ? 'main-container off-canvas' : 'main-container';
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
