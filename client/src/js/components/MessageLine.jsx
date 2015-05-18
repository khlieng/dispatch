var React = require('react');
var Autolinker = require('autolinker');

var PureMixin = require('../mixins/pure');

var MessageLine = React.createClass({
	mixins: [PureMixin],

	render() {
		var line = Autolinker.link(this.props.line, { keepOriginalText: true });
		var messageClass = 'message';
		var style = {
			paddingLeft: window.messageIndent + 'px'
		};

		if (this.props.type) {
			messageClass += ' message-' + this.props.type;
		}

		return (
			<p className={messageClass} style={style}>
				<span dangerouslySetInnerHTML={{ __html: line }}></span>
			</p>
		);
	}
});

module.exports = MessageLine;