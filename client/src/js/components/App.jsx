var React = require('react');
var RouteHandler = require('react-router').RouteHandler;

var TabList = require('./TabList.jsx');

var App = React.createClass({
	render: function() {
		return (
			<div>
				<TabList />
				<RouteHandler />
			</div>
		);
	}
});

module.exports = App;