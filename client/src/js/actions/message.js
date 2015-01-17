var Reflux = require('reflux');

var messageActions = Reflux.createActions([
	'send',
	'add',
	'selectTab'
]);

messageActions.send.preEmit = function() {

};

module.exports = messageActions;