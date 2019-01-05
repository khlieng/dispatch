import { createSelector } from 'reselect';
import createReducer from 'utils/createReducer';
import * as actions from './actions';

export const getModals = state => state.modals;

export const getHasOpenModals = createSelector(
  getModals,
  modals => {
    const keys = Object.keys(modals);

    for (let i = 0; i < keys.length; i++) {
      if (modals[keys[i]].isOpen) {
        return true;
      }
    }
    return false;
  }
);

export default createReducer(
  {},
  {
    [actions.OPEN_MODAL](state, { name, payload = {} }) {
      state[name] = {
        isOpen: true,
        payload
      };
    },

    [actions.CLOSE_MODAL](state, { name }) {
      state[name].isOpen = false;
    }
  }
);

export function openModal(name, payload) {
  return {
    type: actions.OPEN_MODAL,
    name,
    payload
  };
}

export function closeModal(name) {
  return {
    type: actions.CLOSE_MODAL,
    name
  };
}
