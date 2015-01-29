var React = require('react');
var Reflux = require('reflux');
var _ = require('lodash');

var channelStore = require('../stores/channel');
var selectedTabStore = require('../stores/selectedTab');

var UserList = React.createClass({
	mixins: [
		Reflux.connect(channelStore, 'channels'),
		Reflux.connect(selectedTabStore, 'selectedTab')
	],

	getInitialState: function() {
		return {
			channels: channelStore.getState(),
			selectedTab: selectedTabStore.getState()
		};
	},

	render: function() {
		var users = null;
		var tab = this.state.selectedTab;

		if (tab.channel && this.state.channels[tab.server]) {
			var channel = this.state.channels[tab.server][tab.channel];
			if (channel) {
				users = _.map(channel.users, function(user) {
					return <p>{user.renderName}</p>;
				});
			}
		}

		return (
			<div className="userlist">{users}</div>
		);
	}
});

module.exports = UserList;