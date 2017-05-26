import React, { PureComponent } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';
import ChatTitle from '../components/ChatTitle';
import Search from '../components/Search';
import MessageBox from '../components/MessageBox';
import MessageInput from '../components/MessageInput';
import UserList from '../components/UserList';
import { getSelectedTabTitle } from '../state';
import { getSelectedChannel, getSelectedChannelUsers, part } from '../state/channels';
import { getCurrentInputHistoryEntry, addInputHistory, resetInputHistory,
  incrementInputHistory, decrementInputHistory } from '../state/input';
import { getSelectedMessages, getHasMoreMessages,
  runCommand, sendMessage, fetchMessages } from '../state/messages';
import { openPrivateChat, closePrivateChat } from '../state/privateChats';
import { getSearch, searchMessages, toggleSearch } from '../state/search';
import { getCurrentNick, disconnect } from '../state/servers';
import { getSelectedTab, select } from '../state/tab';
import { getShowUserList, toggleUserList } from '../state/ui';

class Chat extends PureComponent {
  handleSearch = phrase => {
    const { dispatch, tab } = this.props;
    if (tab.isChannel()) {
      dispatch(searchMessages(tab.server, tab.name, phrase));
    }
  };

  handleMessageNickClick = message => {
    const { tab } = this.props;

    this.props.openPrivateChat(tab.server, message.from);
    this.props.select(tab.server, message.from);
  };

  handleFetchMore = () => this.props.dispatch(fetchMessages());

  render() {
    const { title, tab, channel, search, currentInputHistoryEntry,
      messages, hasMoreMessages, users, showUserList, nick, inputActions } = this.props;

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
          title={title}
          tab={tab}
          channel={channel}
          toggleSearch={this.props.toggleSearch}
          toggleUserList={this.props.toggleUserList}
          disconnect={this.props.disconnect}
          part={this.props.part}
          closePrivateChat={this.props.closePrivateChat}
        />
        <Search
          search={search}
          onSearch={this.handleSearch}
        />
        <MessageBox
          messages={messages}
          hasMoreMessages={hasMoreMessages}
          tab={tab}
          onNickClick={this.handleMessageNickClick}
          onFetchMore={this.handleFetchMore}
        />
        <MessageInput
          tab={tab}
          channel={channel}
          currentHistoryEntry={currentInputHistoryEntry}
          nick={nick}
          runCommand={this.props.runCommand}
          sendMessage={this.props.sendMessage}
          {...inputActions}
        />
        <UserList
          users={users}
          tab={tab}
          showUserList={showUserList}
          select={this.props.select}
          openPrivateChat={this.props.openPrivateChat}
        />
      </div>
    );
  }
}

const mapState = createStructuredSelector({
  channel: getSelectedChannel,
  currentInputHistoryEntry: getCurrentInputHistoryEntry,
  hasMoreMessages: getHasMoreMessages,
  messages: getSelectedMessages,
  nick: getCurrentNick,
  search: getSearch,
  showUserList: getShowUserList,
  tab: getSelectedTab,
  title: getSelectedTabTitle,
  users: getSelectedChannelUsers
});

function mapDispatch(dispatch) {
  return {
    dispatch,
    ...bindActionCreators({
      closePrivateChat,
      disconnect,
      openPrivateChat,
      part,
      runCommand,
      searchMessages,
      select,
      sendMessage,
      toggleSearch,
      toggleUserList
    }, dispatch),
    inputActions: bindActionCreators({
      add: addInputHistory,
      reset: resetInputHistory,
      increment: incrementInputHistory,
      decrement: decrementInputHistory
    }, dispatch)
  };
}

export default connect(mapState, mapDispatch)(Chat);
