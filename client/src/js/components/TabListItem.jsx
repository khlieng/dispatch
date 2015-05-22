var React = require('react');
var Reflux = require('reflux');

var selectedTabStore = require('../stores/selectedTab');
var tabActions = require('../actions/tab');
var PureMixin = require('../mixins/pure');

var TabListItem = React.createClass({
	mixins: [
		PureMixin,
		Reflux.connect(selectedTabStore, 'tab')
	],

	handleClick() {
		tabActions.select(this.props.server, this.props.channel);
	},

	render() {
		var classes = [];

		if (!this.props.channel) {
			classes.push('tab-server');
		}

		if (this.props.server === this.state.tab.server &&
			this.props.channel === this.state.tab.channel) {
			classes.push('selected');
		}

		return (
			<p className={classes.join(' ')} onClick={this.handleClick}>{this.props.name}</p>
		);
	}
});

module.exports = TabListItem;