import { createStructuredSelector } from 'reselect';
import Settings from 'components/pages/Settings';
import { appSet } from 'state/app';
import {
  getSettings,
  setSetting,
  setCert,
  setKey,
  uploadCert
} from 'state/settings';
import connect from 'utils/connect';

const mapState = createStructuredSelector({
  settings: getSettings,
  installable: state => state.app.installable,
  version: state => state.app.version
});

const mapDispatch = {
  onCertChange: setCert,
  onKeyChange: setKey,
  uploadCert,
  setSetting,
  onInstall: () => appSet('installable', null)
};

export default connect(
  mapState,
  mapDispatch
)(Settings);
