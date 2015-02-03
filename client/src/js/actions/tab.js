var Reflux = require('reflux');

var routeActions = require('./route');

var tabActions = Reflux.createActions([
	'select'
]);

tabActions.select.preEmit = function() {
	routeActions.navigate('app');
};

module.exports = tabActions;