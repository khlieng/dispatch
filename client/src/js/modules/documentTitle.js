import capitalize from 'lodash/capitalize';
import { getRouter } from 'state';
import { getCurrentServerName } from 'state/servers';
import { observe } from 'util/observe';

export default function documentTitle({ store }) {
  observe(store, [getRouter, getCurrentServerName], (router, serverName) => {
    let title;

    if (router.route === 'chat') {
      const { name } = router.params;
      if (name) {
        title = `${name} @ ${serverName}`;
      } else {
        title = serverName;
      }
    } else {
      title = capitalize(router.route);
    }

    document.title = `${title} | Dispatch`;
  });
}
