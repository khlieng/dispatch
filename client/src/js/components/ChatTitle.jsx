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

        if (tab.channel) {
            var channel = this.state.channels[tab.server][tab.channel];
            if (channel) {
                title = tab.channel
                title += ' [';
                title += channel.users.length;
                title += ']';

                if (channel.topic) {
                    title += ': ' + channel.topic;
                }
            }
        } else {
            title = tab.server;
        }

        return (
            <div className="chat-title-bar">
                <span className="chat-title" title={title}>{title}</span>
            </div>
        );
    }
});

module.exports = ChatTitle;