var Reflux = require('reflux');
var _ = require('lodash');

var actions = require('../actions/inputHistory');

var HISTORY_MAX_LENGTH = 128;

var history = [];
var index = -1;

var inputHistoryStore = Reflux.createStore({
	init: function() {
		this.listenToMany(actions);
	},

	add: function(line) {
		if (line.trim() && line !== history[0]) {
			history.unshift(line);
			if (history.length > HISTORY_MAX_LENGTH) {
				history.pop();
			}
			this.trigger(history[index]);
		}
	},

	reset: function() {
		index = -1;
		this.trigger(history[index]);
	},

	increment: function() {
		if (index !== history.length - 1) {
			index++;
			this.trigger(history[index]);
		}
	},

	decrement: function() {
		if (index !== -1) {
			index--;
			this.trigger(history[index]);
		}
	},

	getState: function() {
		if (index !== -1) {
			return history[index];
		}
		return null;
	}
});

module.exports = inputHistoryStore;