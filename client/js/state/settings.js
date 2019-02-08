import createReducer from 'utils/createReducer';
import * as actions from './actions';

export const getSettings = state => state.settings;

export default createReducer(
  {},
  {
    [actions.UPLOAD_CERT](state) {
      state.uploadingCert = true;
    },

    [actions.socket.CERT_SUCCESS](state) {
      state.uploadingCert = false;
      delete state.certFile;
      delete state.cert;
      delete state.keyFile;
      delete state.key;
    },

    [actions.socket.CERT_FAIL](state, action) {
      state.uploadingCert = false;
      state.certError = action.message;
    },

    [actions.SET_CERT_ERROR](state, action) {
      state.uploadingCert = false;
      state.certError = action.message;
    },

    [actions.SET_CERT](state, action) {
      state.certFile = action.fileName;
      state.cert = action.cert;
    },

    [actions.SET_KEY](state, action) {
      state.keyFile = action.fileName;
      state.key = action.key;
    },

    [actions.SETTINGS_SET](state, { key, value, settings }) {
      if (settings) {
        Object.assign(state, settings);
      } else {
        state[key] = value;
      }
    }
  }
);

export function setCertError(message) {
  return {
    type: actions.SET_CERT_ERROR,
    message
  };
}

export function uploadCert() {
  return (dispatch, getState) => {
    const { settings } = getState();
    if (settings.cert && settings.key) {
      dispatch({
        type: actions.UPLOAD_CERT,
        socket: {
          type: 'cert',
          data: {
            cert: settings.cert,
            key: settings.key
          }
        }
      });
    } else {
      dispatch(setCertError('Missing certificate or key'));
    }
  };
}

export function setCert(fileName, cert) {
  return {
    type: actions.SET_CERT,
    fileName,
    cert
  };
}

export function setKey(fileName, key) {
  return {
    type: actions.SET_KEY,
    fileName,
    key
  };
}

export function setSetting(key, value) {
  return {
    type: actions.SETTINGS_SET,
    key,
    value,
    socket: {
      type: 'settings_set',
      data: {
        [key]: value
      },
      debounce: {
        delay: 250,
        key: `settings:${key}`
      }
    }
  };
}

export function setSettings(settings, local = false) {
  const action = {
    type: actions.SETTINGS_SET,
    settings
  };

  if (!local) {
    action.socket = {
      type: 'settings_set',
      data: settings
    };
  }

  return action;
}
