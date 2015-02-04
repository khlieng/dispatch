var React = require('react');
var Reflux = require('reflux');

var channelStore = require('../stores/channel');
var selectedTabStore = require('../stores/selectedTab');
var messageActions = require('../actions/message');
var channelActions = require('../actions/channel');
var serverActions = require('../actions/server');
var tabActions = require('../actions/tab');

function dispatchCommand(cmd, channel, server) {
	var params = cmd.slice(1).split(' ');

	switch (params[0].toLowerCase()) {
		case 'join':
			if (params[1]) {
				channelActions.join([params[1]], server);
				tabActions.select(server, params[1]);
			}
			break;

		case 'part':
			if (channel) {
				channelActions.part([channel], server);
			}
			break;

		case 'me':
			if (params.length > 1) {
				messageActions.send('\x01ACTION ' + params.slice(1).join(' ') + '\x01', channel, server);
			}
			break;

		case 'topic':
			var topic = channelStore.getTopic(server, channel);
			if (topic) {
				messageActions.inform(topic, server, channel);
			} else {
				messageActions.inform('No topic set', server, channel);
			}
			break;

		case 'nick':
			if (params[1]) {
				serverActions.setNick(params[1], server);
			}
			break;
	}
}

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
				dispatchCommand(e.target.value, tab.channel, tab.server);
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