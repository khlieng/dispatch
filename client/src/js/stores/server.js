var Reflux = require('reflux');
var _ = require('lodash');

var actions = require('../actions/server');
var tabActions = require('../actions/tab');

var servers = {};

var serverStore = Reflux.createStore({
	init: function() {
		this.listenToMany(actions);
	},

	connect: function(server, nick, username, tls, name) {
		var i = server.indexOf(':');
		if (i > 0) {
			server = server.slice(0, i);
		}

		servers[server] = {
			address: server,
			nick: nick,
			username: username,
			name: name || server
		};

		this.trigger(servers);
		tabActions.select(server);
	},

	load: function(storedServers) {
		_.each(storedServers, function(server) {
			servers[server.address] = server;
		});
		this.trigger(servers);
	},

	getNick: function(server) {
		return servers[server].nick;
	},

	getName: function(server) {
		if (servers[server]) {
			return servers[server].name;
		}
		return null;
	},

	getState: function() {
		return servers;
	}
});

module.exports = serverStore;