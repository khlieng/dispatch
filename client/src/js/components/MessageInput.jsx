var React = require('react');
var Reflux = require('reflux');

var inputHistoryStore = require('../stores/inputHistory');
var selectedTabStore = require('../stores/selectedTab');
var messageActions = require('../actions/message');
var inputHistoryActions = require('../actions/inputHistory');
var tabActions = require('../actions/tab');

var MessageInput = React.createClass({
	mixins: [
		Reflux.connect(selectedTabStore, 'selectedTab'),
		Reflux.connect(inputHistoryStore, 'history'),
		Reflux.listenTo(tabActions.select, 'tabSelected')
	],

	getInitialState: function() {
		return {
			selectedTab: selectedTabStore.getState(),
			history: inputHistoryStore.getState(),
			value: ''
		};
	},

	componentDidMount: function() {
		this.refs.input.getDOMNode().focus();
	},

	tabSelected: function() {
		this.refs.input.getDOMNode().focus();
	},

	handleKey: function(e) {
		if (e.which === 13 && e.target.value) {
			var tab = this.state.selectedTab;

			if (e.target.value[0] === '/') {
				messageActions.command(e.target.value, tab.channel, tab.server);
			} else {
				messageActions.send(e.target.value, tab.channel, tab.server);
			}

			inputHistoryActions.add(e.target.value);
			inputHistoryActions.reset();
			this.setState({ value: '' });
		} else if (e.which === 38) {
			e.preventDefault();
			inputHistoryActions.increment();
		} else if (e.which === 40) {
			inputHistoryActions.decrement();
		} else if (e.key === 'Backspace' || e.key === 'Delete') {
			inputHistoryActions.reset();
		} else if (e.key === 'Unidentified') {
			inputHistoryActions.reset();
		}
	},

	handleChange: function(e) {
		this.setState({ value: e.target.value });
	},

	render: function() {
		return (
			<div className="message-input-wrap">
				<input 
					ref="input" 
					className="message-input" 
					type="text" 
					value={this.state.history || this.state.value}
					onKeyDown={this.handleKey}
					onChange={this.handleChange} />
			</div>
		);
	}
});

module.exports = MessageInput;