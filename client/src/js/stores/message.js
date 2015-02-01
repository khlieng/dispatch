var Reflux = require('reflux');
var _ = require('lodash');

var serverStore = require('./server');
var channelStore = require('./channel');
var actions = require('../actions/message');

var messages = {};

function addMessage(message, dest) {
	message.time = new Date();

	if (!(message.server in messages)) {
		messages[message.server] = {};
		messages[message.server][dest] = [message];
	} else if (!(dest in messages[message.server])) {
		messages[message.server][dest] = [message];
	} else {
		messages[message.server][dest].push(message);
	}
}

var messageStore = Reflux.createStore({
	init: function() {
		this.listenToMany(actions);
	},

	send: function(message, to, server) {
		addMessage({
			server: server,
			from: serverStore.getNick(server),
			to: to,
			message: message
		}, to);

		this.trigger(messages);
	},

	add: function(message) {
		var dest = message.to || message.from;
		if (message.from && message.from.indexOf('.') !== -1) {
			dest = message.server;
		}

		addMessage(message, dest);
		this.trigger(messages);
	},

	broadcast: function(message, server, user) {
		_.each(channelStore.getChannels(server), function(channel, channelName) {
			if (!user || (user && _.find(channel.users, { nick: user }))) {
				addMessage({
					server: server,
					to: channelName,
					message: message,
					type: 'info'
				}, channelName);
			}
		});
		this.trigger(messages);
	},

	inform: function(message, server, channel) {
		addMessage({
			server: server,
			to: channel,
			message: message,
			type: 'info'
		}, channel);
		this.trigger(messages);
	},

	getState: function() {
		return messages;
	}
});

module.exports = messageStore;