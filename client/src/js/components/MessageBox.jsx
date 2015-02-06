var React = require('react');
var Reflux = require('reflux');
var _ = require('lodash');

var util = require('../util');
var messageStore = require('../stores/message');
var selectedTabStore = require('../stores/selectedTab');

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
		var tab = this.state.selectedTab;
		var dest = tab.channel || tab.server;

		var messages = _.map(messageStore.getMessages(tab.server, dest), function(message) {
			var messageClass = 'message';

			if (message.type) {
				messageClass += ' message-' + message.type;
			}

			return (
				<p className={messageClass}>
					<span className="message-time">{util.timestamp(message.time)}</span>
					{ message.from ? <span className="message-sender">{message.from}</span> : null }
					{message.message}
				</p>
			);
		});
		
		return (
			<div className="messagebox">{messages}</div>
		);
	}
});

module.exports = MessageBox;