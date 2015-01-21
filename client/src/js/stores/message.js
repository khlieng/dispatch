var Reflux = require('reflux');

var actions = require('../actions/message.js');

var messages = {};

function addMessage(message, dest) {
	if (!(dest in messages)) {
		messages[dest] = [message];
	} else {
		messages[dest].push(message);
	}
}

var messageStore = Reflux.createStore({
	init: function() {
		this.listenToMany(actions);
	},

	send: function(message, to, server) {
		addMessage({
			server: server,
			from: 'self',
			to: to,
			message: message
		}, to);

		this.trigger(messages);
	},

	add: function(message) {
		var dest = message.to || message.from;
		if (message.from.indexOf('.') !== -1) {
			dest = message.server;
		}

		addMessage(message, dest);
		this.trigger(messages);
	},

	getState: function() {
		return messages;
	}
});

module.exports = messageStore;