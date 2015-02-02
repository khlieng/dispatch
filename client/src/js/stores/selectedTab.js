var Reflux = require('reflux');
var _ = require('lodash');

var serverStore = require('./server');
var actions = require('../actions/tab');
var channelActions = require('../actions/channel');

var selectedTab = {};
var stored = localStorage.selectedTab;

if (stored) {
	selectedTab = JSON.parse(stored);
}

var selectedTabStore = Reflux.createStore({
	init: function() {
		this.listenToMany(actions);
		this.listenTo(channelActions.part, 'part');
	},

	select: function(server, channel) {
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
		var self = this;
		if (server === selectedTab.server) {
			_.each(channels, function(channel) {
				if (channel === selectedTab.channel) {
					delete selectedTab.channel;
					self.trigger(selectedTab);
					return;
				}
			});
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