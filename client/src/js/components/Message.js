import React, { PureComponent } from 'react';

export default class Message extends PureComponent {
  handleSenderClick = () => {
    const { message, openPrivateChat, select } = this.props;

    openPrivateChat(message.server, message.from);
    select(message.server, message.from, true);
  };

  render() {
    const { message } = this.props;
    const className = message.type ? `message message-${message.type}` : 'message';
    const style = {
      paddingLeft: `${window.messageIndent + 15}px`,
      textIndent: `-${window.messageIndent}px`,
      ...this.props.style
    };

    return (
      <p className={className} style={style}>
        <span className="message-time">{message.time}</span>
        {message.from &&
          <span className="message-sender" onClick={this.handleSenderClick}>
            {' '}{message.from}
          </span>
        }{' '}{message.message}
      </p>
    );
  }
}
