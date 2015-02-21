var Reflux = require('reflux');

var socket = require('../socket');

var serverActions = Reflux.createActions([
	'connect',
	'disconnect',
	'whois',
	'away',
	'setNick',
	'load'
]);

serverActions.connect.preEmit = (server, nick, opts) => {
	socket.send('connect', {
		server,
		nick,
		username: opts.username || nick,
		password: opts.password,
		realname: opts.realname || nick,
		tls: opts.tls || false,
		name: opts.name || server
	});
};

serverActions.disconnect.preEmit = (server) => {
	socket.send('quit', { server });
};

serverActions.whois.preEmit = (user, server) => {
	socket.send('whois', { server, user });
};

serverActions.away.preEmit = (message, server) => {
	socket.send('away', { server, message });
};

serverActions.setNick.preEmit = (nick, server) => {
	socket.send('nick', {
		server,
		new: nick
	});
};

module.exports = serverActions;