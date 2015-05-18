var Reflux = require('reflux');
var _ = require('lodash');

var actions = require('../actions/server');
var tabActions = require('../actions/tab');

var servers = {};

var serverStore = Reflux.createStore({
	init() {
		this.listenToMany(actions);
	},

	connect(server, nick, opts) {
		var i = server.indexOf(':');
		if (i > 0) {
			server = server.slice(0, i);
		}

		servers[server] = {
			address: server,
			nick: nick,
			name: opts.name || server
		};

		this.trigger(servers);
		tabActions.select(server);
	},

	disconnect(server) {
		delete servers[server];
		this.trigger(servers);
	},

	setNick(nick, server) {
		servers[server].nick = nick;
		this.trigger(servers);
	},

	load(storedServers) {
		_.each(storedServers, function(server) {
			servers[server.address] = server;
		});
		this.trigger(servers);
	},

	getNick(server) {
		if (servers[server]) {
			return servers[server].nick;
		}
		return null;
	},

	getName(server) {
		if (servers[server]) {
			return servers[server].name;
		}
		return null;
	},

	getState() {
		return servers;
	}
});

module.exports = serverStore;