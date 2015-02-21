var Reflux = require('reflux');
var _ = require('lodash');

var util = require('../util');
var messageStore = require('./message');
var selectedTabStore = require('./selectedTab');
var messageActions = require('../actions/message');

var width = window.innerWidth;
window.charWidth = util.stringWidth(' ', '16px Droid Sans Mono');
window.messageIndent = 6 * charWidth;

var tab = selectedTabStore.getState();
var messages;

function wrap() {
	messages = messageStore.getMessages(tab.server, tab.channel || tab.server);
	util.wrapMessages(messages, width, charWidth, messageIndent);
}

wrap();

var messageLineStore = Reflux.createStore({
	init: function() {
		this.listenTo(messageActions.setWrapWidth, 'setWrapWidth');
		this.listenTo(messageStore, 'messagesChanged');
		this.listenTo(selectedTabStore, 'selectedTabChanged');
	},

	setWrapWidth: function(w) {
		width = w;

		util.wrapMessages(messages, width, charWidth, messageIndent);
		this.trigger(messages);
	},

	messagesChanged: function() {
		wrap();
		this.trigger(messages);
	},

	selectedTabChanged: function(selectedTab) {
		tab = selectedTab;
		
		wrap();		
		this.trigger(messages);
	},

	getState: function() {
		return messages;
	}
});

module.exports = messageLineStore;