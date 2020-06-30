import React, { memo, useState, useEffect, useRef } from 'react';
import Modal from 'react-modal';
import { useSelector, useDispatch } from 'react-redux';
import { FiUsers, FiX } from 'react-icons/fi';
import Text from 'components/Text';
import useModal from 'components/modals/useModal';
import Button from 'components/ui/Button';
import { join } from 'state/channels';
import { searchChannels } from 'state/channelSearch';
import { linkify } from 'utils';

const Channel = memo(({ network, name, topic, userCount, joined }) => {
  const dispatch = useDispatch();

  const handleClick = () => dispatch(join([name], network));

  return (
    <div className="modal-channel-result">
      <div className="modal-channel-result-header">
        <h2 className="modal-channel-name" onClick={handleClick}>
          {name}
        </h2>
        <FiUsers />
        <span className="modal-channel-users">{userCount}</span>
        {joined ? (
          <span style={{ color: '#6bb758' }}>Joined</span>
        ) : (
          <Button
            className="modal-channel-button-join"
            category="normal"
            onClick={handleClick}
          >
            Join
          </Button>
        )}
      </div>
      <p className="modal-channel-topic">
        <Text>{linkify(topic)}</Text>
      </p>
    </div>
  );
});

const AddChannel = () => {
  const [modal, network, closeModal] = useModal('channel');

  const channels = useSelector(state => state.channels);
  const search = useSelector(state => state.channelSearch);
  const dispatch = useDispatch();
  const [q, setQ] = useState('');

  const inputEl = useRef();
  const resultsEl = useRef();
  const prevSearch = useRef('');

  useEffect(() => {
    if (modal.isOpen) {
      dispatch(searchChannels(network, ''));
      setTimeout(() => inputEl.current.focus(), 0);
    } else {
      prevSearch.current = '';
      setQ('');
    }
  }, [modal.isOpen]);

  const handleSearch = e => {
    let nextQ = e.target.value.trim().toLowerCase();
    setQ(nextQ);

    if (nextQ !== q) {
      resultsEl.current.scrollTop = 0;

      while (nextQ.charAt(0) === '#') {
        nextQ = nextQ.slice(1);
      }

      if (nextQ !== prevSearch.current) {
        prevSearch.current = nextQ;
        dispatch(searchChannels(network, nextQ));
      }
    }
  };

  const handleKey = e => {
    if (e.key === 'Enter') {
      let channel = e.target.value.trim();

      if (channel !== '') {
        closeModal(false);

        if (channel.charAt(0) !== '#') {
          channel = `#${channel}`;
        }

        dispatch(join([channel], network));
      }
    }
  };

  const handleLoadMore = () =>
    dispatch(searchChannels(network, q, search.results.length));

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
    <Modal {...modal}>
      <div className="modal-channel-input-wrap">
        <input
          ref={inputEl}
          type="text"
          value={q}
          placeholder="Enter channel name"
          onKeyDown={handleKey}
          onChange={handleSearch}
        />
        <Button
          icon={FiX}
          className="modal-close modal-channel-close"
          onClick={closeModal}
        />
      </div>
      <div ref={resultsEl} className="modal-channel-results">
        {search.results.map(channel => (
          <Channel
            key={`${network} ${channel.name}`}
            network={network}
            joined={channels[network]?.[channel.name]?.joined}
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
    </Modal>
  );
};

export default AddChannel;
