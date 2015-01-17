var Reflux = require('reflux');
var actions = require('../actions/tab.js');

var selectedTab = {};

var selectedTabStore = Reflux.createStore({
	init: function() {
		this.listenToMany(actions);
	},

	select: function(server, channel) {
		selectedTab.server = server;
		selectedTab.channel = channel;
		this.trigger(selectedTab);
	},

	getState: function() {
		return selectedTab;
	}
});

module.exports = selectedTabStore;