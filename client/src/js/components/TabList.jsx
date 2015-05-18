var React = require('react');
var Reflux = require('reflux');
var _ = require('lodash');

var TabListItem = require('./TabListItem.jsx');
var channelStore = require('../stores/channel');
var privateChatStore = require('../stores/privateChat');
var serverStore = require('../stores/server');
var routeActions = require('../actions/route');

var TabList = React.createClass({
	mixins: [
		Reflux.connect(serverStore, 'servers'),
		Reflux.connect(channelStore, 'channels'),
		Reflux.connect(privateChatStore, 'privateChats')
	],

	getInitialState() {
		return {
			servers: serverStore.getState(),
			channels: channelStore.getState(),
			privateChats: privateChatStore.getState()
		};
	},

	handleConnectClick() {
		routeActions.navigate('connect');
	},

	handleSettingsClick() {
		routeActions.navigate('settings');
	},

	render() {
		var tabs = _.map(this.state.channels, (server, address) => {
			var serverTabs = _.map(server, (channel, name) => {
				return (
					<TabListItem 
						server={address} 
						channel={name}
						name={name}>
					</TabListItem>
				);
			});

			_.each(this.state.privateChats[address], (chat, nick) => {
				serverTabs.push(
					<TabListItem 
						server={address} 
						channel={nick}
						name={nick}>
					</TabListItem>
				);	
			});

			serverTabs.unshift(
				<TabListItem 
					server={address} 
					channel={null}
					name={serverStore.getName(address)}>
				</TabListItem>
			);

			return serverTabs;
		});

		return (
			<div className="tablist">
				<button className="button-connect" onClick={this.handleConnectClick}>Connect</button>
				{tabs}
				<div className="side-buttons">
					<i className="icon-user"></i>
					<i className="icon-cog" onClick={this.handleSettingsClick}></i>
				</div>
			</div>
		);
	}
});

module.exports = TabList;