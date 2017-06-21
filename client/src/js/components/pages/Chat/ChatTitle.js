import React, { PureComponent } from 'react';
import { List } from 'immutable';
import Navicon from 'containers/Navicon';
import Editable from 'components/ui/Editable';
import { isValidServerName } from 'state/servers';
import { linkify } from 'util';

export default class ChatTitle extends PureComponent {
  render() {
    const { title, tab, channel, onTitleChange,
      onToggleSearch, onToggleUserList, onCloseClick } = this.props;

    let closeTitle;
    if (tab.isChannel()) {
      closeTitle = 'Leave';
    } else if (tab.name) {
      closeTitle = 'Close';
    } else {
      closeTitle = 'Disconnect';
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
            <span className="chat-topic">{linkify(channel.get('topic')) || null}</span>
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
          <span className="chat-usercount">{channel.get('users', List()).size}</span>
        </div>
      </div>
    );
  }
}
