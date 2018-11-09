import { createStructuredSelector } from 'reselect';
import Settings from 'components/pages/Settings';
import {
  getSettings,
  setSetting,
  setCert,
  setKey,
  uploadCert
} from 'state/settings';
import connect from 'utils/connect';

const mapState = createStructuredSelector({
  settings: getSettings
});

const mapDispatch = {
  onCertChange: setCert,
  onKeyChange: setKey,
  uploadCert,
  setSetting
};

export default connect(
  mapState,
  mapDispatch
)(Settings);
