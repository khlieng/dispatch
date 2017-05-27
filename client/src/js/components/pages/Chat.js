import React, { Component } from 'react';
import ChatTitle from '../ChatTitle';
import Search from '../Search';
import MessageBox from '../MessageBox';
import MessageInput from '../MessageInput';
import UserList from '../UserList';

export default class Chat extends Component {
  handleCloseClick = () => {
    const { tab, part, closePrivateChat, disconnect } = this.props;

    if (tab.isChannel()) {
      part([tab.name], tab.server);
    } else if (tab.name) {
      closePrivateChat(tab.server, tab.name);
    } else {
      disconnect(tab.server);
    }
  };

  handleSearch = phrase => {
    const { tab, searchMessages } = this.props;
    if (tab.isChannel()) {
      searchMessages(tab.server, tab.name, phrase);
    }
  };

  handleNickClick = nick => {
    const { tab, openPrivateChat, select } = this.props;
    openPrivateChat(tab.server, nick);
    select(tab.server, nick);
  };

  render() {
    const {
      channel,
      currentInputHistoryEntry,
      hasMoreMessages,
      messages,
      nick,
      search,
      showUserList,
      tab,
      title,
      users,

      fetchMessages,
      inputActions,
      runCommand,
      sendMessage,
      toggleSearch,
      toggleUserList
    } = this.props;

    let chatClass;
    if (tab.isChannel()) {
      chatClass = 'chat-channel';
    } else if (tab.name) {
      chatClass = 'chat-private';
    } else {
      chatClass = 'chat-server';
    }

    return (
      <div className={chatClass}>
        <ChatTitle
          channel={channel}
          tab={tab}
          title={title}
          onCloseClick={this.handleCloseClick}
          onToggleSearch={toggleSearch}
          onToggleUserList={toggleUserList}
        />
        <Search
          search={search}
          onSearch={this.handleSearch}
        />
        <MessageBox
          hasMoreMessages={hasMoreMessages}
          messages={messages}
          tab={tab}
          onFetchMore={fetchMessages}
          onNickClick={this.handleNickClick}
        />
        <MessageInput
          currentHistoryEntry={currentInputHistoryEntry}
          nick={nick}
          tab={tab}
          onCommand={runCommand}
          onMessage={sendMessage}
          {...inputActions}
        />
        <UserList
          showUserList={showUserList}
          users={users}
          onNickClick={this.handleNickClick}
        />
      </div>
    );
  }
}
