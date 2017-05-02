import React, { PureComponent } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { createSelector, createStructuredSelector } from 'reselect';
import { List, Map } from 'immutable';
import ChatTitle from '../components/ChatTitle';
import Search from '../components/Search';
import MessageBox from '../components/MessageBox';
import MessageInput from '../components/MessageInput';
import UserList from '../components/UserList';
import { part } from '../actions/channel';
import { openPrivateChat, closePrivateChat } from '../actions/privateChat';
import { searchMessages, toggleSearch } from '../actions/search';
import { select, setSelectedTab } from '../actions/tab';
import { runCommand, sendMessage, updateMessageHeight, fetchMessages } from '../actions/message';
import { disconnect } from '../actions/server';
import { setWrapWidth, setCharWidth } from '../actions/environment';
import { stringWidth } from '../util';
import { toggleUserList } from '../actions/ui';
import * as inputHistoryActions from '../actions/inputHistory';
import { getSelectedTab } from '../reducers/tab';
import { getSelectedMessages } from '../reducers/messages';

function updateSelected({ params, dispatch }) {
  if (params.server) {
    dispatch(setSelectedTab(params.server, params.channel || params.user));
  }
}

function updateCharWidth() {
  const charWidth = stringWidth(' ', '16px Roboto Mono');
  window.messageIndent = 6 * charWidth;
  return setCharWidth(charWidth);
}

class Chat extends PureComponent {
  componentWillMount() {
    const { dispatch } = this.props;
    dispatch(updateCharWidth());
    setTimeout(() => dispatch(updateCharWidth()), 1000);
  }

  componentDidMount() {
    updateSelected(this.props);
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.params.server !== this.props.params.server ||
      nextProps.params.channel !== this.props.params.channel ||
      nextProps.params.user !== this.props.params.user) {
      updateSelected(nextProps);
    }
  }

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
    const { title, tab, channel, search, history,
      messages, hasMoreMessages, users, showUserList, inputActions } = this.props;

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
          setWrapWidth={this.props.setWrapWidth}
          updateMessageHeight={this.props.updateMessageHeight}
          onNickClick={this.handleMessageNickClick}
          onFetchMore={this.handleFetchMore}
        />
        <MessageInput
          tab={tab}
          channel={channel}
          history={history}
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

const serverSelector = state => state.servers;
const channelSelector = state => state.channels;
const searchSelector = state => state.search;
const showUserListSelector = state => state.ui.showUserList;
const historySelector = state => {
  if (state.input.index === -1) {
    return null;
  }

  return state.input.history.get(state.input.index);
};

const selectedChannelSelector = createSelector(
  getSelectedTab,
  channelSelector,
  (tab, channels) => channels.getIn([tab.server, tab.name], Map())
);

const usersSelector = createSelector(
  selectedChannelSelector,
  channel => channel.get('users', List())
);

const titleSelector = createSelector(
  getSelectedTab,
  serverSelector,
  (tab, servers) => tab.name || servers.getIn([tab.server, 'name'])
);

const getHasMoreMessages = createSelector(
  getSelectedMessages,
  messages => messages.get(0) && typeof messages.get(0).id === 'string'
);

const mapStateToProps = createStructuredSelector({
  title: titleSelector,
  tab: getSelectedTab,
  channel: selectedChannelSelector,
  messages: getSelectedMessages,
  hasMoreMessages: getHasMoreMessages,
  users: usersSelector,
  showUserList: showUserListSelector,
  search: searchSelector,
  history: historySelector
});

function mapDispatchToProps(dispatch) {
  return {
    dispatch,
    ...bindActionCreators({
      select,
      toggleSearch,
      toggleUserList,
      searchMessages,
      runCommand,
      sendMessage,
      part,
      disconnect,
      openPrivateChat,
      closePrivateChat,
      setWrapWidth,
      updateMessageHeight
    }, dispatch),
    inputActions: bindActionCreators(inputHistoryActions, dispatch)
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(Chat);
