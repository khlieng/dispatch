var React = require('react');
var Reflux = require('reflux');
var _ = require('lodash');
var messageStore = require('../stores/message.js');
var selectedTabStore = require('../stores/selectedTab.js');

var MessageBox = React.createClass({
	mixins: [
		Reflux.connect(messageStore, 'messages'),
		Reflux.connect(selectedTabStore, 'selectedTab')
	],

	getInitialState: function() {
		return {
			messages: messageStore.getState(),
			selectedTab: selectedTabStore.getState()
		};
	},

	componentWillUpdate: function() {
		var el = this.getDOMNode();
		this.autoScroll = el.scrollTop + el.offsetHeight === el.scrollHeight;
	},

	componentDidUpdate: function() {
		if (this.autoScroll) {
			var el = this.getDOMNode();
			el.scrollTop = el.scrollHeight;
		}
	},

	render: function() {
		var tab = this.state.selectedTab.channel || this.state.selectedTab.server;
		var messages = _.map(this.state.messages[tab], function(message) {
			return <p>{message.from ? message.from + ': ' : null}{message.message}</p>;
		});

		return (
			<div className="messagebox">{messages}</div>
		);
	}
});

module.exports = MessageBox;