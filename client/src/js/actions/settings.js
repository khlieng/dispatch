import base64 from 'base64-arraybuffer';
import * as actions from '../actions';

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
