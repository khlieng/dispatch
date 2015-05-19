var Reflux = require('reflux');
var _ = require('lodash');

var util = require('../util');
var messageStore = require('./message');
var selectedTabStore = require('./selectedTab');
var messageActions = require('../actions/message');

var tab = selectedTabStore.getState();
var width = window.innerWidth;
var messages;
var prev;

function updateCharWidth() {
	window.charWidth = util.stringWidth(' ', '16px Droid Sans Mono');
	window.messageIndent = 6 * charWidth;
}

function wrap() {
	var next = messageStore.getMessages(tab.server, tab.channel || tab.server);
	if (next !== prev) {
		prev = next;
		messages = util.wrapMessages(next, width, charWidth, messageIndent);
		return true;
	}
	return false;
}

var messageLineStore = Reflux.createStore({
	init() {
		updateCharWidth();
		wrap();

		// Temporary hack incase this runs before the font has loaded
		setTimeout(updateCharWidth, 1000);

		this.listenTo(messageActions.setWrapWidth, 'setWrapWidth');
		this.listenTo(messageStore, 'messagesChanged');
		this.listenTo(selectedTabStore, 'selectedTabChanged');
	},

	setWrapWidth(w) {
		width = w;
		messages = util.wrapMessages(messages, width, charWidth, messageIndent);
		this.trigger(messages);
	},

	messagesChanged() {
		if (wrap()) {
			this.trigger(messages);
		}
	},

	selectedTabChanged(selectedTab) {
		tab = selectedTab;
		
		if (wrap()) {
			this.trigger(messages);
		}
	},

	getState() {
		return messages;
	}
});

module.exports = messageLineStore;