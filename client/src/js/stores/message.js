var Reflux = require('reflux');
var actions = require('../actions/message.js');

var messages = {};

var messageStore = Reflux.createStore({
	init: function() {
		this.listenToMany(actions);
	},

	add: function(message) {
		var dest = message.to || message.from;
		if (message.from.indexOf('.') !== -1) {
			dest = message.server;
		}
		
		if (!(dest in messages)) {
			messages[dest] = [message];
		} else {
			messages[dest].push(message);
		}

		this.trigger(messages);
	},

	getState: function() {
		return messages;
	}
});

module.exports = messageStore;