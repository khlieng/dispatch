import React from 'react';
import Reflux from 'reflux';
import { Router } from 'react-router';
import TabList from './TabList.jsx';
import routeActions from '../actions/route';
import tabActions from '../actions/tab';
import PureMixin from '../mixins/pure';

export default React.createClass({
	mixins: [
		PureMixin,
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
		const { history } = this.props;
		if (!replace) {
			history.pushState(null, path);
		} else {
			history.replaceState(null, path);
		}
	},

	hideMenu() {
		this.setState({ menuToggled: false });
	},

	toggleMenu() {
		this.setState({ menuToggled: !this.state.menuToggled });
	},

	render() {
		const mainClass = this.state.menuToggled ? 'main-container off-canvas' : 'main-container';

		return (
			<div>
				<TabList menuToggled={this.state.menuToggled} />
				<div className={mainClass}>
					{this.props.children}
				</div>
			</div>
		);
	}
});
