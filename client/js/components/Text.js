import React from 'react';
import stringToRGB from 'utils/color';

function nickStyle(nick, color) {
  const style = {
    fontWeight: 400
  };

  if (color) {
    style.color = stringToRGB(nick);
  }

  return style;
}

function renderBlock(block, coloredNick, key) {
  switch (block.type) {
    case 'text':
      return block.text;

    case 'link':
      return (
        <a target="_blank" rel="noopener noreferrer" href={block.url} key={key}>
          {block.text}
        </a>
      );

    case 'format':
      return (
        <span style={block.style} key={key}>
          {block.text}
        </span>
      );

    case 'nick':
      return (
        <span
          className="message-sender"
          style={nickStyle(block.text, coloredNick)}
          key={key}
        >
          {block.text}
        </span>
      );

    case 'events':
      return (
        <span className="message-events-more" key={key}>
          {block.text}
        </span>
      );

    default:
      return null;
  }
}

const Text = ({ children, coloredNick }) => {
  if (!children) {
    return null;
  }
  if (children.length > 1) {
    let key = 0;
    return children.map(block => renderBlock(block, coloredNick, key++));
  }
  if (children.length === 1) {
    return renderBlock(children[0], coloredNick);
  }
  return children;
};

export default Text;
