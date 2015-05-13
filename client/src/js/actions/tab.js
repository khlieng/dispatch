var Reflux = require('reflux');

var routeActions = require('./route');

var tabActions = Reflux.createActions([
	'select'
]);

module.exports = tabActions;