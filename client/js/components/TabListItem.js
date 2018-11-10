import React, { memo } from 'react';
import classnames from 'classnames';

const TabListItem = ({
  target,
  content,
  server,
  selected,
  connected,
  onClick
}) => {
  const className = classnames({
    'tab-server': !target,
    success: !target && connected,
    error: !target && !connected,
    selected
  });

  return (
    <p className={className} onClick={() => onClick(server, target)}>
      <span className="tab-content">{content}</span>
    </p>
  );
};

export default memo(TabListItem);
