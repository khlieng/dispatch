var React = require('react');
var Reflux = require('reflux');
var Router = require('react-router');
var RouteHandler = Router.RouteHandler;
var Navigation = Router.Navigation;

var TabList = require('./TabList.jsx');
var routeActions = require('../actions/route');
var tabActions = require('../actions/tab');
var PureMixin = require('../mixins/pure');

var App = React.createClass({
	mixins: [
		PureMixin,
		Navigation,
		Reflux.listenTo(routeActions.navigate, 'navigate'),
		Reflux.listenTo(tabActions.hideMenu, 'hideMenu'),
		Reflux.listenTo(tabActions.toggleMenu, 'toggleMenu')
	],

	getInitialState() {
		return {
			menuToggled: false
		};
	},

	navigate(path, replace) {
		if (!replace) {
			this.transitionTo(path);
		} else {
			this.replaceWith(path);
		}
	},

	hideMenu() {
		this.setState({ menuToggled: false });
	},

	toggleMenu() {
		this.setState({ menuToggled: !this.state.menuToggled });
	},

	render() {
		var mainClass = this.state.menuToggled ? 'main-container off-canvas' : 'main-container';

		return (
			<div>
				<TabList menuToggled={this.state.menuToggled} />
				<div className={mainClass}>
					<RouteHandler />
				</div>
			</div>
		);
	}
});

module.exports = App;