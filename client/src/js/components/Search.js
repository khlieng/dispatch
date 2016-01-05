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

  handleSearch = e => this.props.onSearch(e.target.value)

  render() {
    const { search } = this.props;
    const style = {
      display: search.show ? 'block' : 'none'
    };

    const results = search.results.map(result => (
      <p key={result.id}>
        {timestamp(new Date(result.time * 1000))} {result.from} {result.content}
      </p>
    ));

    return (
      <div className="search" style={style}>
        <input
          ref="input"
          className="search-input"
          type="text"
          onChange={this.handleSearch}
        />
        <div className="search-results">{results}</div>
      </div>
    );
  }
}
