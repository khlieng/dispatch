import React, { PureComponent } from 'react';
import classnames from 'classnames';
import get from 'lodash/get';
import Button from 'components/ui/Button';
import TabListItem from 'containers/TabListItem';

export default class TabList extends PureComponent {
  handleTabClick = (server, target) => this.props.select(server, target);

  handleConnectClick = () => this.props.push('/connect');

  handleSettingsClick = () => this.props.push('/settings');

  render() {
    const {
      tab,
      channels,
      servers,
      privateChats,
      showTabList,
      openModal
    } = this.props;
    const tabs = [];

    const className = classnames('tablist', {
      'off-canvas': showTabList
    });

    channels.forEach(server => {
      const { address } = server;
      const srv = servers[address];
      tabs.push(
        <TabListItem
          key={address}
          server={address}
          content={srv.name}
          selected={tab.server === address && !tab.name}
          connected={srv.status.connected}
          onClick={this.handleTabClick}
        />
      );

      let chanLabel;
      const chanLimit =
        get(srv.features, ['CHANLIMIT', '#'], 0) || srv.features.MAXCHANNELS;
      if (chanLimit > 0) {
        chanLabel = `CHANNELS (${server.channels.length}/${chanLimit})`;
      } else {
        chanLabel = `CHANNELS (${server.channels.length})`;
      }

      tabs.push(
        <div
          key={`${address}-chans}`}
          className="tab-label"
          onClick={() => openModal('channel', { server: address })}
        >
          <span>{chanLabel}</span>
          <Button>+</Button>
        </div>
      );

      server.channels.forEach(name =>
        tabs.push(
          <TabListItem
            key={address + name}
            server={address}
            target={name}
            content={name}
            selected={tab.server === address && tab.name === name}
            onClick={this.handleTabClick}
          />
        )
      );

      if (privateChats[address] && privateChats[address].length > 0) {
        tabs.push(
          <div key={`${address}-pm}`} className="tab-label">
            <span>DIRECT MESSAGES ({privateChats[address].length})</span>
            {/*<Button>+</Button>*/}
          </div>
        );

        privateChats[address].forEach(nick =>
          tabs.push(
            <TabListItem
              key={address + nick}
              server={address}
              target={nick}
              content={nick}
              selected={tab.server === address && tab.name === nick}
              onClick={this.handleTabClick}
            />
          )
        );
      }
    });

    return (
      <div className={className}>
        <div className="tab-container">{tabs}</div>
        <div className="side-buttons">
          <Button onClick={this.handleConnectClick}>+</Button>
          <i className="icon-user" />
          <i className="icon-cog" onClick={this.handleSettingsClick} />
        </div>
      </div>
    );
  }
}
