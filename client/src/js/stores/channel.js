var Reflux = require('reflux');
var _ = require('lodash');

var actions = require('../actions/channel.js');

var channels = {};

function initChannel(server, channel) {
	if (!(server in channels)) {
		channels[server] = {};
		channels[server][channel] = { users: [] };
	} else if (!(channel in channels[server])) {
		channels[server][channel] = { users: [] };
	}
}

var channelStore = Reflux.createStore({
	init: function() {
		this.listenToMany(actions);
	},

	part: function(partChannels, server) {
		_.each(partChannels, function(channel) {
			delete channels[server][channel];
		});
		this.trigger(channels);
	},

	addUser: function(user, server, channel) {
		initChannel(server, channel);
		channels[server][channel].users.push(user);
		this.trigger(channels);
	},

	removeUser: function(user, server, channel) {
		_.pull(channels[server][channel].users, user);
		this.trigger(channels);
	},

	removeUserAll: function(user, server) {
		_.each(channels[server], function(channel) {
			_.pull(channel.users, user);
		});
		this.trigger(channels);
	},

	setUsers: function(users, server, channel) {
		initChannel(server, channel);
		channels[server][channel].users = users;
		this.trigger(channels);
	},

	setTopic: function(topic, server, channel) {
		channels[server][channel].topic = topic;
		this.trigger(channels);
	},

	load: function(storedChannels) {
		_.each(storedChannels, function(channel) {
			initChannel(channel.server, channel.name);
			channels[channel.server][channel.name].users = channel.users;
			channels[channel.server][channel.name].topic = channel.topic;
		});
		this.trigger(channels);
	},

	getState: function() {
		return channels;
	}
});

module.exports = channelStore;