import React from 'react';
import Navicon from 'containers/Navicon';
import Checkbox from 'components/ui/Checkbox';
import FileInput from 'components/ui/FileInput';

const Settings = ({
  settings,
  setSetting,
  onCertChange,
  onKeyChange,
  uploadCert
}) => {
  const status = settings.uploadingCert ? 'Uploading...' : 'Upload';
  const error = settings.certError;

  return (
    <div className="settings-container">
      <div className="settings">
        <Navicon />
        <h1>Settings</h1>
        <div className="settings-section">
          <h2>Visuals</h2>
          <Checkbox
            name="coloredNicks"
            label="Colored nicks"
            checked={settings.coloredNicks}
            onChange={e => setSetting('coloredNicks', e.target.checked)}
          />
        </div>
        <div className="settings-section">
          <h2>Client Certificate</h2>
          <div className="settings-cert">
            <div className="settings-file">
              <p>Certificate</p>
              <FileInput
                name={settings.certFile || 'Select Certificate'}
                onChange={onCertChange}
              />
            </div>
            <div className="settings-file">
              <p>Private Key</p>
              <FileInput
                name={settings.keyFile || 'Select Key'}
                onChange={onKeyChange}
              />
            </div>
            <button className="settings-button" onClick={uploadCert}>
              {status}
            </button>
            {error ? <p className="error">{error}</p> : null}
          </div>
        </div>
      </div>
    </div>
  );
};

export default Settings;
