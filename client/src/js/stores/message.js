var Reflux = require('reflux');
var Immutable = require('immutable');
var _ = require('lodash');

var serverStore = require('./server');
var channelStore = require('./channel');
var actions = require('../actions/message');
var serverActions = require('../actions/server');
var channelActions = require('../actions/channel');

var messages = Immutable.Map();
var empty = Immutable.List();

var Message = Immutable.Record({
	id: null,
	server: null,
	from: null,
	to: null,
	message: '',
	time: null,
	type: null,
	lines: []
});

function addMessage(message, dest) {
	message.time = new Date();

	if (message.message.indexOf('\x01ACTION') === 0) {
		var from = message.from;
		message.from = null;
		message.type = 'action';
		message.message = from + message.message.slice(7);
	}

	messages = messages.updateIn([message.server, dest], empty, list => list.push(new Message(message)));
}

var messageStore = Reflux.createStore({
	init() {
		this.listenToMany(actions);
		this.listenTo(serverActions.disconnect, 'disconnect');
		this.listenTo(channelActions.part, 'part');
	},

	send(message, to, server) {
		addMessage({
			server: server,
			from: serverStore.getNick(server),
			to: to,
			message: message
		}, to);

		this.trigger(messages);
	},

	add(message) {
		var dest = message.to || message.from;
		if (message.from && message.from.indexOf('.') !== -1) {
			dest = message.server;
		}

		addMessage(message, dest);
		this.trigger(messages);
	},

	broadcast(message, server, user) {
		channelStore.getChannels(server).forEach((channel, channelName) => {
			if (!user || (user && channel.get('users').find(u => u.nick === user))) {
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

	inform(message, server, channel) {
		if (_.isArray(message)) {
			_.each(message, (msg) => {
				addMessage({
					server: server,
					to: channel,
					message: msg,
					type: 'info'
				}, channel || server);
			});
		} else {
			addMessage({
				server: server,
				to: channel,
				message: message,
				type: 'info'
			}, channel || server);
		}

		this.trigger(messages);
	},

	disconnect(server) {
		messages = messages.delete(server);
		this.trigger(messages);
	},

	part(channels, server) {
		_.each(channels, function(channel) {
			messages = messages.deleteIn([server, channel]);
		});
		this.trigger(messages);
	},

	getMessages(server, dest) {
		return messages.getIn([server, dest]) || empty;
	},

	getState() {
		return messages;
	}
});

module.exports = messageStore;