import React, { memo } from 'react';
import classnames from 'classnames';
import stringToRGB from 'utils/color';

const Message = ({ message, coloredNick, style, onNickClick }) => {
  const className = classnames('message', {
    [`message-${message.type}`]: message.type
  });

  if (message.type === 'date') {
    return (
      <div className={className} style={style}>
        {message.content}
        <hr />
      </div>
    );
  }

  style = {
    ...style,
    paddingLeft: `${window.messageIndent + 15}px`,
    textIndent: `-${window.messageIndent}px`
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
