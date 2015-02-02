var React = require('react');
var Reflux = require('reflux');
var _ = require('lodash');

var UserListItem = require('./UserListItem.jsx');
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
		var style = {};

		if (!tab.channel || tab.channel[0] !== '#') {
			style.display = 'none';
		} else {
			users = _.map(channelStore.getUsers(tab.server, tab.channel), function(user) {
				return <UserListItem key={user.nick} user={user} />;
			});
		}

		return (
			<div className="userlist" style={style}>{users}</div>
		);
	}
});

module.exports = UserList;