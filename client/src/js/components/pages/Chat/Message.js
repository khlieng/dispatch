import React, { PureComponent } from 'react';

export default class Message extends PureComponent {
  handleNickClick = () => this.props.onNickClick(this.props.message.from);

  render() {
    const { message } = this.props;
    const className = message.type
      ? `message message-${message.type}`
      : 'message';
    const style = {
      paddingLeft: `${window.messageIndent + 15}px`,
      textIndent: `-${window.messageIndent}px`,
      ...this.props.style
    };

    return (
      <p className={className} style={style}>
        <span className="message-time">{message.time}</span>
        {message.from && (
          <span className="message-sender" onClick={this.handleNickClick}>
            {' '}
            {message.from}
          </span>
        )}{' '}
        {message.content}
      </p>
    );
  }
}
