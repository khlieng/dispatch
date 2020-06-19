import React from 'react';
import classnames from 'classnames';

function splitContent(content) {
  let start = 0;
  while (content[start] === '#') {
    start++;
  }

  if (start > 0) {
    return [content.slice(0, start), content.slice(start)];
  }
  return [null, content];
}

const TabListItem = ({
  target,
  content,
  network,
  selected,
  connected,
  joined,
  error,
  onClick
}) => {
  const className = classnames({
    'tab-network': !target,
    success: !target && connected,
    error: (!target && !connected) || (!joined && error),
    disabled: !!target && !error && joined === false,
    selected
  });

  const [prefix, name] = splitContent(content);

  return (
    <p className={className} onClick={() => onClick(network, target)}>
      <span className="tab-content">
        {prefix && <span className="tab-prefix">{prefix}</span>}
        {name}
      </span>
    </p>
  );
};

export default TabListItem;
