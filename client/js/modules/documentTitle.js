import capitalize from 'lodash/capitalize';
import { getRouter } from 'state';
import { getCurrentNetworkName } from 'state/networks';
import { observe } from 'utils/observe';

export default function documentTitle({ store }) {
  observe(store, [getRouter, getCurrentNetworkName], (router, networkName) => {
    let title;

    if (router.route === 'chat') {
      const { network, name } = router.params;
      if (name) {
        title = `${name} @ ${networkName || network}`;
      } else {
        title = networkName || network;
      }
    } else {
      title = capitalize(router.route);
    }

    if (title) {
      document.title = `${title} | Dispatch`;
    } else {
      document.title = 'Dispatch';
    }
  });
}
