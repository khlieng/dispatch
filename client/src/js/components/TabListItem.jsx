var React = require('react');
var Reflux = require('reflux');

var selectedTabStore = require('../stores/selectedTab');
var tabActions = require('../actions/tab');

var TabListItem = React.createClass({
	mixins: [Reflux.connect(selectedTabStore)],

	getInitialState: function() {
		return selectedTabStore.getState();
	},

	handleClick: function() {
		tabActions.select(this.props.server, this.props.channel);
	},

	render: function() {
		var classes = [];

		if (!this.props.channel) {
			classes.push('tab-server');
		}

		if (this.props.server === this.state.server &&
			this.props.channel === this.state.channel) {
			classes.push('selected');
		}

		return (
			<p className={classes.join(' ')} onClick={this.handleClick}>{this.props.name}</p>
		);
	}
});

module.exports = TabListItem;