import Cookie from 'js-cookie';
import { socket as socketActions } from '../state/actions';
import { getWrapWidth, setConnectDefaults } from '../state/app';
import { addMessages } from '../state/messages';
import { select, updateSelection } from '../state/tab';
import { find } from '../util';
import { when } from '../util/observe';
import { replace } from '../util/router';

export default function initialState({ store }) {
  const env = JSON.parse(document.getElementById('env').innerHTML);

  store.dispatch(setConnectDefaults(env.defaults));

  if (env.servers) {
    store.dispatch({
      type: socketActions.SERVERS,
      data: env.servers
    });

    if (!store.getState().router.route) {
      const tab = Cookie.get('tab');
      if (tab) {
        const [server, name = null] = tab.split(':');

        if (find(env.servers, srv => srv.host === server)) {
          store.dispatch(select(server, name, true));
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
