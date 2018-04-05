import React, { PureComponent } from 'react';
import { timestamp, linkify } from 'utils';

export default class SearchResult extends PureComponent {
  render() {
    const { result } = this.props;
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
        <span>{' '}{linkify(result.content)}</span>
      </p>
    );
  }
}
