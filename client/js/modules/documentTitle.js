import capitalize from 'lodash/capitalize';
import { getRouter } from 'state';
import { getCurrentServerName } from 'state/servers';
import { observe } from 'utils/observe';

export default function documentTitle({ store }) {
  observe(store, [getRouter, getCurrentServerName], (router, serverName) => {
    let title;

    if (router.route === 'chat') {
      const { server, name } = router.params;
      if (name) {
        title = `${name} @ ${serverName || server}`;
      } else {
        title = serverName || server;
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
