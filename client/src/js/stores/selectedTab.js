var Reflux = require('reflux');
var _ = require('lodash');

var actions = require('../actions/tab.js');
var channelActions = require('../actions/channel.js');

var selectedTab = {};

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

module.exports = selectedTabStore;