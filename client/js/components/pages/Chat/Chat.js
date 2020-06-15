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
      part([tab.name], tab.network);
    } else if (tab.name) {
      closePrivateChat(tab.network, tab.name);
    } else {
      disconnect(tab.network);
    }
  };

  handleSearch = phrase => {
    const { tab, searchMessages } = this.props;
    if (isChannel(tab)) {
      searchMessages(tab.network, tab.name, phrase);
    }
  };

  handleNickClick = nick => {
    const { tab, openPrivateChat, select } = this.props;
    openPrivateChat(tab.network, nick);
    select(tab.network, nick);
  };

  handleTitleChange = title => {
    const { setNetworkName, tab } = this.props;
    setNetworkName(title, tab.network);
  };

  handleNickChange = nick => {
    const { setNick, tab } = this.props;
    setNick(nick, tab.network, true);
  };

  handleNickEditDone = nick => {
    const { setNick, tab } = this.props;
    setNick(nick, tab.network);
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
      error,
      tab,
      title,
      users,

      addFetchedMessages,
      fetchMessages,
      inputActions,
      openModal,
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
      chatClass = 'chat-network';
    }

    return (
      <div className={chatClass}>
        <ChatTitle
          channel={channel}
          error={error}
          tab={tab}
          title={title}
          openModal={openModal}
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
