import Cookie from 'js-cookie';
import debounce from 'lodash/debounce';
import { getSelectedTab } from 'state/tab';
import { stringifyTab } from 'utils';
import { observe } from 'utils/observe';

const saveTab = debounce(
  tab => Cookie.set('tab', stringifyTab(tab), { expires: 30 }),
  1000
);

export default function storage({ store }) {
  observe(store, getSelectedTab, tab => {
    if (tab.network) {
      saveTab(tab);
    }
  });
}
