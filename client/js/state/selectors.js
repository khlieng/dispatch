import { createSelector } from 'reselect';
import get from 'lodash/get';
import { getServers } from './servers';
import { getSelectedTab } from './tab';

// eslint-disable-next-line import/prefer-default-export
export const getSelectedTabTitle = createSelector(
  getSelectedTab,
  getServers,
  (tab, servers) => tab.name || get(servers, [tab.server, 'name'])
);
