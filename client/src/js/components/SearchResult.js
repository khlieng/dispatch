import React, { Component } from 'react';
import Autolinker from 'autolinker';
import pure from 'pure-render-decorator';
import { timestamp } from '../util';

@pure
export default class Search extends Component {
  render() {
    const { result } = this.props;
    const content = Autolinker.link(result.content, { stripPrefix: false });

    const style = {
      paddingLeft: `${window.messageIndent}px`,
      textIndent: `-${window.messageIndent}px`
    };

    return (
      <p className="search-result" style={style}>
        <span className="message-time">{timestamp(new Date(result.time * 1000))}</span>
        <span>
          {' '}
          <span className="message-sender">{result.from}</span>
        </span>
        <span dangerouslySetInnerHTML={{ __html: ` ${content}` }}></span>
      </p>
    );
  }
}
