import React, { Component } from 'react';
import { connect } from 'react-redux';
import pure from 'pure-render-decorator';
import Navicon from '../components/Navicon';
import FileInput from '../components/FileInput';
import { setCert, setKey, uploadCert } from '../actions/settings';

@pure
class Settings extends Component {
  handleCertChange = (name, data) => this.props.dispatch(setCert(name, data));
  handleKeyChange = (name, data) => this.props.dispatch(setKey(name, data));
  handleCertUpload = () => this.props.dispatch(uploadCert());

  render() {
    const { settings } = this.props;
    const status = settings.get('uploadingCert') ? 'Uploading...' : 'Upload';
    const error = settings.get('certError');

    return (
      <div className="settings">
        <Navicon />
        <h1>Settings</h1>
        <h2>Client Certificate</h2>
        <div>
          <p>Certificate</p>
          <FileInput
            name={settings.get('certFile') || 'Select Certificate'}
            onChange={this.handleCertChange}
          />
        </div>
        <div>
          <p>Private Key</p>
          <FileInput
            name={settings.get('keyFile') || 'Select Key'}
            onChange={this.handleKeyChange}
          />
        </div>
        <button onClick={this.handleCertUpload}>{status}</button>
        { error ? <p className="error">{error}</p> : null }
      </div>
    );
  }
}

export default connect(state => ({
  settings: state.settings
}))(Settings);
