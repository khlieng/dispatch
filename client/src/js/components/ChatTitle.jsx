var React = require('react');
var Reflux = require('reflux');

var channelStore = require('../stores/channel');
var selectedTabStore = require('../stores/selectedTab');
var serverActions = require('../actions/server');
var channelActions = require('../actions/channel');
var searchActions = require('../actions/search');
var privateChatActions = require('../actions/privateChat');

var ChatTitle = React.createClass({
    mixins: [
        Reflux.listenTo(channelStore, 'channelsChanged'),
        Reflux.listenTo(selectedTabStore, 'selectedTabChanged')
    ],

    getInitialState() {
        var tab = selectedTabStore.getState();

        return {
            usercount: channelStore.getUsers(tab.server, tab.channel).size,
            selectedTab: tab
        };
    },

    channelsChanged() {
        var tab = this.state.selectedTab;
        
        this.setState({ usercount: channelStore.getUsers(tab.server, tab.channel).size });
    },

    selectedTabChanged(tab) {
        this.setState({
            selectedTab: tab,
            usercount: channelStore.getUsers(tab.server, tab.channel).size
        });
    },

    handleLeaveClick() {
        var tab = this.state.selectedTab;

        if (!tab.channel) {
            serverActions.disconnect(tab.server);
        } else if (tab.channel[0] === '#') {
            channelActions.part([tab.channel], tab.server);
        } else {
            privateChatActions.close(tab.server, tab.channel);
        }
    },

    render() {
        var tab = this.state.selectedTab;
        var leaveTitle;

        if (!tab.channel) {
            leaveTitle = 'Disconnect';
        } else if (tab.channel[0] !== '#') {
            leaveTitle = 'Close';
        } else {
            leaveTitle = 'Leave';
        }

        return (
            <div>
                <div className="chat-title-bar">
                    <span className="chat-title">{tab.name}</span>
                    <i className="icon-search" title="Search" onClick={searchActions.toggle}></i>
                    <i 
                        className="icon-logout button-leave" 
                        title={leaveTitle} 
                        onClick={this.handleLeaveClick}></i>
                </div>
                <div className="userlist-bar">
                    <i className="icon-user"></i>
                    <span className="chat-usercount">{this.state.usercount || null}</span>
                </div>
            </div>
        );
    }
});

module.exports = ChatTitle;