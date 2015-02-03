var React = require('react');
var Reflux = require('reflux');
var _ = require('lodash');

var channelStore = require('../stores/channel');
var privateChatStore = require('../stores/privateChat');
var serverStore = require('../stores/server');
var selectedTabStore = require('../stores/selectedTab');
var tabActions = require('../actions/tab');
var routeActions = require('../actions/route');

var TabList = React.createClass({
	mixins: [
		Reflux.connect(channelStore, 'channels'),
		Reflux.connect(privateChatStore, 'privateChats'),
		Reflux.connect(selectedTabStore, 'selectedTab'),
		Reflux.connect(serverStore, 'servers')
	],

	getInitialState: function() {
		return {
			channels: channelStore.getState(),
			privateChats: privateChatStore.getState(),
			selectedTab: selectedTabStore.getState(),
			servers: serverStore.getState()
		};
	},

	handleConnectClick: function() {
		routeActions.navigate('connect');
	},

	handleSettingsClick: function() {
		routeActions.navigate('settings');
	},

	render: function() {
		var self = this;
		var tabClass;
		var selected = this.state.selectedTab;

		var tabs = _.map(this.state.channels, function(server, address) {
			var channels = _.map(server, function(channel, name) {				
				if (address === selected.server &&
					name === selected.channel) {
					tabClass = 'selected';
				} else {
					tabClass = '';
				}

				return (
					<p 
						className={tabClass} 
						onClick={tabActions.select.bind(null, address, name)}>
							{name}
					</p>
				);
			});

			_.each(self.state.privateChats[address], function(chat, nick) {
				if (address === selected.server &&
					nick === selected.channel) {
					tabClass = 'selected';
				} else {
					tabClass = '';
				}

				channels.push(
					<p 
						className={tabClass} 
						onClick={tabActions.select.bind(null, address, nick)}>
							{nick}
					</p>
				);	
			});

			if (address === selected.server &&
				selected.channel === null) {
				tabClass = 'tab-server selected';
			} else {
				tabClass = 'tab-server';
			}

			channels.unshift(
				<p 
					className={tabClass} 
					onClick={tabActions.select.bind(null, address, null)}>
						{serverStore.getName(address)}
				</p>
			);

			return channels;
		});

		return (
			<div className="tablist">
				<button className="button-connect" onClick={this.handleConnectClick}>Connect</button>
				{tabs}
				<div className="side-buttons">
					<button onClick={this.handleSettingsClick}>Settings</button>
				</div>
			</div>
		);
	}
});

module.exports = TabList;