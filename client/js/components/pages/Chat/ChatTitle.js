import React, { memo } from 'react';
import { FiUsers, FiSearch, FiX } from 'react-icons/fi';
import Navicon from 'components/ui/Navicon';
import Button from 'components/ui/Button';
import Editable from 'components/ui/Editable';
import { isValidNetworkName } from 'state/networks';
import { isChannel } from 'utils';

const ChatTitle = ({
  error,
  title,
  tab,
  channel,
  openModal,
  onTitleChange,
  onToggleSearch,
  onToggleUserList,
  onCloseClick
}) => {
  let closeTitle;
  if (isChannel(tab)) {
    closeTitle = 'Leave';
  } else if (tab.name) {
    closeTitle = 'Close';
  } else {
    closeTitle = 'Disconnect';
  }

  let networkError = null;
  if (!tab.name && error) {
    networkError = <span className="chat-topic error">Error: {error}</span>;
  }

  return (
    <div>
      <div className="chat-title-bar">
        <Navicon />
        <Editable
          className="chat-title"
          editable={!tab.name}
          value={title}
          validate={isValidNetworkName}
          onChange={onTitleChange}
        >
          <span className="chat-title">{title}</span>
        </Editable>
        <div className="chat-topic-wrap">
          {channel?.topic && (
            <span
              className="chat-topic"
              onClick={() => openModal('topic', channel.name)}
            >
              {channel.topic}
            </span>
          )}
          {networkError}
        </div>
        {tab.name && (
          <Button
            icon={FiSearch}
            title="Search"
            aria-label="Search"
            onClick={onToggleSearch}
          />
        )}
        <Button
          icon={FiX}
          title={closeTitle}
          aria-label={closeTitle}
          onClick={onCloseClick}
        />
        <Button
          icon={FiUsers}
          className="button-userlist"
          aria-label="Users"
          onClick={onToggleUserList}
        />
      </div>
      <div className="userlist-bar">
        <FiUsers />
        {channel?.users.length}
      </div>
    </div>
  );
};

export default memo(ChatTitle);
