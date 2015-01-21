var React = require('react');
var Reflux = require('reflux');
var _ = require('lodash');

var serverStore = require('../stores/server');
var channelStore = require('../stores/channel');
var selectedTabStore = require('../stores/selectedTab');
var tabActions = require('../actions/tab');

var TabList = React.createClass({
	mixins: [
		Reflux.connect(serverStore, 'servers'),
		Reflux.connect(channelStore, 'channels'),
		Reflux.connect(selectedTabStore, 'selectedTab')
	],

	getInitialState: function() {
		return {
			servers: serverStore.getState(),
			channels: channelStore.getState(),
			selectedTab: selectedTabStore.getState()
		};
	},

	render: function() {
		var self = this;
		var tabs = _.map(this.state.channels, function(server, address) {
			var channels = _.map(server, function(channel, name) {
				return <p onClick={tabActions.select.bind(null, address, name)}>{name}</p>;
			});
			channels.unshift(<p onClick={tabActions.select.bind(null, address, null)}>{address}</p>);
			return channels;
		});

		return (
			<div className="tablist">{tabs}</div>
		);
	}
});

module.exports = TabList;