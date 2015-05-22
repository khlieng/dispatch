var React = require('react');
var Reflux = require('reflux');
var _ = require('lodash');
var Infinite = require('react-infinite');
var Autolinker = require('autolinker');

var MessageHeader = require('./MessageHeader.jsx');
var MessageLine = require('./MessageLine.jsx');
var messageLineStore = require('../stores/messageLine');
var selectedTabStore = require('../stores/selectedTab');
var messageActions = require('../actions/message');
var PureMixin = require('../mixins/pure');

var MessageBox = React.createClass({
	mixins: [
		PureMixin,
		Reflux.connect(messageLineStore, 'messages'),
		Reflux.connect(selectedTabStore, 'selectedTab')
	],

	getInitialState() {
		return {
			height: window.innerHeight - 100
		};
	},

	componentDidMount() {
		window.addEventListener('resize', this.handleResize);
	},

	componentWillUnmount() {
		window.removeEventListener('resize', this.handleResize);
	},

	componentWillUpdate() {
		var el = this.refs.list.getDOMNode();
		this.autoScroll = el.scrollTop + el.offsetHeight === el.scrollHeight;
	},

	componentDidUpdate() {
		this.updateWidth();

		if (this.autoScroll) {
			var el = this.refs.list.getDOMNode();
			el.scrollTop = el.scrollHeight;
		}
	},

	handleResize() {
		this.updateWidth();
		this.setState({ height: window.innerHeight - 100 });
	},

	updateWidth() {
		var width = this.refs.list.getDOMNode().firstChild.offsetWidth;

		if (this.width !== width) {
			this.width = width;
			messageActions.setWrapWidth(width);
		}
	},

	render() {
		var tab = this.state.selectedTab;
		var dest = tab.channel || tab.server;
		var lines = [];

		this.state.messages.forEach((message, j) => {
			var key = message.server + dest + j;			

			lines.push(<MessageHeader key={key} message={message} />);

			for (var i = 1; i < message.lines.length; i++) {
				lines.push(
					<MessageLine key={key + '-' + i} type={message.type} line={message.lines[i]} />
				);
			}
		});

		if (lines.length !== 1) {
			return (
				<div className="messagebox">
					<Infinite ref="list" containerHeight={this.state.height} elementHeight={24}>
						{lines}
					</Infinite>
				</div>
			);
		} else {
			return (
				<div className="messagebox">
					<div ref="list">{lines}</div>
				</div>
			);
		}
	}
});

module.exports = MessageBox;