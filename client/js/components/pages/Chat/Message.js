import React, { memo } from 'react';
import classnames from 'classnames';
import stringToRGB from 'utils/color';

const Message = ({ message, coloredNick, onNickClick }) => {
  const className = classnames('message', {
    [`message-${message.type}`]: message.type
  });

  if (message.type === 'date') {
    return (
      <div className={className}>
        {message.content}
        <hr />
      </div>
    );
  }

  const style = {
    paddingLeft: `${message.indent + 15}px`,
    textIndent: `-${message.indent}px`
  };

  const senderStyle = {};
  if (message.from && coloredNick) {
    senderStyle.color = stringToRGB(message.from);
  }

  return (
    <p className={className} style={style}>
      <span className="message-time">{message.time} </span>
      {message.from && (
        <span
          className="message-sender"
          style={senderStyle}
          onClick={() => onNickClick(message.from)}
        >
          {message.from}
        </span>
      )}
      <span> {message.content}</span>
    </p>
  );
};

export default memo(Message);
