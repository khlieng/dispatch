import React, { Component } from 'react';
import pure from 'pure-render-decorator';
import SearchResult from './SearchResult';

@pure
export default class Search extends Component {
  componentDidUpdate(prevProps) {
    if (!prevProps.search.show && this.props.search.show) {
      this.refs.input.focus();
    }
  }

  handleSearch = e => this.props.onSearch(e.target.value);

  render() {
    const { search } = this.props;
    const style = {
      display: search.show ? 'block' : 'none'
    };

    const results = search.results.map(result => (
      <SearchResult key={result.id} result={result} />
    ));

    return (
      <div className="search" style={style}>
        <div className="search-input-wrap">
          <i className="icon-search" />
          <input
            ref="input"
            className="search-input"
            type="text"
            onChange={this.handleSearch}
          />
        </div>
        <div className="search-results">{results}</div>
      </div>
    );
  }
}
