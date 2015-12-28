import React, { Component } from 'react';
import Autolinker from 'autolinker';
import { timestamp } from '../util';

export default class MessageHeader extends Component {
  shouldComponentUpdate(nextProps) {
    return nextProps.message.lines[0] !== this.props.message.lines[0];
  }

  handleSenderClick = () => {
    const { message, openPrivateChat, select } = this.props;

    openPrivateChat(message.server, message.from);
    select(message.server, message.from, true);
  }

  render() {
    const { message } = this.props;
    const line = Autolinker.link(message.lines[0], { stripPrefix: false });
    let sender = null;
    let messageClass = 'message';

    if (message.from) {
      sender = (
        <span>
          {' '}
          <span className="message-sender" onClick={this.handleSenderClick}>
            {message.from}
          </span>
        </span>
      );
    }

    if (message.type) {
      messageClass += ' message-' + message.type;
    }

    return (
      <p className={messageClass}>
        <span className="message-time">{timestamp(message.time)}</span>
        {sender}
        <span dangerouslySetInnerHTML={{ __html: ' ' + line }}></span>
      </p>
    );
  }
}
