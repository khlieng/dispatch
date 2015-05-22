var Reflux = require('reflux');
var Immutable = require('immutable');
var _ = require('lodash');

var serverStore = require('./server');
var actions = require('../actions/tab');
var channelActions = require('../actions/channel');
var serverActions = require('../actions/server');
var routeActions = require('../actions/route');
var privateChatActions = require('../actions/privateChat');

var Tab = Immutable.Record({
	server: null,
	channel: null,
	name: null
});

var selectedTab = new Tab();
var history = [];

function selectPrevTab() {
	history.pop();

	if (history.length > 0) {
		selectedTab = history[history.length - 1];
		return true;
	}

	return false;
}

function updateChannelName(name) {
	selectedTab = selectedTab.set('channel', name).set('name', name);
	history[history.length - 1] = selectedTab;
}

function getState() {
	return selectedTab;
}

var selectedTabStore = Reflux.createStore({
	init() {
		this.listenToMany(actions);
		this.listenTo(channelActions.part, 'part');
		this.listenTo(privateChatActions.close, 'close');
		this.listenTo(serverActions.disconnect, 'disconnect');
		this.listenTo(channelActions.addUser, 'userAdded');
		this.listenTo(channelActions.load, 'loadChannels');
		this.listenTo(serverActions.load, 'loadServers');
		this.listenTo(routeActions.navigate, 'navigate');
	},

	select(server, channel = null) {
		selectedTab = new Tab({ 
			server, 
			channel, 
			name: channel || serverStore.getName(server) 
		});

		history.push(selectedTab);

		this.trigger(selectedTab);
	},

	part(channels, server) {
		if (server === selectedTab.server && 
			channels.indexOf(selectedTab.channel) !== -1) {
			if (!selectPrevTab()) {
				selectedTab = selectedTab
					.set('channel', null)
					.set('name', serverStore.getName(server));
			}
			
			this.trigger(selectedTab);
		}
	},

	close(server, nick) {
		if (server === selectedTab.server &&
			nick === selectedTab.channel) {
			if (!selectPrevTab()) {
				selectedTab = selectedTab
					.set('channel', null)
					.set('name', serverStore.getName(server));
			}
			
			this.trigger(selectedTab);
		}
	},

	disconnect(server) {
		if (server === selectedTab.server) {
			_.remove(history, { server: server });

			if (!selectPrevTab()) {
				selectedTab = new Tab();
			}

			this.trigger(selectedTab);
		}
	},

	userAdded(user, server, channel) {
		if (selectedTab.channel &&
			server === selectedTab.server &&
			user === serverStore.getNick(server) &&
			channel.toLowerCase().indexOf(selectedTab.channel.toLowerCase()) !== -1) {
			// Update the selected channel incase the casing is different
			updateChannelName(channel);
			this.trigger(selectedTab);
		}
	},

	loadChannels(channels) {
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

	loadServers(servers) {
		var server = _.find(servers, { address: selectedTab.server });

		if (server && !selectedTab.channel) {
			selectedTab = selectedTab.set('name', server.name);
			history[history.length - 1] = selectedTab;

			this.trigger(selectedTab);
		}
	},

	navigate(route) {
		if (route.indexOf('.') === -1 && selectedTab.server) {
			selectedTab = new Tab();
			this.trigger(selectedTab);
		}
	},

	getServer() {
		return selectedTab.server;
	},

	getChannel() {
		return selectedTab.channel;
	},

	getInitialState: getState,
	getState
});

selectedTabStore.listen(selectedTab => {
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
	} else if (serverStore.getState().size === 0) {
		routeActions.navigate('connect');
	}
	
	localStorage.selectedTab = JSON.stringify(selectedTab);
});

module.exports = selectedTabStore;