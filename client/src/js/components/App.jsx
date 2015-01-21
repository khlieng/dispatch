var React = require('react');

var TabList = require('./TabList.jsx');
var Chat = require('./Chat.jsx');

var App = React.createClass({
	render: function() {
		return (
			<div>
				<TabList />
				<Chat />
			</div>
		);
	}
});

module.exports = App;