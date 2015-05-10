var Reflux = require('reflux');

var actions = require('../actions/search');

var state = {
	show: false,
	results: []
};

var searchStore = Reflux.createStore({
	init: function() {
		this.listenToMany(actions);
	},

	searchDone: function(results) {
		state.results = results;
		this.trigger(state);
	},

	toggle: function() {
		state.show = !state.show;
		this.trigger(state);
	},

	getState: function() {
		return state;
	}
});

module.exports = searchStore;