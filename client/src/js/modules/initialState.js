import { setEnvironment } from '../actions/environment';
import { addMessages } from '../actions/message';
import { select, updateSelection } from '../actions/tab';
import { find } from '../util';
import { initWidthUpdates } from '../util/messageHeight';
import { replace } from '../util/router';

export default function initialState({ store }) {
  const env = JSON.parse(document.getElementById('env').innerHTML);

  store.dispatch(setEnvironment('connect_defaults', env.defaults));

  if (env.servers) {
    store.dispatch({
      type: 'SOCKET_SERVERS',
      data: env.servers
    });

    if (!store.getState().router.route) {
      let tab = localStorage.tab;
      if (tab) {
        tab = JSON.parse(tab);

        if (find(env.servers, server => server.host === tab.server)) {
          store.dispatch(select(tab.server, tab.name, true));
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
      type: 'SOCKET_CHANNELS',
      data: env.channels
    });
  }

  if (env.users) {
    store.dispatch({
      type: 'SOCKET_USERS',
      ...env.users
    });
  }

  initWidthUpdates(store, () => {
    if (env.messages) {
      const { messages, server, to, next } = env.messages;
      store.dispatch(addMessages(messages, server, to, false, next));
    }
  });
}
