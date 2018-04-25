import React from 'react';
import Navicon from 'containers/Navicon';
import FileInput from 'components/ui/FileInput';

const Settings = ({ settings, onCertChange, onKeyChange, uploadCert }) => {
  const status = settings.uploadingCert ? 'Uploading...' : 'Upload';
  const error = settings.certError;

  return (
    <div className="settings">
      <Navicon />
      <h1>Settings</h1>
      <h2>Client Certificate</h2>
      <div>
        <p>Certificate</p>
        <FileInput
          name={settings.certFile || 'Select Certificate'}
          onChange={onCertChange}
        />
      </div>
      <div>
        <p>Private Key</p>
        <FileInput
          name={settings.keyFile || 'Select Key'}
          onChange={onKeyChange}
        />
      </div>
      <button onClick={uploadCert}>{status}</button>
      {error ? <p className="error">{error}</p> : null}
    </div>
  );
};

export default Settings;
