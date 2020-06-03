import { socketAction } from 'state/actions';
import { setConnected } from 'state/app';
import {
  print,
  addMessage,
  addMessages,
  addEvent,
  broadcastEvent
} from 'state/messages';
import { openModal } from 'state/modals';
import { reconnect } from 'state/servers';
import { select } from 'state/tab';
import { find } from 'utils';

function findChannels(state, server, user) {
  const channels = [];

  Object.keys(state.channels[server]).forEach(channel => {
    if (find(state.channels[server][channel].users, u => u.nick === user)) {
      channels.push(channel);
    }
  });

  return channels;
}

export default function handleSocket({
  socket,
  store: { dispatch, getState }
}) {
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
      dispatch(addEvent(server, channels[0], 'join', user));
    },

    part({ user, server, channel, reason }) {
      dispatch(addEvent(server, channel, 'part', user, reason));
    },

    quit({ user, server, reason }) {
      const channels = findChannels(getState(), server, user);
      dispatch(broadcastEvent(server, channels, 'quit', user, reason));
    },

    nick({ server, oldNick, newNick }) {
      if (oldNick) {
        const channels = findChannels(getState(), server, oldNick);
        dispatch(broadcastEvent(server, channels, 'nick', oldNick, newNick));
      }
    },

    topic({ server, channel, topic, nick }) {
      if (nick) {
        dispatch(addEvent(server, channel, 'topic', nick, topic));
        /* if (topic) {
          dispatch(inform(`${nick} changed the topic to:`, server, channel));
          dispatch(print(topic, server, channel));
        } else {
          dispatch(inform(`${nick} cleared the topic`, server, channel));
        } */
      }
    },

    motd({ content, server }) {
      dispatch(
        addMessages(
          content.map(line => ({ content: line })),
          server
        )
      );
    },

    whois(data) {
      const tab = getState().tab.selected;

      dispatch(
        print(
          [
            `Nick: ${data.nick}`,
            `Username: ${data.username}`,
            `Realname: ${data.realname}`,
            `Host: ${data.host}`,
            `Server: ${data.server}`,
            `Channels: ${data.channels}`
          ],
          tab.server,
          tab.name
        )
      );
    },

    print(message) {
      const tab = getState().tab.selected;
      dispatch(addMessage(message, tab.server, tab.name));
    },

    error({ server, target, message }) {
      dispatch(addMessage({ content: message, type: 'error' }, server, target));
    },

    connection_update({ server, errorType }) {
      if (errorType === 'verify') {
        dispatch(
          openModal('confirm', {
            question:
              'The server is using a self-signed certificate, continue anyway?',
            onConfirm: () =>
              dispatch(
                reconnect(server, {
                  skipVerify: true
                })
              )
          })
        );
      }
    },

    dcc_send({ server, from, filename, url }) {
      const serverName = getState().servers[server]?.name || server;

      dispatch(
        openModal('confirm', {
          question: `${from} on ${serverName} is sending you: ${filename}`,
          confirmation: 'Download',
          onConfirm: () => {
            const a = document.createElement('a');
            a.href = url;
            a.click();
          }
        })
      );
    },

    _connected(connected) {
      dispatch(setConnected(connected));
    }
  };

  const afterHandlers = {
    channel_forward(forward) {
      const { selected } = getState().tab;

      if (selected.server === forward.server && selected.name === forward.old) {
        dispatch(select(forward.server, forward.new, true));
      }
    }
  };

  socket.onMessage((type, data) => {
    let action;
    if (Array.isArray(data)) {
      action = { type: socketAction(type), data: [...data] };
    } else {
      action = { ...data, type: socketAction(type) };
    }

    if (type in handlers) {
      handlers[type](data);
    }

    if (type.charAt(0) === '_') {
      return;
    }

    dispatch(action);

    if (type in afterHandlers) {
      afterHandlers[type](data);
    }
  });
}
