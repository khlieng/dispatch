import { createSelector } from 'reselect';
import { getServers } from './servers';
import { getSelectedTab } from './tab';

// eslint-disable-next-line import/prefer-default-export
export const getSelectedTabTitle = createSelector(
  getSelectedTab,
  getServers,
  (tab, servers) => tab.name || servers.getIn([tab.server, 'name'])
);
