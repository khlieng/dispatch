import { Map } from 'immutable';
import createReducer from '../util/createReducer';
import * as actions from '../actions';

export default createReducer(Map(), {
  [actions.UPLOAD_CERT](state) {
    return state.set('uploadingCert', true);
  },

  [actions.SOCKET_CERT_SUCCESS]() {
    return Map({ uploadingCert: false });
  },

  [actions.SOCKET_CERT_FAIL](state, action) {
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
