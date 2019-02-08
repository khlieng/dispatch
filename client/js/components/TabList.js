import React, { PureComponent } from 'react';
import classnames from 'classnames';
import get from 'lodash/get';
import Button from 'components/ui/Button';
import TabListItem from 'containers/TabListItem';
import { count } from 'utils';

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

      const chanCount = count(server.channels, c => c.joined);
      const chanLimit =
        get(srv.features, ['CHANLIMIT', '#'], 0) || srv.features.MAXCHANNELS;

      let chanLabel;
      if (chanLimit > 0) {
        chanLabel = (
          <span>
            <span className="success">{chanCount}</span>/{chanLimit}
          </span>
        );
      } else if (chanCount > 0) {
        chanLabel = <span className="success">{chanCount}</span>;
      }

      tabs.push(
        <div
          key={`${address}-chans}`}
          className="tab-label"
          onClick={() => openModal('channel', { server: address })}
        >
          <span>CHANNELS {chanLabel}</span>
          <Button title="Join Channel">+</Button>
        </div>
      );

      server.channels.forEach(({ name, joined }) =>
        tabs.push(
          <TabListItem
            key={address + name}
            server={address}
            target={name}
            content={name}
            joined={joined}
            selected={tab.server === address && tab.name === name}
            onClick={this.handleTabClick}
          />
        )
      );

      if (privateChats[address] && privateChats[address].length > 0) {
        tabs.push(
          <div key={`${address}-pm}`} className="tab-label">
            <span>
              DIRECT MESSAGES{' '}
              <span style={{ color: '#6bb758' }}>
                {privateChats[address].length}
              </span>
            </span>
            {/* <Button>+</Button> */}
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
