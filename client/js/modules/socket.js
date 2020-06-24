import { socketAction } from 'state/actions';
import { kicked } from 'state/channels';
import {
  print,
  addMessage,
  addMessages,
  addEvent,
  broadcastEvent
} from 'state/messages';
import { openModal } from 'state/modals';
import { reconnect } from 'state/networks';
import { select } from 'state/tab';
import { find } from 'utils';

function findChannels(state, network, user) {
  const channels = [];

  Object.keys(state.channels[network]).forEach(channel => {
    if (find(state.channels[network][channel].users, u => u.nick === user)) {
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
      dispatch(addMessage(message, message.network, message.to));
      return false;
    },

    pm(message) {
      dispatch(addMessage(message, message.network, message.from));
      return false;
    },

    messages({ messages, network, to, prepend, next }) {
      dispatch(addMessages(messages, network, to, prepend, next));
      return false;
    },

    join({ user, network, channels }) {
      dispatch(addEvent(network, channels[0], 'join', user));
    },

    part({ user, network, channel, reason }) {
      dispatch(addEvent(network, channel, 'part', user, reason));
    },

    quit({ user, network, reason }) {
      const channels = findChannels(getState(), network, user);
      dispatch(broadcastEvent(network, channels, 'quit', user, reason));
    },

    kick({ network, channel, sender, user, reason }) {
      dispatch(kicked(network, channel, user));
      dispatch(addEvent(network, channel, 'kick', user, sender, reason));
    },

    nick({ network, oldNick, newNick }) {
      if (oldNick) {
        const channels = findChannels(getState(), network, oldNick);
        dispatch(broadcastEvent(network, channels, 'nick', oldNick, newNick));
      }
    },

    topic({ network, channel, topic, nick }) {
      if (nick) {
        dispatch(addEvent(network, channel, 'topic', nick, topic));
      }
    },

    motd({ content, network }) {
      dispatch(
        addMessages(
          content.map(line => ({ content: line })),
          network
        )
      );
      return false;
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
          tab.network,
          tab.name
        )
      );
      return false;
    },

    print(message) {
      const tab = getState().tab.selected;
      dispatch(addMessage(message, tab.network, tab.name));
      return false;
    },

    error({ network, target, message }) {
      const state = getState();
      const tab = state.tab.selected;

      if (network === tab.network) {
        // Print it in the current channel if the error happened on
        // the current network
        target = tab.name;
      } else if (!state.channels[network]?.[target]) {
        // Print it the network tab if the target does not exist
        target = null;
      }

      dispatch(
        addMessage({ content: message, type: 'error' }, network, target)
      );
      return false;
    },

    connection_update({ network, errorType }) {
      if (errorType === 'verify') {
        dispatch(
          openModal('confirm', {
            question:
              'The network is using a self-signed certificate, continue anyway?',
            onConfirm: () =>
              dispatch(
                reconnect(network, {
                  skipVerify: true
                })
              )
          })
        );
      }
    },

    dcc_send({ network, from, filename, size, url }) {
      const networkName = getState().networks[network]?.name || network;

      dispatch(
        openModal('confirm', {
          question: `${from} on ${networkName} is sending you (${size}): ${filename}`,
          confirmation: 'Download',
          onConfirm: () => {
            const a = document.createElement('a');
            a.href = url;
            a.click();
          }
        })
      );
    }
  };

  const afterHandlers = {
    channel_forward(forward) {
      const { selected } = getState().tab;

      if (
        selected.network === forward.network &&
        selected.name === forward.old
      ) {
        dispatch(select(forward.network, forward.new, true));
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

    if (handlers[type]?.(data) === false) {
      return;
    }

    dispatch(action);

    afterHandlers[type]?.(data);
  });
}
