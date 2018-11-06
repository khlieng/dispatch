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
  store.dispatch(appSet('hexIP', env.hexIP));
  store.dispatch(setSettings(env.settings, true));

  if (env.servers) {
    store.dispatch({
      type: socketActions.SERVERS,
      data: env.servers
    });

    if (!store.getState().router.route) {
      const tab = Cookie.get('tab');
      if (tab) {
        const [server, name = null] = tab.split(/;(.+)/);

        if (
          name &&
          find(
            env.channels,
            chan => chan.server === server && chan.name === name
          )
        ) {
          store.dispatch(select(server, name, true));
        } else if (find(env.servers, srv => srv.host === server)) {
          store.dispatch(select(server, null, true));
        } else {
          store.dispatch(updateSelection());
        }
      } else {
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

export default function initialState(ctx) {
  if (window.__env__) {
    window.__env__.then(env => loadState(ctx, env));
  } else {
    const env = JSON.parse(document.getElementById('env').innerHTML);
    loadState(ctx, env);
  }
}
