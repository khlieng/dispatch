import { COMMAND } from 'state/actions';
import { join, part, invite, kick, setTopic } from 'state/channels';
import { sendMessage, raw } from 'state/messages';
import { setNick, disconnect, whois, away } from 'state/servers';
import { select } from 'state/tab';
import { find } from 'utils';
import createCommandMiddleware, { beforeHandler, notFoundHandler } from './middleware/command';

const help = [
  '/join <channel> - Join a channel',
  '/part [channel] - Leave the current or specified channel',
  '/nick <nick> - Change nick',
  '/quit - Disconnect from the current server',
  '/me <message> - Send action message',
  '/topic [topic] - Show or set topic in the current channel',
  '/msg <target> <message> - Send message to the specified channel or user',
  '/say <message> - Send message to the current chat',
  '/invite <nick> [channel] - Invite user to the current or specified channel',
  '/kick <nick> - Kick user from the current channel',
  '/whois <nick> - Get information about user',
  '/away [message] - Set or clear away message',
  '/raw [message] - Send raw IRC message to the current server',
  '/help [command]... - Print help for all or the specified command(s)'
];

const text = content => ({ content });
const error = content => ({ content, type: 'error' });
const prompt = content => ({ content, type: 'prompt' });
const findHelp = cmd => find(help, line => line.slice(1, line.indexOf(' ')) === cmd);

export default createCommandMiddleware(COMMAND, {
  join({ dispatch, server }, channel) {
    if (channel) {
      if (channel[0] !== '#') {
        return error('Bad channel name');
      }
      dispatch(join([channel], server));
      dispatch(select(server, channel));
    } else {
      return error('Missing channel');
    }
  },

  part({ dispatch, server, channel, isChannel }, partChannel) {
    if (partChannel) {
      dispatch(part([partChannel], server));
    } else if (isChannel) {
      dispatch(part([channel], server));
    } else {
      return error('This is not a channel');
    }
  },

  nick({ dispatch, server }, nick) {
    if (nick) {
      dispatch(setNick(nick, server));
    } else {
      return error('Missing nick');
    }
  },

  quit({ dispatch, server }) {
    dispatch(disconnect(server));
  },

  me({ dispatch, server, channel }, ...message) {
    const msg = message.join(' ');
    if (msg !== '') {
      dispatch(sendMessage(`\x01ACTION ${msg}\x01`, channel, server));
    } else {
      return error('Messages can not be empty');
    }
  },

  topic({ dispatch, getState, server, channel }, ...newTopic) {
    if (newTopic.length > 0) {
      dispatch(setTopic(newTopic.join(' '), channel, server));
    } else {
      const topic = getState().channels.getIn([server, channel, 'topic']);
      if (topic) {
        return text(topic);
      }
      return 'No topic set';
    }
  },

  msg({ dispatch, server }, target, ...message) {
    if (!target) {
      return error('Missing nick/channel');
    }

    const msg = message.join(' ');
    if (msg !== '') {
      dispatch(sendMessage(message.join(' '), target, server));
      dispatch(select(server, target));
    } else {
      return error('Messages can not be empty');
    }
  },

  say({ dispatch, server, channel }, ...message) {
    if (!channel) {
      return error('Messages can only be sent to channels or users');
    }

    const msg = message.join(' ');
    if (msg !== '') {
      dispatch(sendMessage(message.join(' '), channel, server));
    } else {
      return error('Messages can not be empty');
    }
  },

  invite({ dispatch, server, channel, isChannel }, user, inviteChannel) {
    if (!inviteChannel && !isChannel) {
      return error('This is not a channel');
    }

    if (user && inviteChannel) {
      dispatch(invite(user, inviteChannel, server));
    } else if (user && channel) {
      dispatch(invite(user, channel, server));
    } else {
      return error('Missing nick');
    }
  },

  kick({ dispatch, server, channel, isChannel }, user) {
    if (!isChannel) {
      return error('This is not a channel');
    }

    if (user) {
      dispatch(kick(user, channel, server));
    } else {
      return error('Missing nick');
    }
  },

  whois({ dispatch, server }, user) {
    if (user) {
      dispatch(whois(user, server));
    } else {
      return error('Missing nick');
    }
  },

  away({ dispatch, server }, ...message) {
    const msg = message.join(' ');
    dispatch(away(msg, server));
    if (msg !== '') {
      return 'Away message set';
    }
    return 'Away message cleared';
  },

  raw({ dispatch, server }, ...message) {
    if (message.length > 0 && message[0] !== '') {
      const cmd = `${message[0].toUpperCase()} ${message.slice(1).join(' ')}`;
      dispatch(raw(cmd, server));
      return prompt(`=> ${cmd}`);
    }
    return [
      prompt('=> /raw'),
      error('Missing message')
    ];
  },

  help(_, ...commands) {
    if (commands.length > 0) {
      const cmdHelp = commands.filter(findHelp).map(findHelp);
      if (cmdHelp.length > 0) {
        return text(cmdHelp);
      }
      return error('Unable to find any help :(');
    }
    return text(help);
  },

  [beforeHandler](_, command, ...params) {
    if (command !== 'raw') {
      return prompt(`=> /${command} ${params.join(' ')}`);
    }
  },

  [notFoundHandler](ctx, command, ...params) {
    if (command === command.toUpperCase()) {
      return this.raw(ctx, command, ...params);
    }
    return error(`=> /${command}: No such command`);
  }
});
