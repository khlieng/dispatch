var Reflux = require('reflux');

var inputHistoryActions = Reflux.createActions([
	'add',
	'reset',
	'increment',
	'decrement'
]);

module.exports = inputHistoryActions;