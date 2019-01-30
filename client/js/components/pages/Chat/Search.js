import React, { memo, useRef, useEffect } from 'react';
import SearchResult from './SearchResult';

const Search = ({ search, onSearch }) => {
  const inputEl = useRef();

  useEffect(() => {
    if (search.show) {
      inputEl.current.focus();
    }
  }, [search.show]);

  const style = {
    display: search.show ? 'block' : 'none'
  };

  let i = 0;
  const results = search.results.map(result => (
    <SearchResult key={i++} result={result} />
  ));

  return (
    <div className="search" style={style}>
      <div className="search-input-wrap">
        <i className="icon-search" />
        <input
          ref={inputEl}
          className="search-input"
          type="text"
          onChange={e => onSearch(e.target.value)}
        />
      </div>
      <div className="search-results">{results}</div>
    </div>
  );
};

export default memo(Search);
