var React = require('react');
var Reflux = require('reflux');
var Router = require('react-router');
var RouteHandler = Router.RouteHandler;
var Navigation = Router.Navigation;

var TabList = require('./TabList.jsx');
var routeActions = require('../actions/route');
var PureMixin = require('../mixins/pure');

var App = React.createClass({
	mixins: [
		PureMixin,
		Navigation,
		Reflux.listenTo(routeActions.navigate, 'navigate')
	],

	navigate(path, replace) {
		if (!replace) {
			this.transitionTo(path);
		} else {
			this.replaceWith(path);
		}
	},

	render() {
		return (
			<div>
				<TabList />
				<RouteHandler />
			</div>
		);
	}
});

module.exports = App;