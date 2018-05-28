import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';
import Settings from 'components/pages/Settings';
import { getSettings, setCert, setKey, uploadCert } from 'state/settings';

const mapState = createStructuredSelector({
  settings: getSettings
});

const mapDispatch = {
  onCertChange: setCert,
  onKeyChange: setKey,
  uploadCert
};

export default connect(
  mapState,
  mapDispatch
)(Settings);
