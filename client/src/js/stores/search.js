var Reflux = require('reflux');
var Immutable = require('immutable');

var actions = require('../actions/search');

var state = Immutable.Map({
	show: false,
	results: Immutable.List()
});

var searchStore = Reflux.createStore({
	init() {
		this.listenToMany(actions);
	},

	searchDone(results) {
		state = state.set('results', Immutable.List(results));
		this.trigger(state);
	},

	toggle() {
		state = state.update('show', show => !show);
		this.trigger(state);
	},

	getState() {
		return state;
	}
});

module.exports = searchStore;