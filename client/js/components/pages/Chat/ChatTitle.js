import React, { memo } from 'react';
import Navicon from 'containers/Navicon';
import Editable from 'components/ui/Editable';
import { isValidServerName } from 'state/servers';
import { isChannel, linkify } from 'utils';

const ChatTitle = ({
  status,
  title,
  tab,
  channel,
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
          <span className="chat-topic">
            {channel && linkify(channel.topic)}
          </span>
          {serverError}
        </div>
        <i className="icon-search" title="Search" onClick={onToggleSearch} />
        <i
          className="icon-cancel button-leave"
          title={closeTitle}
          onClick={onCloseClick}
        />
        <i className="icon-user button-userlist" onClick={onToggleUserList} />
      </div>
      <div className="userlist-bar">
        <i className="icon-user" />
        <span className="chat-usercount">
          {channel && channel.users.length}
        </span>
      </div>
    </div>
  );
};

export default memo(ChatTitle);
