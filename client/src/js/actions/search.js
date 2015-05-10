var Reflux = require('reflux');

var socket = require('../socket');

var searchActions = Reflux.createActions([
	'search',
	'searchDone',
	'toggle'
]);

searchActions.search.preEmit = (server, channel, phrase) => {
	socket.send('search', { server, channel, phrase });
};

module.exports = searchActions;