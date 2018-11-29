/* eslint-disable no-underscore-dangle */
import Cookie from 'js-cookie';
import { socket as socketActions } from 'state/actions';
import { getWrapWidth, setConnectDefaults, appSet } from 'state/app';
import { addMessages } from 'state/messages';
import { setSettings } from 'state/settings';
import { select, updateSelection } from 'state/tab';
import { find } from 'utils';
import { when } from 'utils/observe';
import { replace } from 'utils/router';

function loadState({ store }, env) {
  store.dispatch(setConnectDefaults(env.defaults));
  store.dispatch(
    appSet({
      hexIP: env.hexIP,
      version: env.version
    })
  );
  store.dispatch(setSettings(env.settings, true));

  if (env.servers) {
    store.dispatch({
      type: socketActions.SERVERS,
      data: env.servers
    });

    const { router } = store.getState();

    if (!router.route || router.route === 'chat') {
      const tabs = [];

      if (router.route === 'chat') {
        tabs.push(router.params);
      }

      const cookie = Cookie.get('tab');
      if (cookie) {
        const [server, name = null] = cookie.split(/;(.+)/);
        tabs.push({
          server,
          name
        });
      }

      let found = false;
      let i = 0;

      while (!found) {
        const tab = tabs[i];
        i++;

        if (
          tab.name &&
          find(
            env.channels,
            chan => chan.server === tab.server && chan.name === tab.name
          )
        ) {
          found = true;
          store.dispatch(select(tab.server, tab.name, true));
        } else if (find(env.servers, srv => srv.host === tab.server)) {
          found = true;
          store.dispatch(select(tab.server, null, true));
        }
      }

      if (!found) {
        store.dispatch(updateSelection());
      }
    }
  } else {
    store.dispatch(replace('/connect'));
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
