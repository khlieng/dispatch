var Reflux = require('reflux');

var routeActions = require('./route');

var tabActions = Reflux.createActions([
	'select'
]);

tabActions.select.preEmit = (server, channel) => {
	if (channel) {
		while (channel[0] === '#') {
			channel = channel.slice(1);
		}
		routeActions.navigate('/' + server + '/' + channel);
	} else {
		routeActions.navigate('/' + server);
	}
};

module.exports = tabActions;