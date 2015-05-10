var React = require('react');
var Reflux = require('reflux');
var _ = require('lodash');

var util = require('../util');
var searchStore = require('../stores/search');
var selectedTabStore = require('../stores/selectedTab');
var searchActions = require('../actions/search');

var Search = React.createClass({
	mixins: [
		Reflux.connect(searchStore),
		Reflux.connect(selectedTabStore, 'selectedTab')
	],

	getInitialState: function() {
		var state = _.extend({}, searchStore.getState());
		state.selectedTab = selectedTabStore.getState();

		return state;
	},

	componentDidUpdate: function(prevProps, prevState) {
		if (!prevState.show && this.state.show) {
			this.refs.input.getDOMNode().focus();
		}
	},

	handleChange: function(e) {
		var tab = this.state.selectedTab;

		if (tab.channel) {
			searchActions.search(tab.server, tab.channel, e.target.value);
		}
	},

	render: function() {
		var style = {
			display: this.state.show ? 'block' : 'none'
		};

		var results = _.map(this.state.results, (result) => {
			return (
				<p key={result.id}>{util.timestamp(new Date(result.time * 1000))} {result.from} {result.content}</p>
			);
		});

		return (
			<div className="search" style={style}>
				<input 
					ref="input"
					className="search-input" 
					type="text"
					onChange={this.handleChange} />
				<div className="search-results">{results}</div>
			</div>
		);
	}
});

module.exports = Search;