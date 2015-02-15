var Reflux = require('reflux');
var _ = require('lodash');

var util = require('../util');
var messageStore = require('./message');
var selectedTabStore = require('./selectedTab');
var messageActions = require('../actions/message');

var width = window.innerWidth;
var charWidth = util.stringWidth(' ', '16px Ubuntu Mono');

var tab = selectedTabStore.getState();
var messages;
var lines;

wrap();

function wrap() {
	messages = messageStore.getMessages(tab.server, tab.channel || tab.server);

	lines = util.wrap(_.map(messages, function(message) {
		var line = util.timestamp(message.time);
		if (message.from) {
			line += ' ' + message.from;
		}
		line += ' ' + message.message;

		return line;
	}), width, charWidth);
}

var messageLineStore = Reflux.createStore({
	init: function() {
		this.listenTo(messageActions.setWrapWidth, 'setWrapWidth');
		this.listenTo(messageStore, 'messagesChanged');
		this.listenTo(selectedTabStore, 'selectedTabChanged');
	},

	setWrapWidth: function(w) {
		width = w;
		
		wrap();
		this.trigger(lines);
	},

	messagesChanged: function() {
		wrap();
		this.trigger(lines);
	},

	selectedTabChanged: function(selectedTab) {
		tab = selectedTab;
		
		wrap();		
		this.trigger(lines);
	},

	getState: function() {
		return lines;
	}
});

module.exports = messageLineStore;