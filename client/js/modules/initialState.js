import { socket as socketActions } from 'state/actions';
import { getConnected, getWrapWidth, appSet } from 'state/app';
import { searchChannels } from 'state/channelSearch';
import { addMessages } from 'state/messages';
import { setSettings } from 'state/settings';
import { when } from 'utils/observe';

function loadState({ store }, env) {
  store.dispatch(setSettings(env.settings, true));

  if (env.servers) {
    store.dispatch({
      type: socketActions.SERVERS,
      data: env.servers
    });

    when(store, getConnected, () =>
      // Cache top channels for each server
      env.servers.forEach(({ host }) =>
        store.dispatch(searchChannels(host, ''))
      )
    );
  }

  if (env.channels) {
    store.dispatch({
      type: socketActions.CHANNELS,
      data: env.channels
    });
  }

  if (env.openDMs) {
    store.dispatch({
      type: 'PRIVATE_CHATS',
      privateChats: env.openDMs
    });
  }

  if (env.users) {
    store.dispatch({
      type: socketActions.USERS,
      ...env.users
    });
  }

  store.dispatch(
    appSet({
      connectDefaults: env.defaults,
      initialized: true,
      hexIP: env.hexIP,
      version: env.version
    })
  );

  if (env.messages) {
    // Wait until wrapWidth gets initialized so that height calculations
    // only happen once for these messages
    when(store, getWrapWidth, () => {
      const { messages, server, to, next } = env.messages;
      store.dispatch(addMessages(messages, server, to, false, next));
    });
  }
}

/* eslint-disable no-underscore-dangle */
export default async function initialState(ctx) {
  const env = await window.__init__;
  ctx.socket.connect();
  loadState(ctx, env);
}
