import React, { Component } from 'react';
import Route from 'containers/Route';
import Chat from 'containers/Chat';
import Connect from 'containers/Connect';
import Settings from 'containers/Settings';
import TabList from 'components/TabList';
import classnames from 'classnames';

export default class App extends Component {
  handleClick = () => {
    const { showTabList, hideMenu } = this.props;
    if (showTabList) {
      hideMenu();
    }
  };

  render() {
    const {
      connected,
      tab,
      channels,
      servers,
      privateChats,
      showTabList,
      select,
      push
    } = this.props;

    const mainClass = classnames('main-container', {
      'off-canvas': showTabList
    });

    return (
      <div className="wrap">
        {!connected && (
          <div className="app-info">
            Connection lost, attempting to reconnect...
          </div>
        )}
        <div className="app-container" onClick={this.handleClick}>
          <TabList
            tab={tab}
            channels={channels}
            servers={servers}
            privateChats={privateChats}
            showTabList={showTabList}
            select={select}
            push={push}
          />
          <div className={mainClass}>
            <Route name="chat">
              <Chat />
            </Route>
            <Route name="connect">
              <Connect />
            </Route>
            <Route name="settings">
              <Settings />
            </Route>
          </div>
        </div>
      </div>
    );
  }
}
