var Reflux = require('reflux');
var _ = require('lodash');

var actions = require('../actions/server');

var servers = {};

var serverStore = Reflux.createStore({
	init: function() {
		this.listenToMany(actions);
	},

	connect: function(server, nick, username) {
		servers[server] = {
			address: server,
			nick: nick,
			username: username
		};
		this.trigger(servers);
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

	getState: function() {
		return servers;
	}
});

module.exports = serverStore;