var React = require('react');
var Reflux = require('reflux');
var _ = require('lodash');

var TabListItem = require('./TabListItem.jsx');
var channelStore = require('../stores/channel');
var privateChatStore = require('../stores/privateChat');
var serverStore = require('../stores/server');
var routeActions = require('../actions/route');
var PureMixin = require('../mixins/pure');

var TabList = React.createClass({
	mixins: [
		PureMixin,
		Reflux.connect(serverStore, 'servers'),
		Reflux.connect(channelStore, 'channels'),
		Reflux.connect(privateChatStore, 'privateChats')
	],

	handleConnectClick() {
		routeActions.navigate('connect');
	},

	handleSettingsClick() {
		routeActions.navigate('settings');
	},

	render() {
		var tabs = [];

		this.state.channels.forEach((server, address) => {
			tabs.push(
				<TabListItem 
					key={address}
					server={address} 
					channel={null}
					name={this.state.servers.getIn([address, 'name'])}>
				</TabListItem>
			);

			server.forEach((channel, name) => {
				tabs.push(
					<TabListItem
						key={address + name} 
						server={address} 
						channel={name}
						name={name}>
					</TabListItem>
				);
			});

			if (this.state.privateChats.has(address)) {
				this.state.privateChats.get(address).forEach(nick => {
					tabs.push(
						<TabListItem 
							key={address + nick}
							server={address} 
							channel={nick}
							name={nick}>
						</TabListItem>
					);
				});
			}
		});

		return (
			<div className="tablist">
				<button className="button-connect" onClick={this.handleConnectClick}>Connect</button>
				<div className="tab-container">{tabs}</div>
				<div className="side-buttons">
					<i className="icon-user"></i>
					<i className="icon-cog" onClick={this.handleSettingsClick}></i>
				</div>
			</div>
		);
	}
});

module.exports = TabList;