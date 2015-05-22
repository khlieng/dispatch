var Reflux = require('reflux');
var Immutable = require('immutable');
var _ = require('lodash');

var actions = require('../actions/server');
var tabActions = require('../actions/tab');

var servers = Immutable.Map();
var Server = Immutable.Record({
	nick: null,
	name: null
});

function getState() {
	return servers;
}

var serverStore = Reflux.createStore({
	init() {
		this.listenToMany(actions);
	},

	connect(server, nick, opts) {
		var i = server.indexOf(':');
		if (i > 0) {
			server = server.slice(0, i);
		}

		servers = servers.set(server, new Server({
			nick: nick,
			name: opts.name || server
		}));

		this.trigger(servers);
		tabActions.select(server);
	},

	disconnect(server) {
		servers = servers.delete(server);
		this.trigger(servers);
	},

	setNick(nick, server) {
		servers = servers.update(server, s => s.set('nick', nick));
		this.trigger(servers);
	},

	load(storedServers) {
		_.each(storedServers, function(server) {
			servers = servers.set(server.address, new Server(server));
		});
		this.trigger(servers);
	},

	getNick(server) {
		return servers.getIn([server, 'nick']);
	},

	getName(server) {
		return servers.getIn([server, 'name']);
	},

	getInitialState: getState,
	getState
});

module.exports = serverStore;