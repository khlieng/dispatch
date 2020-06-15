import { COMMAND } from 'state/actions';
import { join, part, invite, kick, setTopic } from 'state/channels';
import { sendMessage, raw } from 'state/messages';
import { setNick, disconnect, whois, away } from 'state/networks';
import { openPrivateChat } from 'state/privateChats';
import { select } from 'state/tab';
import { find, isChannel } from 'utils';
import createCommandMiddleware, {
  beforeHandler,
  notFoundHandler
} from './middleware/command';

const help = [
  '/join <channel> - Join a channel',
  '/part [channel] - Leave the current or specified channel',
  '/nick <nick> - Change nick',
  '/quit - Disconnect from the current network',
  '/me <message> - Send action message',
  '/topic [topic] - Show or set topic in the current channel',
  '/msg <target> <message> - Send message to the specified channel or user',
  '/say <message> - Send message to the current chat',
  '/invite <nick> [channel] - Invite user to the current or specified channel',
  '/kick <nick> - Kick user from the current channel',
  '/whois <nick> - Get information about user',
  '/away [message] - Set or clear away message',
  '/raw [message] - Send raw IRC message to the current network',
  '/help [command]... - Print help for all or the specified command(s)'
];

const text = content => ({ content });
const error = content => ({ content, type: 'error' });
const prompt = content => ({ content, type: 'prompt' });
const findHelp = cmd =>
  find(help, line => line.slice(1, line.indexOf(' ')) === cmd);

export default createCommandMiddleware(COMMAND, {
  join({ dispatch, network }, channel) {
    if (channel) {
      if (channel[0] !== '#') {
        return error('Bad channel name');
      }
      dispatch(join([channel], network));
      dispatch(select(network, channel));
    } else {
      return error('Missing channel');
    }
  },

  part({ dispatch, network, channel, inChannel }, partChannel) {
    if (partChannel) {
      dispatch(part([partChannel], network));
    } else if (inChannel) {
      dispatch(part([channel], network));
    } else {
      return error('This is not a channel');
    }
  },

  nick({ dispatch, network }, nick) {
    if (nick) {
      dispatch(setNick(nick, network));
    } else {
      return error('Missing nick');
    }
  },

  quit({ dispatch, network }) {
    dispatch(disconnect(network));
  },

  me({ dispatch, network, channel }, ...message) {
    const msg = message.join(' ');
    if (msg !== '') {
      dispatch(sendMessage(`\x01ACTION ${msg}\x01`, channel, network));
    } else {
      return error('Messages can not be empty');
    }
  },

  topic({ dispatch, getState, network, channel }, ...newTopic) {
    if (newTopic.length > 0) {
      dispatch(setTopic(newTopic.join(' '), channel, network));
      return;
    }
    if (channel) {
      const { topic } = getState().channels[network][channel];
      if (topic) {
        return text(topic);
      }
    }
    return 'No topic set';
  },

  msg({ dispatch, network }, target, ...message) {
    if (!target) {
      return error('Missing nick/channel');
    }

    const msg = message.join(' ');
    if (msg !== '') {
      dispatch(sendMessage(message.join(' '), target, network));
      if (!isChannel(target)) {
        dispatch(openPrivateChat(network, target));
      }
      dispatch(select(network, target));
    } else {
      return error('Messages can not be empty');
    }
  },

  say({ dispatch, network, channel }, ...message) {
    if (!channel) {
      return error('Messages can only be sent to channels or users');
    }

    const msg = message.join(' ');
    if (msg !== '') {
      dispatch(sendMessage(message.join(' '), channel, network));
    } else {
      return error('Messages can not be empty');
    }
  },

  invite({ dispatch, network, channel, inChannel }, user, inviteChannel) {
    if (!inviteChannel && !inChannel) {
      return error('This is not a channel');
    }

    if (user && inviteChannel) {
      dispatch(invite(user, inviteChannel, network));
    } else if (user && channel) {
      dispatch(invite(user, channel, network));
    } else {
      return error('Missing nick');
    }
  },

  kick({ dispatch, network, channel, inChannel }, user) {
    if (!inChannel) {
      return error('This is not a channel');
    }

    if (user) {
      dispatch(kick(user, channel, network));
    } else {
      return error('Missing nick');
    }
  },

  whois({ dispatch, network }, user) {
    if (user) {
      dispatch(whois(user, network));
    } else {
      return error('Missing nick');
    }
  },

  away({ dispatch, network }, ...message) {
    const msg = message.join(' ');
    dispatch(away(msg, network));
    if (msg !== '') {
      return 'Away message set';
    }
    return 'Away message cleared';
  },

  raw({ dispatch, network }, ...message) {
    if (message.length > 0 && message[0] !== '') {
      const cmd = `${message[0].toUpperCase()} ${message.slice(1).join(' ')}`;
      dispatch(raw(cmd, network));
      return prompt(`=> ${cmd}`);
    }
    return [prompt('=> /raw'), error('Missing message')];
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
    return this.raw(ctx, command, ...params);
  }
});
