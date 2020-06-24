import { INIT } from 'state/actions';
import { getConnected, getWrapWidth } from 'state/app';
import { searchChannels } from 'state/channelSearch';
import { addMessages } from 'state/messages';
import { when } from 'utils/observe';

function loadState({ store }, env) {
  store.dispatch({
    type: INIT,
    settings: env.settings,
    networks: env.networks,
    channels: env.channels,
    openDMs: env.openDMs,
    users: env.users,
    app: {
      connectDefaults: env.defaults,
      initialized: true,
      hexIP: env.hexIP,
      version: env.version
    }
  });

  if (env.messages) {
    // Wait until wrapWidth gets initialized so that height calculations
    // only happen once for these messages
    when(store, getWrapWidth, () => {
      const { messages, network, to, next } = env.messages;
      store.dispatch(addMessages(messages, network, to, false, next));
    });
  }

  if (env.networks) {
    when(store, getConnected, () =>
      // Cache top channels for each network
      env.networks.forEach(({ host }) =>
        store.dispatch(searchChannels(host, ''))
      )
    );
  }
}

/* eslint-disable no-underscore-dangle */
export default async function initialState(ctx) {
  const env = await window.__init__;
  ctx.socket.connect();
  loadState(ctx, env);
}
