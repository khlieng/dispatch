import React, { Component } from 'react';
import pure from 'pure-render-decorator';
import { timestamp } from '../util';

@pure
export default class Search extends Component {
  componentDidUpdate(prevProps) {
    if (!prevProps.search.show && this.props.search.show) {
      this.refs.input.focus();
    }
  }

  render() {
    const { search, onSearch } = this.props;
    const results = search.results.map(result => {
      return (
        <p key={result.id}>{timestamp(new Date(result.time * 1000))} {result.from} {result.content}</p>
      );
    });

    const style = {
      display: search.show ? 'block' : 'none'
    };

    return (
      <div className="search" style={style}>
        <input
          ref="input"
          className="search-input"
          type="text"
          onChange={e => onSearch(e.target.value)}
        />
        <div className="search-results">{results}</div>
      </div>
    );
  }
}
