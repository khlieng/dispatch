import React from 'react';
import classnames from 'classnames';

const TabListItem = ({
  target,
  content,
  server,
  selected,
  connected,
  joined,
  error,
  onClick
}) => {
  const className = classnames({
    'tab-server': !target,
    success: !target && connected,
    error: (!target && !connected) || (!joined && error),
    disabled: !!target && !error && joined === false,
    selected
  });

  return (
    <p className={className} onClick={() => onClick(server, target)}>
      <span className="tab-content">{content}</span>
    </p>
  );
};

export default TabListItem;
