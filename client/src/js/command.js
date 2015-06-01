var _ = require('lodash');

var channelStore = require('./stores/channel');
var channelActions = require('./actions/channel');
var messageActions = require('./actions/message');
var serverActions = require('./actions/server');
var tabActions = require('./actions/tab');

messageActions.command.listen(function(line, channel, server) {
	var params = line.slice(1).split(' ');
	var command = params[0].toLowerCase();

	switch (command) {
		case 'nick':
			if (params[1]) {
				serverActions.setNick(params[1], server);
			}
			break;

		case 'quit':
			serverActions.disconnect(server);
			break;

		case 'join':
			if (params[1]) {
				channelActions.join([params[1]], server);
				tabActions.select(server, params[1]);
			}
			break;

		case 'part':
			if (params[1]) {
				channelActions.part([params[1]], server);
			} else if (channel) {
				channelActions.part([channel], server);
			}
			break;

		case 'me':
			if (params.length > 1) {
				messageActions.send('\x01ACTION ' + params.slice(1).join(' ') + '\x01', channel, server);
			}
			break;

		case 'topic':
			let topic = channelStore.getTopic(server, channel);
			if (topic) {
				messageActions.add({
					server: server,
					to: channel,
					message: topic
				});
			} else {
				messageActions.inform('No topic set', server, channel);
			}
			break;

		case 'invite':
			if (params[1] && params[2] && server) {
				channelActions.invite(params[1], params[2], server);
			} else if (params[1] && channel) {
				channelActions.invite(params[1], channel, server);
			}
			break;

		case 'kick':
			if (params[1] && channel) {
				channelActions.kick(params[1], channel, server);
			}
			break;


		case 'msg':
			if (params.length > 2) {
				let dest = params[1];
				let message = params.slice(2).join(' ');

				messageActions.send(message, dest, server);
			}
			break;

		case 'say':
			if (params.length > 1) {
				let message = params.slice(1).join(' ');

				messageActions.send(message, channel, server);
			}
			break;

		case 'whois':
			if (params[1]) {
				serverActions.whois(params[1], server);
			}
			break;

		case 'away':
			serverActions.away(params[1], server);
			break;

		case 'help':
			messageActions.inform([
				_.escape('/join <channel> - Join a channel'),
				'/part [channel] - Leave the current or entered channel',
				_.escape('/nick <nick> - Change nick'),
				'/quit - Disconnect from the current server',
				_.escape('/me <message> - Send action message'),
				'/topic - Show topic for the current channel',
				_.escape('/msg <target> <message> - Send message to the entered channel or user'),
				_.escape('/say <message> - Send message to the current chat'),
				'/away [message] - Set or clear away message'
			], server, channel);
			break;
	}
});