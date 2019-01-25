import React, { memo, useState, useEffect, useCallback, useRef } from 'react';
import get from 'lodash/get';
import withModal from 'components/modals/withModal';
import Button from 'components/ui/Button';
import { join } from 'state/channels';
import { select } from 'state/tab';
import { searchChannels } from 'state/channelSearch';

const Channel = memo(({ server, name, topic, userCount, joined, ...props }) => {
  const handleJoinClick = useCallback(() => props.join([name], server), []);

  return (
    <div className="modal-channel-result">
      <div className="modal-channel-result-header">
        <h2 className="modal-channel-name" onClick={handleJoinClick}>
          {name}
        </h2>
        <span className="modal-channel-users">
          <i className="icon-user" />
          {userCount}
        </span>
        {joined ? (
          <span style={{ color: '#6bb758' }}>Joined</span>
        ) : (
          <Button
            className="modal-channel-button-join"
            category="normal"
            onClick={handleJoinClick}
          >
            Join
          </Button>
        )}
      </div>
      <p className="modal-channel-topic">{topic}</p>
    </div>
  );
});

const AddChannel = ({ search, payload: { server }, onClose, ...props }) => {
  const [q, setQ] = useState('');

  const inputEl = useRef();
  const resultsEl = useRef();
  const prevSearch = useRef('');

  useEffect(() => {
    inputEl.current.focus();
    props.searchChannels(server, '');
  }, []);

  const handleSearch = useCallback(
    e => {
      let nextQ = e.target.value.trim().toLowerCase();
      setQ(nextQ);

      if (nextQ !== q) {
        resultsEl.current.scrollTop = 0;

        while (nextQ.charAt(0) === '#') {
          nextQ = nextQ.slice(1);
        }

        if (nextQ !== prevSearch.current) {
          prevSearch.current = nextQ;
          props.searchChannels(server, nextQ);
        }
      }
    },
    [q]
  );

  const handleKey = useCallback(e => {
    if (e.key === 'Enter') {
      let channel = e.target.value.trim();

      if (channel !== '') {
        onClose(false);

        if (channel.charAt(0) !== '#') {
          channel = `#${channel}`;
        }

        props.join([channel], server);
        props.select(server, channel);
      }
    }
  }, []);

  const handleLoadMore = useCallback(
    () => props.searchChannels(server, q, search.results.length),
    [q, search.results.length]
  );

  let hasMore = !search.end;
  if (hasMore) {
    if (search.results.length < 10) {
      hasMore = false;
    } else if (
      search.results.length > 10 &&
      (search.results.length - 10) % 50 !== 0
    ) {
      hasMore = false;
    }
  }

  return (
    <>
      <div className="modal-channel-input-wrap">
        <input
          ref={inputEl}
          type="text"
          value={q}
          placeholder="Enter channel name"
          onKeyDown={handleKey}
          onChange={handleSearch}
        />
        <i
          className="icon-cancel modal-close modal-channel-close"
          onClick={onClose}
        />
      </div>
      <div ref={resultsEl} className="modal-channel-results">
        {search.results.map(channel => (
          <Channel
            key={`${server} ${channel.name}`}
            server={server}
            join={props.join}
            joined={get(
              props.channels,
              [server, channel.name, 'joined'],
              false
            )}
            {...channel}
          />
        ))}
        {hasMore && (
          <Button
            className="modal-channel-button-more"
            onClick={handleLoadMore}
          >
            Load more
          </Button>
        )}
      </div>
    </>
  );
};

export default withModal({
  name: 'channel',
  state: {
    channels: state => state.channels,
    search: state => state.channelSearch
  },
  actions: { searchChannels, join, select }
})(AddChannel);
