var React = require('react');
var Reflux = require('reflux');
var _ = require('lodash');
var Infinite = require('react-infinite');

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

	handleResize: function() {
		this.setState({ height: window.innerHeight - 100 });
	},

	render: function() {
		var users = [];
		var tab = this.state.selectedTab;
		var style = {};

		if (!tab.channel || tab.channel[0] !== '#') {
			style.display = 'none';
		} else {
			users = _.map(channelStore.getUsers(tab.server, tab.channel), (user) => {
				return <UserListItem key={user.nick} user={user} />;
			});
		}

		return (
			<div className="userlist" style={style}>
				<Infinite containerHeight={this.state.height} elementHeight={24}>
					{users}
				</Infinite>
			</div>
		);
	}
});

module.exports = UserList;