var React = require('react');
var Reflux = require('reflux');
var Autolinker = require('autolinker');

var util = require('../util');
var privateChatActions = require('../actions/privateChat');
var tabActions = require('../actions/tab');

var MessageHeader = React.createClass({
	handleSenderClick: function() {
		var message = this.props.message;

		privateChatActions.open(message.server, message.from);
		tabActions.select(message.server, message.from);
	},

	render: function() {
		var message = this.props.message;
		var sender = null;
		var messageClass = 'message';

		if (message.from) {
			sender = (
				<span 
					className="message-sender" 
					style={{ marginLeft: window.charWidth + 'px' }}
					onClick={this.handleSenderClick}>
						{message.from}
				</span>
			);
		}

		if (message.type) {
			messageClass += ' message-' + message.type;
		}

		return (
			<p className={messageClass}>
				<span className="message-time">{util.timestamp(message.time)}</span>
				{sender}
				<span dangerouslySetInnerHTML={{ __html: ' ' + Autolinker.link(message.lines[0]) }}></span>
			</p>
		);
	}
});

module.exports = MessageHeader;