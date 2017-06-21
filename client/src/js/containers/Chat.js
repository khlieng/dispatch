import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';
import Chat from '../components/pages/Chat';
import { getSelectedTabTitle } from '../state';
import { getSelectedChannel, getSelectedChannelUsers, part } from '../state/channels';
import { getCurrentInputHistoryEntry, addInputHistory, resetInputHistory,
  incrementInputHistory, decrementInputHistory } from '../state/input';
import { getSelectedMessages, getHasMoreMessages,
  runCommand, sendMessage, fetchMessages, addFetchedMessages } from '../state/messages';
import { openPrivateChat, closePrivateChat } from '../state/privateChats';
import { getSearch, searchMessages, toggleSearch } from '../state/search';
import { getCurrentNick, disconnect, setNick, setServerName } from '../state/servers';
import { getSelectedTab, select } from '../state/tab';
import { getShowUserList, toggleUserList } from '../state/ui';

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

const mapDispatch = dispatch => ({
  ...bindActionCreators({
    addFetchedMessages,
    closePrivateChat,
    disconnect,
    fetchMessages,
    openPrivateChat,
    part,
    runCommand,
    searchMessages,
    select,
    sendMessage,
    setNick,
    setServerName,
    toggleSearch,
    toggleUserList
  }, dispatch),

  inputActions: bindActionCreators({
    add: addInputHistory,
    reset: resetInputHistory,
    increment: incrementInputHistory,
    decrement: decrementInputHistory
  }, dispatch)
});

export default connect(mapState, mapDispatch)(Chat);
