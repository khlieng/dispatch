var React = require('react');
var Reflux = require('reflux');

var channelStore = require('../stores/channel');
var selectedTabStore = require('../stores/selectedTab');
var channelActions = require('../actions/channel');
var searchActions = require('../actions/search');
var privateChatActions = require('../actions/privateChat');

var ChatTitle = React.createClass({
    mixins: [
        Reflux.connect(channelStore, 'channels'),
        Reflux.connect(selectedTabStore, 'selectedTab')
    ],

    getInitialState: function() {
        return {
            channels: channelStore.getState(),
            selectedTab: selectedTabStore.getState()
        };
    },

    handleLeaveClick: function() {
        var tab = this.state.selectedTab;

        if (tab.channel[0] === '#') {
            channelActions.part([tab.channel], tab.server);
        } else if (tab.channel) {
            privateChatActions.close(tab.server, tab.channel);
        }
    },

    render: function() {
        var tab = this.state.selectedTab;
        var usercount = channelStore.getUsers(tab.server, tab.channel).length;
        var iconStyle = {};
        var userListStyle = {};

        if (!tab.channel) {
            iconStyle.display = 'none';
            userListStyle.display = 'none';
        } else if (tab.channel[0] !== '#') {
            userListStyle.display = 'none';
        }

        return (
            <div>
                <div className="chat-title-bar">
                    <span className="chat-title">{tab.name}</span>
                    <i className="icon-search" title="Search" style={iconStyle} onClick={searchActions.toggle}></i>
                    <i className="icon-logout button-leave" title="Leave" style={iconStyle} onClick={this.handleLeaveClick}></i>
                </div>
                <div className="userlist-bar">
                    <i className="icon-user" style={userListStyle}></i>
                    <span className="chat-usercount" style={userListStyle}>{usercount || null}</span>
                </div>
            </div>
        );
    }
});

module.exports = ChatTitle;