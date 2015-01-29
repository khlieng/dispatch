var React = require('react');
var Reflux = require('reflux');

var channelStore = require('../stores/channel');
var selectedTabStore = require('../stores/selectedTab');

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

    render: function() {
        var tab = this.state.selectedTab;
        var title;
        var topic;
        var usercount;

        if (tab.channel && this.state.channels[tab.server]) {
            var channel = this.state.channels[tab.server][tab.channel];
            if (channel) {
                title = tab.channel
                usercount = channel.users.length;
                topic = channel.topic || '';
            }
        } else {
            title = tab.server;
        }

        return (
            <div className="chat-title-bar">
                <div>
                    <span className="chat-title">{title}</span>
                    <span className="chat-topic" title={topic}>{topic}</span>
                </div>
                <span className="chat-usercount">{usercount}</span>
            </div>
        );
    }
});

module.exports = ChatTitle;