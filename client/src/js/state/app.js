import { Record } from 'immutable';
import createReducer from 'utils/createReducer';
import * as actions from './actions';

export const getApp = state => state.app;
export const getConnected = state => state.app.connected;
export const getWrapWidth = state => state.app.wrapWidth;
export const getCharWidth = state => state.app.charWidth;
export const getWindowWidth = state => state.app.windowWidth;
export const getConnectDefaults = state => state.app.connectDefaults;

const ConnectDefaults = Record({
  name: '',
  address: '',
  channels: [],
  ssl: false,
  password: false,
  readonly: false,
  showDetails: false
});

const App = Record({
  connected: true,
  wrapWidth: 0,
  charWidth: 0,
  windowWidth: 0,
  connectDefaults: new ConnectDefaults()
});

export default createReducer(new App(), {
  [actions.APP_SET](state, action) {
    return state.set(action.key, action.value);
  },

  [actions.UPDATE_MESSAGE_HEIGHT](state, action) {
    return state
      .set('wrapWidth', action.wrapWidth)
      .set('charWidth', action.charWidth)
      .set('windowWidth', action.windowWidth);
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

export function setConnectDefaults(defaults) {
  return appSet('connectDefaults', new ConnectDefaults(defaults));
}
