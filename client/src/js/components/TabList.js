import React, { PureComponent } from 'react';
import TabListItem from './TabListItem';

export default class TabList extends PureComponent {
  handleTabClick = (server, target) => this.props.select(server, target);
  handleConnectClick = () => this.props.push('/connect');
  handleSettingsClick = () => this.props.push('/settings');

  render() {
    const { tab, channels, servers, privateChats, showTabList } = this.props;
    const className = showTabList ? 'tablist off-canvas' : 'tablist';
    const tabs = [];

    channels.forEach((server, address) => {
      tabs.push(
        <TabListItem
          key={address}
          server={address}
          content={servers.getIn([address, 'name'])}
          selected={tab.server === address && tab.name === null}
          connected={servers.getIn([address, 'connected'])}
          onClick={this.handleTabClick}
        />
      );

      server.forEach((channel, name) => tabs.push(
        <TabListItem
          key={address + name}
          server={address}
          target={name}
          content={name}
          selected={tab.server === address && tab.name === name}
          onClick={this.handleTabClick}
        />
      ));

      if (privateChats.has(address) && privateChats.get(address).size > 0) {
        tabs.push(<div key={`${address}-pm}`} className="tab-label">Private messages</div>);

        privateChats.get(address).forEach(nick => tabs.push(
          <TabListItem
            key={address + nick}
            server={address}
            target={nick}
            content={nick}
            selected={tab.server === address && tab.name === nick}
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
          <i className="icon-user" />
          <i className="icon-cog" onClick={this.handleSettingsClick} />
        </div>
      </div>
    );
  }
}
