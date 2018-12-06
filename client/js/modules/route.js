import Cookie from 'js-cookie';
import { select, updateSelection, tabExists } from 'state/tab';
import { observe, when } from 'utils/observe';

export default function route({ store }) {
  let first = true;

  when(
    store,
    state => state.app.initialized,
    () =>
      observe(
        store,
        state => state.router,
        router => {
          if (!router.route || router.route === 'chat') {
            const state = store.getState();
            let redirect = true;
            const tabs = [];

            if (router.route === 'chat') {
              if (tabExists(router.params, state)) {
                redirect = false;
              } else {
                tabs.push(router.params);
              }
            }

            if (redirect && first) {
              const cookie = Cookie.get('tab');
              if (cookie) {
                const [server, name = null] = cookie.split(/;(.+)/);
                tabs.unshift({ server, name });
              }
            }

            if (redirect) {
              let found = false;

              for (let i = 0; i < tabs.length; i++) {
                const tab = tabs[i];
                if (tabExists(tab, state)) {
                  store.dispatch(select(tab.server, tab.name, true));
                  found = true;
                  break;
                }
              }

              if (!found) {
                store.dispatch(updateSelection());
              }
            }

            first = false;
          }
        }
      )
  );
}
