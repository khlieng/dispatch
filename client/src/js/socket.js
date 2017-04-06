import { routeActions } from 'react-router-redux';
import { broadcast, inform, addMessage, addMessages } from './actions/message';
import { select } from './actions/tab';
import { normalizeChannel } from './util';

function withReason(message, reason) {
  return message + (reason ? ` (${reason})` : '');
}

export default function handleSocket(socket, { dispatch, getState }) {
  socket.onAny((event, data) => {
    const type = `SOCKET_${event.toUpperCase()}`;
    if (Array.isArray(data)) {
      dispatch({ type, data });
    } else {
      dispatch({ type, ...data });
    }
  });

  socket.on('message', message => dispatch(addMessage(message)));
  socket.on('pm', message => dispatch(addMessage(message)));

  socket.on('join', data => {
    const state = getState();
    const { server, channel } = state.tab.selected;
    if (server && channel) {
      const { nick } = state.servers.get(server);
      const [joinedChannel] = data.channels;
      if (server === data.server &&
        nick === data.user &&
        channel !== joinedChannel &&
        normalizeChannel(channel) === normalizeChannel(joinedChannel)) {
        dispatch(select(server, joinedChannel));
      }
    }
  });

  socket.on('servers', data => {
    if (!data) {
      dispatch(routeActions.replace('/connect'));
    }
  });

  socket.on('join', ({ user, server, channels }) =>
    dispatch(inform(`${user} joined the channel`, server, channels[0]))
  );

  socket.on('part', ({ user, server, channel, reason }) =>
    dispatch(inform(withReason(`${user} left the channel`, reason), server, channel))
  );

  socket.on('quit', ({ user, server, reason, channels }) =>
    dispatch(broadcast(withReason(`${user} quit`, reason), server, channels))
  );

  socket.on('nick', data =>
    dispatch(broadcast(`${data.old} changed nick to ${data.new}`, data.server, data.channels))
  );

  socket.on('motd', ({ content, server }) =>
    dispatch(addMessages(content.map(line => ({
      server,
      to: server,
      message: line
    }))))
  );

  socket.on('whois', data => {
    const tab = getState().tab.selected;

    dispatch(inform([
      `Nick: ${data.nick}`,
      `Username: ${data.username}`,
      `Realname: ${data.realname}`,
      `Host: ${data.host}`,
      `Server: ${data.server}`,
      `Channels: ${data.channels}`
    ], tab.server, tab.channel));
  });

  socket.on('print', ({ server, message }) => dispatch(inform(message, server)));
}
