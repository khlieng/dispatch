var React = require('react');
var Reflux = require('reflux');

var selectedTabStore = require('../stores/selectedTab');
var messageActions = require('../actions/message');
var tabActions = require('../actions/tab');

var MessageInput = React.createClass({
	mixins: [
		Reflux.connect(selectedTabStore, 'selectedTab')
	],

	getInitialState: function() {
		return {
			selectedTab: selectedTabStore.getState()
		};
	},

	handleKey: function(e) {
		if (e.which === 13 && e.target.value) {
			var tab = this.state.selectedTab;

			if (e.target.value[0] === '/') {
				messageActions.command(e.target.value, tab.channel, tab.server);
			} else {
				messageActions.send(e.target.value, tab.channel, tab.server);
			}
			e.target.value = '';
		}
	},

	render: function() {
		return (
			<div className="message-input-wrap">
				<input className="message-input" type="text" onKeyDown={this.handleKey} />
			</div>
		);
	}
});

module.exports = MessageInput;