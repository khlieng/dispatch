import { createSelector } from 'reselect';
import get from 'lodash/get';
import { getNetworks } from './networks';
import { getSelectedTab } from './tab';

// eslint-disable-next-line import/prefer-default-export
export const getSelectedTabTitle = createSelector(
  getSelectedTab,
  getNetworks,
  (tab, networks) => tab.name || get(networks, [tab.network, 'name'])
);
