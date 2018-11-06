import React, { useState } from 'react';
import cn from 'classnames';

const AppInfo = ({ type, children, dismissible }) => {
  const [dismissed, setDismissed] = useState(false);

  if (!dismissed) {
    const handleDismiss = () => {
      if (dismissible) {
        setDismissed(true);
      }
    };

    const className = cn('app-info', {
      [`app-info-${type}`]: type
    });

    return (
      <div className={className} onClick={handleDismiss}>
        {children}
      </div>
    );
  }

  return null;
};

export default AppInfo;
