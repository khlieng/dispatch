var React = require('react');

var tabActions = require('../actions/tab');

var Navicon = React.createClass({
	render() {
		return (
			<i className="icon-menu navicon" onClick={tabActions.toggleMenu}></i>
		);
	}
});

module.exports = Navicon;