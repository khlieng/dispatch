var React = require('react');
var Reflux = require('reflux');
var Router = require('react-router');

var ChatTitle = require('./ChatTitle.jsx');
var Search = require('./Search.jsx');
var MessageBox = require('./MessageBox.jsx');
var MessageInput = require('./MessageInput.jsx');
var UserList = require('./UserList.jsx');
var tabActions = require('../actions/tab');

var Chat = React.createClass({
    mixins: [Router.State],

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
        return (
            <div>
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