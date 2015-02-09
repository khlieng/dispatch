var Reflux = require('reflux');
var _ = require('lodash');

var serverStore = require('./server');
var actions = require('../actions/tab');
var channelActions = require('../actions/channel');
var serverActions = require('../actions/server');

var selectedTab = {};
var stored = localStorage.selectedTab;

if (stored) {
	selectedTab = JSON.parse(stored);
}

var selectedTabStore = Reflux.createStore({
	init: function() {
		this.listenToMany(actions);
		this.listenTo(channelActions.part, 'part');
		this.listenTo(serverActions.disconnect, 'disconnect');
		this.listenTo(channelActions.load, 'loadChannels');
		this.listenTo(serverActions.load, 'loadServers');
	},

	select: function(server, channel) {
		selectedTab.server = server;
		selectedTab.channel = channel || null;

		if (channel) {
			selectedTab.name = channel;
		} else {
			selectedTab.name = serverStore.getName(server);
		}

		this.trigger(selectedTab);
	},

	part: function(channels, server) {
		if (server === selectedTab.server && 
			channels.indexOf(selectedTab.channel) !== -1) {
			selectedTab.channel = null;
			selectedTab.name = null;
			this.trigger(selectedTab);
		}
	},

	disconnect: function(server) {
		if (server === selectedTab.server) {
			selectedTab = {};
			this.trigger(selectedTab);
		}
	},

	loadChannels: function(channels) {
		if (selectedTab.channel && !_.find(channels, { name: selectedTab.channel })) {
			selectedTab.channel = null;
			selectedTab.name = null;
			this.trigger(selectedTab);
		}
	},

	loadServers: function(servers) {
		if (selectedTab.server && !_.find(servers, { address: selectedTab.server })) {
			selectedTab = {};
			this.trigger(selectedTab);
		}
	},

	getServer: function() {
		return selectedTab.server;
	},

	getChannel: function() {
		return selectedTab.channel;
	},

	getState: function() {
		return selectedTab;
	}
});

selectedTabStore.listen(function(selectedTab) {
	localStorage.selectedTab = JSON.stringify(selectedTab);
});

module.exports = selectedTabStore;