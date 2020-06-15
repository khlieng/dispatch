import React, { PureComponent } from 'react';
import classnames from 'classnames';
import get from 'lodash/get';
import { FiPlus, FiUser, FiSettings } from 'react-icons/fi';
import Button from 'components/ui/Button';
import TabListItem from 'containers/TabListItem';
import { count } from 'utils';

export default class TabList extends PureComponent {
  handleTabClick = (network, target) => this.props.select(network, target);

  handleConnectClick = () => this.props.push('/connect');

  handleSettingsClick = () => this.props.push('/settings');

  render() {
    const {
      tab,
      channels,
      networks,
      privateChats,
      showTabList,
      openModal
    } = this.props;
    const tabs = [];

    const className = classnames('tablist', {
      'off-canvas': showTabList
    });

    channels.forEach(network => {
      const { address } = network;
      const srv = networks[address];
      tabs.push(
        <TabListItem
          key={address}
          network={address}
          content={srv.name}
          selected={tab.network === address && !tab.name}
          connected={srv.connected}
          onClick={this.handleTabClick}
        />
      );

      const chanCount = count(network.channels, c => c.joined);
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
          onClick={() => openModal('channel', address)}
        >
          <span>CHANNELS {chanLabel}</span>
          <Button title="Join Channel">+</Button>
        </div>
      );

      network.channels.forEach(({ name, joined }) =>
        tabs.push(
          <TabListItem
            key={address + name}
            network={address}
            target={name}
            content={name}
            joined={joined}
            selected={tab.network === address && tab.name === name}
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
              network={address}
              target={nick}
              content={nick}
              selected={tab.network === address && tab.name === nick}
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
          <Button
            icon={FiPlus}
            aria-label="Connect"
            onClick={this.handleConnectClick}
          />
          <Button icon={FiUser} aria-label="User" />
          <Button
            icon={FiSettings}
            aria-label="Settings"
            onClick={this.handleSettingsClick}
          />
        </div>
      </div>
    );
  }
}
