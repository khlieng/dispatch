import { socket as socketActions } from 'state/actions';
import { getWrapWidth, appSet } from 'state/app';
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
