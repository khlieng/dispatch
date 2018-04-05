import { Map } from 'immutable';
import base64 from 'base64-arraybuffer';
import createReducer from 'utils/createReducer';
import * as actions from './actions';

export const getSettings = state => state.settings;

export default createReducer(Map(), {
  [actions.UPLOAD_CERT](state) {
    return state.set('uploadingCert', true);
  },

  [actions.socket.CERT_SUCCESS]() {
    return Map({ uploadingCert: false });
  },

  [actions.socket.CERT_FAIL](state, action) {
    return state.merge({
      uploadingCert: false,
      certError: action.message
    });
  },

  [actions.SET_CERT_ERROR](state, action) {
    return state.merge({
      uploadingCert: false,
      certError: action.message
    });
  },

  [actions.SET_CERT](state, action) {
    return state.merge({
      certFile: action.fileName,
      cert: action.cert
    });
  },

  [actions.SET_KEY](state, action) {
    return state.merge({
      keyFile: action.fileName,
      key: action.key
    });
  }
});

export function setCertError(message) {
  return {
    type: actions.SET_CERT_ERROR,
    message
  };
}

export function uploadCert() {
  return (dispatch, getState) => {
    const { settings } = getState();
    if (settings.has('cert') && settings.has('key')) {
      dispatch({
        type: actions.UPLOAD_CERT,
        socket: {
          type: 'cert',
          data: {
            cert: settings.get('cert'),
            key: settings.get('key')
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
    cert: base64.encode(cert)
  };
}

export function setKey(fileName, key) {
  return {
    type: actions.SET_KEY,
    fileName,
    key: base64.encode(key)
  };
}
