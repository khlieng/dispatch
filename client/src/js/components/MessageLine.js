import React, { Component } from 'react';
import Autolinker from 'autolinker';
import pure from 'pure-render-decorator';

@pure
export default class MessageLine extends Component {
  render() {
    const line = Autolinker.link(this.props.line, { stripPrefix: false });

    let messageClass = 'message';
    if (this.props.type) {
      messageClass += ' message-' + this.props.type;
    }

    const style = {
      paddingLeft: window.messageIndent + 'px'
    };

    return (
      <p className={messageClass} style={style}>
        <span dangerouslySetInnerHTML={{ __html: line }}></span>
      </p>
    );
  }
}
