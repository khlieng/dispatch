import React, { Component } from 'react';
import { isChannel } from 'utils';
import ChatTitle from './ChatTitle';
import Search from './Search';
import MessageBox from './MessageBox';
import MessageInput from './MessageInput';
import UserList from './UserList';

export default class Chat extends Component {
  handleCloseClick = () => {
    const { tab, part, closePrivateChat, disconnect } = this.props;

    if (isChannel(tab)) {
      part([tab.name], tab.server);
    } else if (tab.name) {
      closePrivateChat(tab.server, tab.name);
    } else {
      disconnect(tab.server);
    }
  };

  handleSearch = phrase => {
    const { tab, searchMessages } = this.props;
    if (isChannel(tab)) {
      searchMessages(tab.server, tab.name, phrase);
    }
  };

  handleNickClick = nick => {
    const { tab, openPrivateChat, select } = this.props;
    openPrivateChat(tab.server, nick);
    select(tab.server, nick);
  };

  handleTitleChange = title => {
    const { setServerName, tab } = this.props;
    setServerName(title, tab.server);
  };

  handleNickChange = nick => {
    const { setNick, tab } = this.props;
    setNick(nick, tab.server, true);
  };

  handleNickEditDone = nick => {
    const { setNick, tab } = this.props;
    setNick(nick, tab.server);
  };

  render() {
    const {
      channel,
      coloredNicks,
      currentInputHistoryEntry,
      hasMoreMessages,
      messages,
      nick,
      search,
      showUserList,
      status,
      tab,
      title,
      users,

      addFetchedMessages,
      fetchMessages,
      inputActions,
      runCommand,
      sendMessage,
      toggleSearch,
      toggleUserList
    } = this.props;
    let chatClass;
    if (isChannel(tab)) {
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
          status={status}
          tab={tab}
          title={title}
          onCloseClick={this.handleCloseClick}
          onTitleChange={this.handleTitleChange}
          onToggleSearch={toggleSearch}
          onToggleUserList={toggleUserList}
        />
        <Search search={search} onSearch={this.handleSearch} />
        <MessageBox
          coloredNicks={coloredNicks}
          hasMoreMessages={hasMoreMessages}
          messages={messages}
          tab={tab}
          hideTopDate={search.show}
          onAddMore={addFetchedMessages}
          onFetchMore={fetchMessages}
          onNickClick={this.handleNickClick}
        />
        <MessageInput
          currentHistoryEntry={currentInputHistoryEntry}
          nick={nick}
          tab={tab}
          onCommand={runCommand}
          onMessage={sendMessage}
          onNickChange={this.handleNickChange}
          onNickEditDone={this.handleNickEditDone}
          {...inputActions}
        />
        <UserList
          coloredNicks={coloredNicks}
          showUserList={showUserList}
          users={users}
          onNickClick={this.handleNickClick}
        />
      </div>
    );
  }
}
