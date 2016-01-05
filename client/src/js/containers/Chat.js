import React, { Component } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { createSelector } from 'reselect';
import { List, Map } from 'immutable';
import pure from 'pure-render-decorator';
import ChatTitle from '../components/ChatTitle';
import Search from '../components/Search';
import MessageBox from '../components/MessageBox';
import MessageInput from '../components/MessageInput';
import UserList from '../components/UserList';
import { part } from '../actions/channel';
import { openPrivateChat, closePrivateChat } from '../actions/privateChat';
import { searchMessages, toggleSearch } from '../actions/search';
import { select, setSelectedChannel, setSelectedUser } from '../actions/tab';
import { runCommand, sendMessage } from '../actions/message';
import { disconnect } from '../actions/server';
import * as inputHistoryActions from '../actions/inputHistory';
import { setWrapWidth, setCharWidth } from '../actions/environment';
import { stringWidth, wrapMessages } from '../util';

function updateSelected({ params, dispatch }) {
  if (params.channel) {
    dispatch(setSelectedChannel(params.server, params.channel));
  } else if (params.user) {
    dispatch(setSelectedUser(params.server, params.user));
  } else if (params.server) {
    dispatch(setSelectedChannel(params.server));
  }
}

function updateCharWidth() {
  const charWidth = stringWidth(' ', '16px Droid Sans Mono');
  window.messageIndent = 6 * charWidth;
  return setCharWidth(charWidth);
}

@pure
class Chat extends Component {
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
    if (tab.channel) {
      dispatch(searchMessages(tab.server, tab.channel, phrase));
    }
  }

  render() {
    const { tab, channel, search, history, dispatch } = this.props;

    let chatClass;
    if (tab.channel) {
      chatClass = 'chat-channel';
    } else if (tab.user) {
      chatClass = 'chat-private';
    } else {
      chatClass = 'chat-server';
    }

    return (
      <div className={chatClass}>
        <ChatTitle {...this.props } />
        <Search
          search={search}
          onSearch={this.handleSearch}
        />
        <MessageBox {...this.props } />
        <MessageInput
          tab={tab}
          channel={channel}
          runCommand={this.props.runCommand}
          sendMessage={this.props.sendMessage}
          history={history}
          {...bindActionCreators(inputHistoryActions, dispatch)}
        />
        <UserList {...this.props} />
      </div>
    );
  }
}

const tabSelector = state => state.tab.selected;
const messageSelector = state => state.messages;

const selectedMessagesSelector = createSelector(
  tabSelector,
  messageSelector,
  (tab, messages) => messages.getIn([tab.server, tab.channel || tab.user || tab.server], List())
);

const wrapWidthSelector = state => state.environment.get('wrapWidth');
const charWidthSelector = state => state.environment.get('charWidth');

const wrappedMessagesSelector = createSelector(
  selectedMessagesSelector,
  wrapWidthSelector,
  charWidthSelector,
  (messages, width, charWidth) => wrapMessages(messages, width, charWidth, 6 * charWidth)
);

function mapStateToProps(state) {
  const tab = state.tab.selected;
  const channel = state.channels.getIn([tab.server, tab.channel], Map());

  let title;
  if (tab.channel) {
    title = tab.channel;
  } else if (tab.user) {
    title = tab.user;
  } else {
    title = state.servers.getIn([tab.server, 'name']);
  }

  return {
    title,
    search: state.search,
    users: channel.get('users', List()),
    history: state.input.index === -1 ? null : state.input.history.get(state.input.index),
    messages: wrappedMessagesSelector(state),
    channel,
    tab
  };
}

function mapDispatchToProps(dispatch) {
  return {
    dispatch,
    ...bindActionCreators({
      select,
      toggleSearch,
      searchMessages,
      runCommand,
      sendMessage,
      part,
      disconnect,
      openPrivateChat,
      closePrivateChat,
      setWrapWidth
    }, dispatch)
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(Chat);
