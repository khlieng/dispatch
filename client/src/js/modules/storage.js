import Cookie from 'js-cookie';
import debounce from 'lodash/debounce';
import observe from '../util/observe';
import { getSelectedTab } from '../reducers/tab';

const saveTab = debounce(tab =>
  Cookie.set('tab', tab.toString(), { expires: 30 })
, 3000);

export default function storage({ store }) {
  observe(store, getSelectedTab, tab => {
    if (tab.server) {
      saveTab(tab);
    }
  });
}
