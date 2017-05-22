import createCommandMiddleware from './middleware/command';
import { COMMAND } from './actions';
import { setNick, disconnect, whois, away } from './actions/server';
import { join, part, invite, kick } from './actions/channel';
import { select } from './actions/tab';
import { sendMessage, addMessage, raw } from './actions/message';

const help = [
  '/join <channel> - Join a channel',
  '/part [channel] - Leave the current or specified channel',
  '/nick <nick> - Change nick',
  '/quit - Disconnect from the current server',
  '/me <message> - Send action message',
  '/topic - Show topic for the current channel',
  '/msg <target> <message> - Send message to the specified channel or user',
  '/say <message> - Send message to the current chat',
  '/invite <user> [channel] - Invite user to the current or specified channel',
  '/kick <user> - Kick user from the current channel',
  '/whois <user> - Get information about user',
  '/away [message] - Set or clear away message',
  '/raw [message] - Send raw IRC message to the current server'
];

export default createCommandMiddleware(COMMAND, {
  join({ dispatch, server }, channel) {
    if (channel) {
      dispatch(join([channel], server));
      dispatch(select(server, channel));
    }
  },

  part({ dispatch, server, channel }, partChannel) {
    if (partChannel) {
      dispatch(part([partChannel], server));
    } else {
      dispatch(part([channel], server));
    }
  },

  nick({ dispatch, server }, nick) {
    if (nick) {
      dispatch(setNick(nick, server));
    }
  },

  quit({ dispatch, server }) {
    dispatch(disconnect(server));
  },

  me({ dispatch, server, channel }, ...params) {
    if (params.length > 0) {
      dispatch(sendMessage(`\x01ACTION ${params.join(' ')}\x01`, channel, server));
    }
  },

  topic({ dispatch, getState, server, channel }) {
    const topic = getState().channels.getIn([server, channel, 'topic']);
    if (topic) {
      dispatch(addMessage({
        server,
        to: channel,
        content: topic
      }));
    } else {
      return 'No topic set';
    }
  },

  msg({ dispatch, server }, target, ...message) {
    if (target && message) {
      dispatch(sendMessage(message.join(' '), target, server));
    }
  },

  say({ dispatch, server, channel }, ...message) {
    if (channel && message) {
      dispatch(sendMessage(message.join(' '), channel, server));
    }
  },

  invite({ dispatch, server, channel }, user, inviteChannel) {
    if (user && inviteChannel) {
      dispatch(invite(user, inviteChannel, server));
    } else if (user && channel) {
      dispatch(invite(user, channel, server));
    }
  },

  kick({ dispatch, server, channel }, user) {
    if (user && channel) {
      dispatch(kick(user, channel, server));
    }
  },

  whois({ dispatch, server }, user) {
    if (user) {
      dispatch(whois(user, server));
    }
  },

  away({ dispatch, server }, message) {
    dispatch(away(message, server));
  },

  raw({ dispatch, server }, ...message) {
    if (message) {
      const cmd = message.join(' ');
      dispatch(raw(cmd, server));
      return `=> ${cmd}`;
    }
  },

  help() {
    return help;
  },

  commandNotFound(_, command) {
    return `The command /${command} was not found`;
  }
});
