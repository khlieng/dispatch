import React, { Component } from 'react';
import pure from 'pure-render-decorator';
import TabListItem from './TabListItem';

@pure
export default class TabList extends Component {
  handleConnectClick = () => {
    this.props.pushPath('/connect');
    this.props.hideMenu();
  }

  handleSettingsClick = () => {
    this.props.pushPath('/settings');
    this.props.hideMenu();
  }

  render() {
    const { channels, servers, privateChats, showMenu, select, selected } = this.props;
    const className = showMenu ? 'tablist off-canvas' : 'tablist';
    const tabs = [];

    channels.forEach((server, address) => {
      tabs.push(
        <TabListItem
          key={address}
          server
          content={servers.getIn([address, 'name'])}
          selected={selected.server === address && selected.channel === null && selected.user === null}
          onClick={() => select(address)}
        />
      );

      server.forEach((channel, name) => {
        tabs.push(
          <TabListItem
            key={address + channel.get('name')}
            content={channel.get('name')}
            selected={selected.server === address && selected.channel === name}
            onClick={() => select(address, channel.get('name'))}
          />
        );
      });

      if (privateChats.has(address)) {
        privateChats.get(address).forEach(nick => {
          tabs.push(
            <TabListItem
              key={address + nick}
              content={nick}
              selected={selected.server === address && selected.user === nick}
              onClick={() => select(address, nick, true)}
            />
          );
        });
      }
    });

    return (
      <div className={className}>
        <button className="button-connect" onClick={this.handleConnectClick}>Connect</button>
        <div className="tab-container">{tabs}</div>
        <div className="side-buttons">
          <i className="icon-user"></i>
          <i className="icon-cog" onClick={this.handleSettingsClick}></i>
        </div>
      </div>
    );
  }
}
