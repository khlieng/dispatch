var React = require('react');
var Reflux = require('reflux');
var _ = require('lodash');
var Infinite = require('react-infinite');

var UserListItem = require('./UserListItem.jsx');
var channelStore = require('../stores/channel');
var selectedTabStore = require('../stores/selectedTab');

var UserList = React.createClass({
	mixins: [
		Reflux.listenTo(channelStore, 'channelsChanged'),
		Reflux.listenTo(selectedTabStore, 'selectedTabChanged')
	],

	getInitialState() {
		var tab = selectedTabStore.getState();

		return {
			users: channelStore.getUsers(tab.server, tab.channel),
			selectedTab: tab,
			height: window.innerHeight - 100
		};
	},

	componentDidMount() {
		window.addEventListener('resize', this.handleResize);
	},

	componentWillUnmount() {
		window.removeEventListener('resize', this.handleResize);
	},

	channelsChanged() {
		var tab = this.state.selectedTab;

		this.setState({ users: channelStore.getUsers(tab.server, tab.channel) });
	},

	selectedTabChanged(tab) {
		this.setState({
			selectedTab: tab,
			users: channelStore.getUsers(tab.server, tab.channel)
		});
	},

	handleResize() {
		this.setState({ height: window.innerHeight - 100 });
	},

	render() {
		var tab = this.state.selectedTab;
		var users = [];
		var style = {};

		if (!tab.channel || tab.channel[0] !== '#') {
			style.display = 'none';
		} else {
			users = _.map(this.state.users, (user) => {
				return <UserListItem key={user.nick} user={user} />;
			});
		}

		if (users.length !== 1) {
			return (
				<div className="userlist" style={style}>
					<Infinite containerHeight={this.state.height} elementHeight={24}>
						{users}
					</Infinite>
				</div>
			);
		} else {
			return (
				<div className="userlist" style={style}>
					<div>{users}</div>
				</div>
			);
		}		
	}
});

module.exports = UserList;