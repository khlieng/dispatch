var React = require('react');
var TabList = require('./TabList.jsx');
var MessageBox = require('./MessageBox.jsx');
var UserList = require('./UserList.jsx');

var App = React.createClass({
	render: function() {
		return (
			<div>
				<TabList />
				<MessageBox />
				<UserList />
			</div>
		);
	}
});

module.exports = App;