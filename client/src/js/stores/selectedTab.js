var Reflux = require('reflux');
var _ = require('lodash');

var serverStore = require('./server');
var actions = require('../actions/tab');
var channelActions = require('../actions/channel');
var serverActions = require('../actions/server');

var selectedTab = {};

var selectedTabStore = Reflux.createStore({
	init: function() {
		this.listenToMany(actions);
		this.listenTo(channelActions.part, 'part');
		this.listenTo(serverActions.disconnect, 'disconnect');
		this.listenTo(channelActions.addUser, 'userAdded');
		//this.listenTo(channelActions.load, 'loadChannels');
		this.listenTo(serverActions.load, 'loadServers');
	},

	select: function(server, channel = null) {
		selectedTab.server = server;
		selectedTab.channel = channel;

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

	userAdded: function(user, server, channel) {
		// Update the selected channel incase the casing is different
		if (selectedTab.channel &&
			server === selectedTab.server &&
			user === serverStore.getNick(server) &&
			channel.toLowerCase().indexOf(selectedTab.channel.toLowerCase()) !== -1) {
			selectedTab.channel = channel;
			selectedTab.name = channel;
			this.trigger(selectedTab);
		}
	},

	loadChannels: function(channels) {
		
	},

	loadServers: function(servers) {
		var server = _.find(servers, { address: selectedTab.server });

		if (!server) {
			selectedTab = {};
			this.trigger(selectedTab);
		} else if (!selectedTab.channel) {
			selectedTab.name = server.name;
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