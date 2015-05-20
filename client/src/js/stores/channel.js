var Reflux = require('reflux');
var Immutable = require('immutable');
var _ = require('lodash');

var actions = require('../actions/channel');
var serverActions = require('../actions/server');

var channels = Immutable.Map();
var empty = Immutable.List();

var User = Immutable.Record({
	nick: null,
	renderName: null,
	mode: ''
});

function createUser(nick, mode) {
	return updateRenderName(new User({
		nick: nick,
		renderName: nick,
		mode: mode || ''
	}));
}

function loadUser(nick) {
	var mode;

	if (nick[0] === '@') {
		mode = 'o';
	} else if (nick[0] === '+') {
		mode = 'v';
	}

	if (mode) {
		nick = nick.slice(1);
	}

	return createUser(nick, mode);
}

function updateRenderName(user) {
	var name = user.nick;

	if (user.mode.indexOf('o') !== -1) {
		name = '@' + name;
	} else if (user.mode.indexOf('v') !== -1) {
		name = '+' + name;
	}

	return user.set('renderName', name);
}

function sortUsers(a, b) {
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
}

var channelStore = Reflux.createStore({
	init() {
		this.listenToMany(actions);
		this.listenTo(serverActions.connect, 'addServer');
		this.listenTo(serverActions.disconnect, 'removeServer');
		this.listenTo(serverActions.load, 'loadServers');
	},

	part(partChannels, server) {
		_.each(partChannels, function(channel) {
			channels = channels.deleteIn([server, channel]);
		});
		this.trigger(channels);
	},

	addUser(user, server, channel) {
		channels = channels.updateIn([server, channel, 'users'], empty, users => {
			return users.push(createUser(user)).sort(sortUsers);
		});
		this.trigger(channels);
	},

	removeUser(user, server, channel) {
		if (channels.hasIn([server, channel])) {
			channels = channels.updateIn([server, channel, 'users'], users => users.filter(u => u.nick !== user));
			this.trigger(channels);
		}
	},

	removeUserAll(user, server) {
		channels.get(server).forEach((v, k) => {
			channels = channels.updateIn([server, k, 'users'], users => users.filter(u => u.nick !== user));
		});
		this.trigger(channels);
	},

	renameUser(oldNick, newNick, server) {
		channels.get(server).forEach((v, k) => {
			channels = channels.updateIn([server, k, 'users'], users => {
				var i = users.findIndex(user => user.nick === oldNick);
				return users.update(i, user => {
					return updateRenderName(user.set('nick', newNick));
				}).sort(sortUsers);
			});
		});
		this.trigger(channels);
	},

	setUsers(users, server, channel) {
		users = _.map(users, user => loadUser(user)).sort(sortUsers);
		channels = channels.setIn([server, channel, 'users'], Immutable.List(users));
		this.trigger(channels);
	},

	setTopic(topic, server, channel) {
		channels = channels.setIn([server, channel, 'topic'], topic);
		this.trigger(channels);
	},

	setMode(mode) {
		var i = channels.getIn([mode.server, mode.channel, 'users']).findIndex(u => u.nick === mode.user);

		channels = channels.updateIn([mode.server, mode.channel, 'users', i], user => {
			_.each(mode.remove, function(mode) {
				user = user.set('mode', user.mode.replace(mode, ''));
			});
			
			user = user.set('mode', user.mode + mode.add);

			return updateRenderName(user);
		});
		channels = channels.updateIn([mode.server, mode.channel, 'users'], users => users.sort(sortUsers));
		this.trigger(channels);
	},

	load(storedChannels) {
		_.each(storedChannels, function(channel) {
			channels = channels.setIn([channel.server, channel.name], Immutable.Map({
				users: Immutable.List(),
				topic: channel.topic
			}));
		});
		this.trigger(channels);
	},

	addServer(server) {
		if (!channels.has(server)) {
			channels = channels.set(server, Immutable.Map());
			this.trigger(channels);
		}
	},

	removeServer(server) {
		channels = channels.delete(server);
		this.trigger(channels);
	},

	loadServers(storedServers) {
		_.each(storedServers, function(server) {
			if (!channels.has(server.address)) {
				channels = channels.set(server.address, Immutable.Map());
			}
		});
		this.trigger(channels);
	},

	getChannels(server) {
		return channels.get(server);
	},

	getUsers(server, channel) {
		return channels.getIn([server, channel, 'users']) || empty;
	},

	getTopic(server, channel) {
		return channels.getIn([server, channel, 'topic']);
	},

	getState() {
		return channels;
	}
});

module.exports = channelStore;