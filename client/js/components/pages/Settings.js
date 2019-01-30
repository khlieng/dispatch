import React, { useCallback } from 'react';
import Navicon from 'containers/Navicon';
import Button from 'components/ui/Button';
import Checkbox from 'components/ui/Checkbox';
import FileInput from 'components/ui/FileInput';

const Settings = ({
  settings,
  installable,
  version,
  setSetting,
  onCertChange,
  onKeyChange,
  onInstall,
  uploadCert
}) => {
  const status = settings.uploadingCert ? 'Uploading...' : 'Upload';
  const error = settings.certError;

  const handleInstallClick = useCallback(async () => {
    installable.prompt();
    await installable.userChoice;
    onInstall();
  }, [installable]);

  return (
    <div className="settings-container">
      <div className="settings">
        <Navicon />
        <h1>Settings</h1>
        {installable && (
          <Button className="button-install" onClick={handleInstallClick}>
            <h2>Install</h2>
          </Button>
        )}
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
            <Button
              type="submit"
              className="settings-button"
              onClick={uploadCert}
            >
              {status}
            </Button>
            {error ? <p className="error">{error}</p> : null}
          </div>
        </div>
        {version && (
          <div className="settings-version">
            <p>{version.tag}</p>
            <p>Commit: {version.commit}</p>
            <p>Build Date: {version.date}</p>
          </div>
        )}
      </div>
    </div>
  );
};

export default Settings;
