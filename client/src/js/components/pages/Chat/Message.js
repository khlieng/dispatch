import React, { PureComponent } from 'react';
import classnames from 'classnames';
import stringToHSL from 'utils/color';

export default class Message extends PureComponent {
  handleNickClick = () => this.props.onNickClick(this.props.message.from);

  render() {
    const { message } = this.props;

    const className = classnames('message', {
      [`message-${message.type}`]: message.type
    });

    const style = {
      paddingLeft: `${window.messageIndent + 15}px`,
      textIndent: `-${window.messageIndent}px`,
      ...this.props.style
    };

    const senderStyle = {};
    if (message.from) {
      senderStyle.color = stringToHSL(message.from);
    }

    return (
      <p className={className} style={style}>
        <span className="message-time">{message.time}</span>
        {message.from && (
          <span
            className="message-sender"
            style={senderStyle}
            onClick={this.handleNickClick}
          >
            {' '}
            {message.from}
          </span>
        )}{' '}
        {message.content}
      </p>
    );
  }
}
