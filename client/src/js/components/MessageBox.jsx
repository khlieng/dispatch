var React = require('react');
var Reflux = require('reflux');
var _ = require('lodash');
var Infinite = require('react-infinite');

var util = require('../util');
var messageLineStore = require('../stores/messageLine');
var selectedTabStore = require('../stores/selectedTab');
var messageActions = require('../actions/message');

var MessageBox = React.createClass({
	mixins: [
		Reflux.connect(messageLineStore, 'lines'),
		Reflux.connect(selectedTabStore, 'selectedTab')
	],

	getInitialState: function() {
		return {
			lines: messageLineStore.getState(),
			selectedTab: selectedTabStore.getState(),
			height: window.innerHeight - 100
		};
	},

	componentDidMount: function() {
		window.addEventListener('resize', this.handleResize);
	},

	componentWillUnmount: function() {
		window.removeEventListener('resize', this.handleResize);
	},

	componentWillUpdate: function() {
		var el = this.refs.list.getDOMNode();
		this.autoScroll = el.scrollTop + el.offsetHeight === el.scrollHeight;
	},

	componentDidUpdate: function() {
		this.updateWidth();

		if (this.autoScroll) {
			var el = this.refs.list.getDOMNode();
			el.scrollTop = el.scrollHeight;
		}
	},

	handleResize: function() {
		this.updateWidth();
		this.setState({ height: window.innerHeight - 100 });
	},

	updateWidth: function() {
		var width = this.refs.list.getDOMNode().firstChild.offsetWidth;

		if (this.width !== width) {
			this.width = width;
			messageActions.setWrapWidth(width);
		}
	},

	render: function() {
		var tab = this.state.selectedTab;
		var dest = tab.channel || tab.server;
		var style = {}

		/*var messages = _.map(messageStore.getMessages(tab.server, dest), function(message) {
			var messageClass = 'message';

			if (message.type) {
				messageClass += ' message-' + message.type;
			}

			return (
				<p className={messageClass}>
					<span className="message-time">{util.timestamp(message.time)}</span>
					{ message.from ? <span className="message-sender"> {message.from}</span> : null }
					{' ' + message.message}
				</p>
			);
		});*/

		if (!tab.channel || tab.channel[0] !== '#') {
			style.right = 0;
		}

		var lines = _.map(this.state.lines, function(line) {
			return <p className="message">{line}</p>;
		});

		if (lines.length !== 1) {
			return (
				<div className="messagebox" style={style}>
					<Infinite ref="list" containerHeight={this.state.height} elementHeight={24}>
						{lines}
					</Infinite>
				</div>
			);
		} else {
			return (
				<div className="messagebox" style={style}>
					<div ref="list">{lines}</div>
				</div>
			);
		}
	}
});

module.exports = MessageBox;