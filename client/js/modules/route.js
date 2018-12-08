import { updateSelection } from 'state/tab';
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
            store.dispatch(updateSelection(first));
            first = false;
          }
        }
      )
  );
}
