import React, { PureComponent } from 'react';
import { List } from 'immutable';
import Navicon from '../components/Navicon';
import { linkify } from '../util';

export default class ChatTitle extends PureComponent {
  handleLeaveClick = () => {
    const { tab, disconnect, part, closePrivateChat } = this.props;

    if (tab.channel) {
      part([tab.channel], tab.server);
    } else if (tab.user) {
      closePrivateChat(tab.server, tab.user);
    } else {
      disconnect(tab.server);
    }
  };

  render() {
    const { title, tab, channel, toggleSearch, toggleUserList } = this.props;
    let topic = channel.get('topic');
    topic = topic ? linkify(topic) : null;

    let leaveTitle;
    if (tab.channel) {
      leaveTitle = 'Leave';
    } else if (tab.user) {
      leaveTitle = 'Close';
    } else {
      leaveTitle = 'Disconnect';
    }

    return (
      <div>
        <div className="chat-title-bar">
          <Navicon />
          <span className="chat-title">{title}</span>
          <div className="chat-topic-wrap">
            <span className="chat-topic">{topic}</span>
          </div>
          <i className="icon-search" title="Search" onClick={toggleSearch} />
          <i
            className="icon-logout button-leave"
            title={leaveTitle}
            onClick={this.handleLeaveClick}
          />
          <i className="icon-user button-userlist" onClick={toggleUserList} />
        </div>
        <div className="userlist-bar">
          <i className="icon-user" />
          <span className="chat-usercount">{channel.get('users', List()).size || null}</span>
        </div>
      </div>
    );
  }
}
