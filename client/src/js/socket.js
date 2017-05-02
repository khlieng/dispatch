import { routeActions } from 'react-router-redux';
import { broadcast, inform, addMessage, addMessages } from './actions/message';
import { select } from './actions/tab';
import { normalizeChannel } from './util';

function withReason(message, reason) {
  return message + (reason ? ` (${reason})` : '');
}

function findChannels(state, server, user) {
  const channels = [];

  state.channels.get(server).forEach((channel, channelName) => {
    if (channel.get('users').find(u => u.nick === user)) {
      channels.push(channelName);
    }
  });

  return channels;
}

export default function handleSocket(socket, { dispatch, getState }) {
  const handlers = {
    message(message) {
      dispatch(addMessage(message, message.server, message.to));
    },

    pm(message) {
      dispatch(addMessage(message, message.server, message.from));
    },

    messages({ messages, server, to, prepend, next }) {
      dispatch(addMessages(messages, server, to, prepend, next));
    },

    join({ user, server, channels }) {
      const state = getState();
      const tab = state.tab.selected;
      const [joinedChannel] = channels;
      if (tab.server && tab.name) {
        const { nick } = state.servers.get(tab.server);
        if (tab.server === server &&
          nick === user &&
          tab.name !== joinedChannel &&
          normalizeChannel(tab.name) === normalizeChannel(joinedChannel)) {
          dispatch(select(server, joinedChannel));
        }
      }

      dispatch(inform(`${user} joined the channel`, server, joinedChannel));
    },

    servers(data) {
      if (!data) {
        dispatch(routeActions.replace('/connect'));
      }
    },

    part({ user, server, channel, reason }) {
      dispatch(inform(withReason(`${user} left the channel`, reason), server, channel));
    },

    quit({ user, server, reason }) {
      const channels = findChannels(getState(), server, user);
      dispatch(broadcast(withReason(`${user} quit`, reason), server, channels));
    },

    nick(data) {
      const channels = findChannels(getState(), data.server, data.old);
      dispatch(broadcast(`${data.old} changed nick to ${data.new}`, data.server, channels));
    },

    motd({ content, server }) {
      dispatch(addMessages(content.map(line => ({ content: line })), server));
    },

    whois(data) {
      const tab = getState().tab.selected;

      dispatch(inform([
        `Nick: ${data.nick}`,
        `Username: ${data.username}`,
        `Realname: ${data.realname}`,
        `Host: ${data.host}`,
        `Server: ${data.server}`,
        `Channels: ${data.channels}`
      ], tab.server, tab.name));
    },

    print({ server, message }) {
      dispatch(inform(message, server));
    }
  };

  socket.onMessage((type, data) => {
    if (type in handlers) {
      handlers[type](data);
    }

    type = `SOCKET_${type.toUpperCase()}`;
    if (Array.isArray(data)) {
      dispatch({ type, data });
    } else {
      dispatch({ type, ...data });
    }
  });
}
