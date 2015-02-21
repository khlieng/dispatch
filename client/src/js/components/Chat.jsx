var React = require('react');

var ChatTitle = require('./ChatTitle.jsx');
var MessageBox = require('./MessageBox.jsx');
var MessageInput = require('./MessageInput.jsx');
var UserList = require('./UserList.jsx');

var Chat = React.createClass({
    render: function() {
        return (
            <div>
                <ChatTitle />
                <MessageBox indent={window.messageIndent} />
                <MessageInput />
                <UserList />
            </div>
        );
    }
});

module.exports = Chat;