var Reflux = require('reflux');

var routeActions = require('./route');

var tabActions = Reflux.createActions([
	'select'
]);

tabActions.select.preEmit = (server, channel) => {
	if (channel) {
		routeActions.navigate('/' + server + '/' + channel.slice(1))
	} else {
		routeActions.navigate('/' + server);
	}
};

module.exports = tabActions;