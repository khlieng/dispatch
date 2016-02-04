import React, { Component } from 'react';
import Autolinker from 'autolinker';
import pure from 'pure-render-decorator';

@pure
export default class MessageLine extends Component {
  render() {
    const { line, type } = this.props;
    const content = Autolinker.link(line, { stripPrefix: false });
    const classes = ['message'];

    if (type) {
      classes.push(`message-${type}`);
    }

    const style = {
      paddingLeft: `${window.messageIndent}px`
    };

    return (
      <p className={classes.join(' ')} style={style}>
        <span dangerouslySetInnerHTML={{ __html: content }}></span>
      </p>
    );
  }
}
