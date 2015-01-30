var Reflux = require('reflux');

var privateChatActions = Reflux.createActions([
	'open',
	'close'
]);

module.exports = privateChatActions;