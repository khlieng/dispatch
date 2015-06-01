var React = require('react');
var Autolinker = require('autolinker');

var util = require('../util');
var privateChatActions = require('../actions/privateChat');
var tabActions = require('../actions/tab');

var MessageHeader = React.createClass({
	shouldComponentUpdate(nextProps) {
		return nextProps.message.lines[0] !== this.props.message.lines[0];
	},

	handleSenderClick() {
		var message = this.props.message;

		privateChatActions.open(message.server, message.from);
		tabActions.select(message.server, message.from);
	},

	render() {
		var message = this.props.message;
		var sender = null;
		var messageClass = 'message';
		var line = Autolinker.link(message.lines[0], { keepOriginalText: true });

		if (message.from) {
			sender = (
				<span>
					{' '}
					<span className="message-sender" onClick={this.handleSenderClick}>
						{message.from}
					</span>
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
				<span dangerouslySetInnerHTML={{ __html: ' ' + line }}></span>
			</p>
		);
	}
});

module.exports = MessageHeader;