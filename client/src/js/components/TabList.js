import React, { Component } from 'react';
import pure from 'pure-render-decorator';
import TabListItem from './TabListItem';

@pure
export default class TabList extends Component {
  handleTabClick = (server, target) => {
    this.props.select(server, target, target && target.charAt(0) !== '#');
    this.props.hideMenu();
  }

  handleConnectClick = () => {
    this.props.pushPath('/connect');
    this.props.hideMenu();
  }

  handleSettingsClick = () => {
    this.props.pushPath('/settings');
    this.props.hideMenu();
  }

  render() {
    const { channels, servers, privateChats, showMenu, selected } = this.props;
    const className = showMenu ? 'tablist off-canvas' : 'tablist';
    const tabs = [];

    channels.forEach((server, address) => {
      tabs.push(
        <TabListItem
          key={address}
          server={address}
          content={servers.getIn([address, 'name'])}
          selected={
            selected.server === address &&
            selected.channel === null &&
            selected.user === null
          }
          onClick={this.handleTabClick}
        />
      );

      server.forEach((channel, name) => tabs.push(
        <TabListItem
          key={address + name}
          server={address}
          target={name}
          content={name}
          selected={selected.server === address && selected.channel === name}
          onClick={this.handleTabClick}
        />
      ));

      if (privateChats.has(address)) {
        privateChats.get(address).forEach(nick => tabs.push(
          <TabListItem
            key={address + nick}
            server={address}
            target={nick}
            content={nick}
            selected={selected.server === address && selected.user === nick}
            onClick={this.handleTabClick}
          />
        ));
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
