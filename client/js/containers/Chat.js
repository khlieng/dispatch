import { bindActionCreators } from 'redux';
import { createStructuredSelector } from 'reselect';
import Chat from 'components/pages/Chat';
import { getSelectedTabTitle } from 'state';
import {
  getSelectedChannel,
  getSelectedChannelUsers,
  part
} from 'state/channels';
import {
  getCurrentInputHistoryEntry,
  addInputHistory,
  resetInputHistory,
  incrementInputHistory,
  decrementInputHistory
} from 'state/input';
import {
  getSelectedMessages,
  getHasMoreMessages,
  runCommand,
  sendMessage,
  fetchMessages,
  addFetchedMessages
} from 'state/messages';
import { openModal } from 'state/modals';
import { openPrivateChat, closePrivateChat } from 'state/privateChats';
import { getSearch, searchMessages, toggleSearch } from 'state/search';
import {
  getCurrentNick,
  getCurrentNetworkError,
  disconnect,
  setNick,
  setNetworkName
} from 'state/networks';
import { getSettings } from 'state/settings';
import { getSelectedTab, select } from 'state/tab';
import { getShowUserList, toggleUserList } from 'state/ui';
import connect from 'utils/connect';

const mapState = createStructuredSelector({
  channel: getSelectedChannel,
  currentInputHistoryEntry: getCurrentInputHistoryEntry,
  hasMoreMessages: getHasMoreMessages,
  messages: getSelectedMessages,
  nick: getCurrentNick,
  search: getSearch,
  showUserList: getShowUserList,
  error: getCurrentNetworkError,
  tab: getSelectedTab,
  title: getSelectedTabTitle,
  users: getSelectedChannelUsers,
  coloredNicks: state => getSettings(state).coloredNicks
});

const mapDispatch = dispatch => ({
  ...bindActionCreators(
    {
      addFetchedMessages,
      closePrivateChat,
      disconnect,
      fetchMessages,
      openModal,
      openPrivateChat,
      part,
      runCommand,
      searchMessages,
      select,
      sendMessage,
      setNick,
      setNetworkName,
      toggleSearch,
      toggleUserList
    },
    dispatch
  ),

  inputActions: bindActionCreators(
    {
      add: addInputHistory,
      reset: resetInputHistory,
      increment: incrementInputHistory,
      decrement: decrementInputHistory
    },
    dispatch
  )
});

export default connect(mapState, mapDispatch)(Chat);
