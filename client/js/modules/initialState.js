/* eslint-disable no-underscore-dangle */

import { socket as socketActions } from 'state/actions';
import { getWrapWidth, setConnectDefaults, appSet } from 'state/app';
import { addMessages } from 'state/messages';
import { setSettings } from 'state/settings';
import { when } from 'utils/observe';

function loadState({ store }, env) {
  store.dispatch(setConnectDefaults(env.defaults));
  store.dispatch(setSettings(env.settings, true));

  if (env.servers) {
    store.dispatch({
      type: socketActions.SERVERS,
      data: env.servers
    });
  }

  if (env.channels) {
    store.dispatch({
      type: socketActions.CHANNELS,
      data: env.channels
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
      initialized: true,
      hexIP: env.hexIP,
      version: env.version
    })
  );

  // Wait until wrapWidth gets initialized so that height calculations
  // only happen once for these messages
  when(store, getWrapWidth, () => {
    if (env.messages) {
      const { messages, server, to, next } = env.messages;
      store.dispatch(addMessages(messages, server, to, false, next));
    }
  });
}

export default async function initialState(ctx) {
  const env = await window.__init__;
  ctx.socket.connect();
  loadState(ctx, env);
}
