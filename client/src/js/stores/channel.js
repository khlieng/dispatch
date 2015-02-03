var Reflux = require('reflux');
var _ = require('lodash');

var actions = require('../actions/channel');
var serverActions = require('../actions/server');

var channels = {};

function initChannel(server, channel) {
	if (!(server in channels)) {
		channels[server] = {};
		channels[server][channel] = { users: [] };
	} else if (!(channel in channels[server])) {
		channels[server][channel] = { users: [] };
	}
}

function createUser(nick, mode) {
	return updateRenderName({
		nick: nick,
		renderName: nick,
		mode: mode || ''
	});
}

function loadUser(users, nick) {
	var mode;

	if (nick[0] === '@') {
		mode = 'o';
	} else if (nick[0] === '+') {
		mode = 'v';
	}

	if (mode) {
		nick = nick.slice(1);
	}

	users.push(createUser(nick, mode));
}

function updateRenderName(user) {
	if (user.mode.indexOf('o') !== -1) {
		user.renderName = '@' + user.nick;
	} else if (user.mode.indexOf('v') !== -1) {
		user.renderName = '+' + user.nick;
	} else {
		user.renderName = user.nick;
	}

	return user;
}

function sortUsers(server, channel) {
	channels[server][channel].users.sort(function(a, b) {
		a = a.renderName.toLowerCase();
		b = b.renderName.toLowerCase();

		if (a[0] === '@' && b[0] !== '@') {
			return -1;
		}
		if (b[0] === '@' && a[0] !== '@') {
			return 1;
		}
		if (a[0] === '+' && b[0] !== '+') {
			return -1;
		}
		if (b[0] === '+' && a[0] !== '+') {
			return 1;
		}
		if (a < b) {
			return -1;
		}
		if (a > b) {
			return 1;
		}
		return 0;
	});
}

var channelStore = Reflux.createStore({
	init: function() {
		this.listenToMany(actions);
		this.listenTo(serverActions.connect, 'addServer');
		this.listenTo(serverActions.load, 'loadServers');
	},

	part: function(partChannels, server) {
		_.each(partChannels, function(channel) {
			delete channels[server][channel];
		});
		this.trigger(channels);
	},

	addUser: function(user, server, channel) {
		initChannel(server, channel);
		channels[server][channel].users.push(createUser(user));
		sortUsers(server, channel);
		this.trigger(channels);
	},

	removeUser: function(user, server, channel) {
		if (channels[server][channel]) {
			_.remove(channels[server][channel].users, { nick: user });
			this.trigger(channels);
		}
	},

	removeUserAll: function(user, server) {
		_.each(channels[server], function(channel) {
			_.remove(channel.users, { nick: user });
		});
		this.trigger(channels);
	},

	renameUser: function(oldNick, newNick, server) {
		_.each(channels[server], function(channel, channelName) {
			var user = _.find(channel.users, { nick: oldNick });
			if (user) {
				user.nick = newNick;
				updateRenderName(user);
				sortUsers(server, channelName);
			}
		});
		this.trigger(channels);
	},

	setUsers: function(users, server, channel) {
		initChannel(server, channel);
		var chan = channels[server][channel];

		chan.users = [];

		_.each(users, function(user) {
			loadUser(chan.users, user);
		});

		sortUsers(server, channel);
		this.trigger(channels);
	},

	setTopic: function(topic, server, channel) {
		channels[server][channel].topic = topic;
		this.trigger(channels);
	},

	setMode: function(mode) {
		var user = _.find(channels[mode.server][mode.channel].users, { nick: mode.user });
		if (user) {
			_.each(mode.remove, function(mode) {
				user.mode = user.mode.replace(mode, '');
			});
			_.each(mode.add, function(mode) {
				user.mode += mode;
			});

			updateRenderName(user);
			sortUsers(mode.server, mode.channel);
			this.trigger(channels);
		}
	},

	load: function(storedChannels) {
		_.each(storedChannels, function(channel) {
			initChannel(channel.server, channel.name);
			var chan = channels[channel.server][channel.name];
			
			chan.users = [];
			chan.topic = channel.topic;

			_.each(channel.users, function(user) {
				loadUser(chan.users, user);
			});

			sortUsers(channel.server, channel.name);
		});

		this.trigger(channels);
	},

	addServer: function(server) {
		if (!(server in channels)) {
			channels[server] = {};
			this.trigger(channels);
		}
	},

	loadServers: function(storedServers) {
		_.each(storedServers, function(server) {
			if (!(server.address in channels)) {
				channels[server.address] = {};
			}
		});
		this.trigger(channels);
	},

	getChannels: function(server) {
		return channels[server];
	},

	getUsers: function(server, channel) {
		if (channels[server] && channels[server][channel]) {
			return channels[server][channel].users;
		}
		return [];
	},

	getState: function() {
		return channels;
	}
});

module.exports = channelStore;