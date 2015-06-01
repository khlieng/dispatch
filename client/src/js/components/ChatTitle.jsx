var React = require('react');
var Reflux = require('reflux');
var Autolinker = require('autolinker');

var channelStore = require('../stores/channel');
var selectedTabStore = require('../stores/selectedTab');
var serverActions = require('../actions/server');
var channelActions = require('../actions/channel');
var searchActions = require('../actions/search');
var privateChatActions = require('../actions/privateChat');
var PureMixin = require('../mixins/pure');

function buildState(tab) {
    return {
        selectedTab: tab,
        usercount: channelStore.getUsers(tab.server, tab.channel).size,
        topic: channelStore.getTopic(tab.server, tab.channel)
    };
}

var ChatTitle = React.createClass({
    mixins: [
        PureMixin,
        Reflux.listenTo(channelStore, 'channelsChanged'),
        Reflux.listenTo(selectedTabStore, 'selectedTabChanged')
    ],

    getInitialState() {
        return buildState(selectedTabStore.getState());
    },

    channelsChanged() {
        this.setState(buildState(this.state.selectedTab));
    },

    selectedTabChanged(tab) {
        this.setState(buildState(tab));
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
        var topic = Autolinker.link(this.state.topic || '', { keepOriginalText: true });
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
                    <div className="chat-topic-wrap">
                        <span className="chat-topic" dangerouslySetInnerHTML={{ __html: topic }}></span>
                    </div>
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