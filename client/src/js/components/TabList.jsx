var React = require('react');
var Reflux = require('reflux');
var _ = require('lodash');

var channelStore = require('../stores/channel');
var privateChatStore = require('../stores/privateChat');
var selectedTabStore = require('../stores/selectedTab');
var tabActions = require('../actions/tab');

var TabList = React.createClass({
	mixins: [
		Reflux.connect(channelStore, 'channels'),
		Reflux.connect(privateChatStore, 'privateChats'),
		Reflux.connect(selectedTabStore, 'selectedTab')
	],

	getInitialState: function() {
		return {
			channels: channelStore.getState(),
			privateChats: privateChatStore.getState(),
			selectedTab: selectedTabStore.getState()
		};
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

				return <p className={tabClass} onClick={tabActions.select.bind(null, address, name)}>{name}</p>;
			});

			_.each(self.state.privateChats[address], function(chat, nick) {
				if (address === selected.server &&
					nick === selected.channel) {
					tabClass = 'selected';
				} else {
					tabClass = '';
				}
				
				channels.push(<p className={tabClass} onClick={tabActions.select.bind(null, address, nick)}>{nick}</p>);	
			});

			if (address === selected.server &&
				selected.channel === null) {
				tabClass = 'tab-server selected';
			} else {
				tabClass = 'tab-server';
			}

			channels.unshift(<p className={tabClass} onClick={tabActions.select.bind(null, address, null)}>{address}</p>);

			return channels;
		});

		return (
			<div className="tablist">
				<button className="button-connect">Add Network</button>
				{tabs}
				<div className="side-buttons">
					<button>Settings</button>
				</div>
			</div>
		);
	}
});

module.exports = TabList;