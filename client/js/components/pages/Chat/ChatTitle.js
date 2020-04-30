import React, { memo } from 'react';
import { FiUsers, FiSearch, FiX } from 'react-icons/fi';
import Navicon from 'components/ui/Navicon';
import Button from 'components/ui/Button';
import Editable from 'components/ui/Editable';
import { isValidServerName } from 'state/servers';
import { isChannel } from 'utils';

const ChatTitle = ({
  status,
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

  let serverError = null;
  if (!tab.name && status.error) {
    serverError = (
      <span className="chat-topic error">Error: {status.error}</span>
    );
  }

  return (
    <div>
      <div className="chat-title-bar">
        <Navicon />
        <Editable
          className="chat-title"
          editable={!tab.name}
          value={title}
          validate={isValidServerName}
          onChange={onTitleChange}
        >
          <span className="chat-title">{title}</span>
        </Editable>
        <div className="chat-topic-wrap">
          {channel && channel.topic && (
            <span
              className="chat-topic"
              onClick={() =>
                openModal('topic', {
                  topic: channel.topic,
                  channel: channel.name
                })
              }
            >
              {channel.topic}
            </span>
          )}
          {serverError}
        </div>
        {tab.name && (
          <Button icon={FiSearch} title="Search" onClick={onToggleSearch} />
        )}
        <Button icon={FiX} title={closeTitle} onClick={onCloseClick} />
        <Button
          icon={FiUsers}
          className="button-userlist"
          onClick={onToggleUserList}
        />
      </div>
      <div className="userlist-bar">
        <FiUsers />
        {channel && channel.users.length}
      </div>
    </div>
  );
};

export default memo(ChatTitle);
