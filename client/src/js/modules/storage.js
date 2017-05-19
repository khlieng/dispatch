import debounce from 'lodash/debounce';
import observe from '../util/observe';
import { getSelectedTab } from '../reducers/tab';

const saveTab = debounce(tab => {
  localStorage.tab = JSON.stringify(tab);
}, 3000);

export default function storage({ store }) {
  observe(store, getSelectedTab, tab => {
    if (tab.server) {
      saveTab(tab);
    }
  });
}
