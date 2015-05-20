var Reflux = require('reflux');
var Immutable = require('immutable');

var actions = require('../actions/privateChat');
var messageActions = require('../actions/message');
var serverActions = require('../actions/server');

var privateChats = Immutable.Map();
var empty = Immutable.List();

var privateChatStore = Reflux.createStore({
	init() {
		this.listenToMany(actions);
		this.listenTo(messageActions.add, 'messageAdded');
		this.listenTo(serverActions.disconnect, 'disconnect');
	},

	open(server, nick) {
		privateChats = privateChats.update(server, empty, chats => chats.push(nick));
		this.trigger(privateChats);
	},

	close(server, nick) {
		privateChats = privateChats.update(server, chats => chats.delete(chats.indexOf(nick)));
		this.trigger(privateChats);
	},

	messageAdded(message) {
		if (!message.to && message.from.indexOf('.') === -1) {
			this.open(message.server, message.from);
		}
	},

	disconnect(server) {
		privateChats = privateChats.delete(server);
		this.trigger(privateChats);
	},

	getState() {
		return privateChats;
	}
});

module.exports = privateChatStore;