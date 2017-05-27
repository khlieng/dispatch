import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';
import Settings from '../components/pages/Settings';
import { getSettings, setCert, setKey, uploadCert } from '../state/settings';

const mapState = createStructuredSelector({
  settings: getSettings
});

const mapDispatch = dispatch => ({
  onCertChange(name, data) { dispatch(setCert(name, data)); },
  onKeyChange(name, data) { dispatch(setKey(name, data)); },
  ...bindActionCreators({ uploadCert }, dispatch)
});

export default connect(mapState, mapDispatch)(Settings);
