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
        var usercount = channelStore.getUsers(tab.server, tab.channel).length;

        return (
            <div className="chat-title-bar">
                <div>
                    <span className="chat-title">{tab.name}</span>
                </div>
                <span className="chat-usercount">{usercount || null}</span>
            </div>
        );
    }
});

module.exports = ChatTitle;