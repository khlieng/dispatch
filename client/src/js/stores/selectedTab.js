var Reflux = require('reflux');
var _ = require('lodash');

var serverStore = require('./server');
var actions = require('../actions/tab');
var channelActions = require('../actions/channel');
var serverActions = require('../actions/server');
var routeActions = require('../actions/route');
var privateChatActions = require('../actions/privateChat');

var selectedTab = {};
var history = [];

function selectPrevTab() {
	history.pop();

	if (history.length > 0) {
		selectedTab = _.extend({}, history[history.length - 1]);
		return true;
	}

	return false;
}

function updateChannelName(name) {
	selectedTab.channel = name;
	selectedTab.name = name;
	history[history.length - 1].channel = name;
	history[history.length - 1].name = name; 
}

var selectedTabStore = Reflux.createStore({
	init: function() {
		this.listenToMany(actions);
		this.listenTo(channelActions.part, 'part');
		this.listenTo(privateChatActions.close, 'close');
		this.listenTo(serverActions.disconnect, 'disconnect');
		this.listenTo(channelActions.addUser, 'userAdded');
		this.listenTo(channelActions.load, 'loadChannels');
		this.listenTo(serverActions.load, 'loadServers');
		this.listenTo(routeActions.navigate, 'navigate');
	},

	select: function(server, channel = null) {
		selectedTab.server = server;
		selectedTab.channel = channel;

		if (channel) {
			selectedTab.name = channel;
		} else {
			selectedTab.name = serverStore.getName(server);
		}

		history.push(_.extend({}, selectedTab));

		this.trigger(selectedTab);
	},

	part: function(channels, server) {
		if (server === selectedTab.server && 
			channels.indexOf(selectedTab.channel) !== -1) {
			if (!selectPrevTab()) {
				selectedTab.channel = null;
				selectedTab.name = serverStore.getName(server);
			}
			
			this.trigger(selectedTab);
		}
	},

	close: function(server, nick) {
		if (server === selectedTab.server &&
			nick === selectedTab.channel) {
			if (!selectPrevTab()) {
				selectedTab.channel = null;
				selectedTab.name = serverStore.getName(server);
			}
			
			this.trigger(selectedTab);
		}
	},

	disconnect: function(server) {
		if (server === selectedTab.server) {
			_.remove(history, { server: server });

			if (!selectPrevTab()) {
				selectedTab = {};
			}

			this.trigger(selectedTab);
		}
	},

	userAdded: function(user, server, channel) {
		if (selectedTab.channel &&
			server === selectedTab.server &&
			user === serverStore.getNick(server) &&
			channel.toLowerCase().indexOf(selectedTab.channel.toLowerCase()) !== -1) {
			// Update the selected channel incase the casing is different
			updateChannelName(channel);
			this.trigger(selectedTab);
		}
	},

	loadChannels: function(channels) {
		_.each(channels, (channel) => {
			if (channel.server === selectedTab.server &&
				channel.name !== selectedTab.channel &&
				channel.name.indexOf(selectedTab.channel) !== -1) {
				// Handle double hashtag channel names, only a single hashtag
				// gets added to the channel in the URL on page load
				updateChannelName(channel.name);
				this.trigger(selectedTab);

				return false;
			}
		});
	},

	loadServers: function(servers) {
		var server = _.find(servers, { address: selectedTab.server });

		if (server && !selectedTab.channel) {
			selectedTab.name = server.name;
			history[history.length - 1].name = server.name;

			this.trigger(selectedTab);
		}
	},

	navigate: function(route) {
		if (route.indexOf('.') === -1 && selectedTab.server) {
			selectedTab.server = null;
			selectedTab.channel = null;
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
	var channel = selectedTab.channel;

	if (selectedTab.server) {
		if (channel) {
			while (channel[0] === '#') {
				channel = channel.slice(1);
			}
			routeActions.navigate('/' + selectedTab.server + '/' + channel);
		} else {
			routeActions.navigate('/' + selectedTab.server);
		}
	} else if (_.size(serverStore.getState()) === 0) {
		routeActions.navigate('connect');
	}
	
	localStorage.selectedTab = JSON.stringify(selectedTab);
});

module.exports = selectedTabStore;