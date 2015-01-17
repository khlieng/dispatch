var React = require('react');
var Reflux = require('reflux');
var _ = require('lodash');
var channelStore = require('../stores/channel.js');
var selectedTabStore = require('../stores/selectedTab.js');

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
		
		if (tab.channel) {
			users = _.map(this.state.channels[tab.server][tab.channel].users, function(user) {
				return <p>{user}</p>;
			});
		}

		return (
			<div className="userlist">{users}</div>
		);
	}
});

module.exports = UserList;