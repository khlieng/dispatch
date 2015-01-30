var Reflux = require('reflux');

var actions = require('../actions/privateChat');
var messageActions = require('../actions/message');

var privateChats = {};

function initChat(server, nick) {
	if (!(server in privateChats)) {
		privateChats[server] = {};
		privateChats[server][nick] = {};

		return true;
	} else if (!(nick in privateChats[server])) {
		privateChats[server][nick] = {};

		return true;
	}
	return false;
}

var privateChatStore = Reflux.createStore({
	init: function() {
		this.listenToMany(actions);
		this.listenTo(messageActions.add, 'messageAdded');
	},

	open: function(server, nick) {
		if (initChat(server, nick)) {
			this.trigger(privateChats);
		}
	},

	close: function(server, nick) {
		delete privateChat[server][nick];
		this.trigger(privateChats);
	},

	messageAdded: function(message) {
		if (!message.to && message.from.indexOf('.') === -1) {
			this.open(message.server, message.from);
		}
	},

	getState: function() {
		return privateChats;
	}
});

module.exports = privateChatStore;