import React from 'react';
import Reflux from 'reflux';
import Router from 'react-router';
import ChatTitle from './ChatTitle.jsx';
import Search from './Search.jsx';
import MessageBox from './MessageBox.jsx';
import MessageInput from './MessageInput.jsx';
import UserList from './UserList.jsx';
import selectedTabStore from '../stores/selectedTab';
import tabActions from '../actions/tab';
import PureMixin from '../mixins/pure';

export default React.createClass({
    mixins: [
        PureMixin,
        Router.State,
        Reflux.connect(selectedTabStore, 'selectedTab')
    ],

    componentWillMount() {
        if (!window.loaded) {
            const { params } = this.props;
            if (params.channel) {
                tabActions.select(params.server, '#' + params.channel);
            } else if (params.server) {
                tabActions.select(params.server);
            }
        }
    },

    render() {
        let chatClass;
        const tab = this.state.selectedTab;

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
                <MessageBox />
                <MessageInput />
                <UserList />
            </div>
        );
    }
});
