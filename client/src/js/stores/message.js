var Reflux = require('reflux');

var serverStore = require('./server');
var actions = require('../actions/message');

var messages = {};

function addMessage(message, dest) {
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
			message: message,
			time: new Date()
		}, to);

		this.trigger(messages);
	},

	add: function(message) {
		var dest = message.to || message.from;
		if (message.from.indexOf('.') !== -1) {
			dest = message.server;
		}

		message.time = new Date();

		addMessage(message, dest);
		this.trigger(messages);
	},

	getState: function() {
		return messages;
	}
});

module.exports = messageStore;