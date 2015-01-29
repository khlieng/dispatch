var Reflux = require('reflux');
var _ = require('lodash');

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

	getState: function() {
		return selectedTab;
	}
});

selectedTabStore.listen(function(selected) {
	localStorage.selectedTab = JSON.stringify(selected);
});

module.exports = selectedTabStore;