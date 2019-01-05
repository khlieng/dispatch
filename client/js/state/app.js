import createReducer from 'utils/createReducer';
import * as actions from './actions';

export const getApp = state => state.app;
export const getConnected = state => state.app.connected;
export const getWrapWidth = state => state.app.wrapWidth;
export const getCharWidth = state => state.app.charWidth;
export const getWindowWidth = state => state.app.windowWidth;
export const getConnectDefaults = state => state.app.connectDefaults;

const initialState = {
  connected: true,
  wrapWidth: 0,
  charWidth: 0,
  windowWidth: 0,
  connectDefaults: {
    name: '',
    address: '',
    channels: [],
    ssl: false,
    password: false,
    readonly: false,
    showDetails: false
  },
  hexIP: false,
  newVersionAvailable: false,
  installable: null
};

export default createReducer(initialState, {
  [actions.APP_SET](state, { key, value }) {
    if (typeof key === 'object') {
      Object.assign(state, key);
    } else {
      state[key] = value;
    }
  },

  [actions.UPDATE_MESSAGE_HEIGHT](state, action) {
    state.wrapWidth = action.wrapWidth;
    state.charWidth = action.charWidth;
    state.windowWidth = action.windowWidth;
  }
});

export function appSet(key, value) {
  return {
    type: actions.APP_SET,
    key,
    value
  };
}

export function setConnected(connected) {
  return appSet('connected', connected);
}

export function setCharWidth(width) {
  return appSet('charWidth', width);
}
