import Cookie from 'js-cookie';
import debounce from 'lodash/debounce';
import { getSelectedTab } from 'state/tab';
import { observe } from 'utils/observe';

const saveTab = debounce(
  tab => Cookie.set('tab', tab.toString(), { expires: 30 }),
  1000
);

export default function storage({ store }) {
  observe(store, getSelectedTab, tab => {
    if (tab.isChannel() || (tab.server && !tab.name)) {
      saveTab(tab);
    }
  });
}
