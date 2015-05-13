var React = require('react');
var Reflux = require('reflux');
var Router = require('react-router');

var ChatTitle = require('./ChatTitle.jsx');
var Search = require('./Search.jsx');
var MessageBox = require('./MessageBox.jsx');
var MessageInput = require('./MessageInput.jsx');
var UserList = require('./UserList.jsx');
var selectedTabStore = require('../stores/selectedTab');
var tabActions = require('../actions/tab');

var Chat = React.createClass({
    mixins: [
        Router.State,
        Reflux.connect(selectedTabStore, 'selectedTab')
    ],

    getInitialState: function() {
        return {
            selectedTab: selectedTabStore.getState()
        };
    },

    componentWillMount: function() {
        if (!window.loaded) {
            var p = this.getParams();

            if (p.channel) {
                tabActions.select(p.server, '#' + p.channel);
            } else if (p.server) {
                tabActions.select(p.server);
            }
        }
    },

    render: function() {
        var chatClass;
        var tab = this.state.selectedTab;

        if (!tab.channel) {
            chatClass = 'chat-server';
        } else if (tab.channel[0] !== '#') {
            chatClass = 'chat-private';
        } else {
            chatClass = 'chat-channel';
        }

        return (
            <div className={chatClass}>
                <ChatTitle />
                <Search />
                <MessageBox indent={window.messageIndent} />
                <MessageInput />
                <UserList />
            </div>
        );
    }
});

module.exports = Chat;