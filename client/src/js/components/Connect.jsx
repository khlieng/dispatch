var React = require('react');
var _ = require('lodash');

var serverActions = require('../actions/server');
var channelActions = require('../actions/channel');

var Connect = React.createClass({
	handleSubmit: function(e) {
		e.preventDefault();

		var name = e.target.name.value.trim();
		var address = e.target.address.value.trim();
		var ssl = e.target.ssl.checked;
		var nick = e.target.nick.value.trim();
		var username = e.target.username.value.trim();
		var channels = _.filter(_.map(e.target.channels.value.split(','), _.trim));

		if (address.indexOf('.') > 0 && nick && username) {
			serverActions.connect(address, nick, username, ssl, name);

			if (channels.length > 0) {
				channelActions.join(channels, address);				
			}
		}
	},

	render: function() {
		return (
			<div className="connect">
				<form ref="form" className="connect-form" onSubmit={this.handleSubmit}>
					<h1>Connect</h1>
					<input name="name" type="text" placeholder="Name" defaultValue="Freenode" />
					<input name="address" type="text" placeholder="Address" defaultValue="irc.freenode.net" />
					<label><input name="ssl" type="checkbox" />SSL</label>
					<input name="nick" type="text" placeholder="Nick" onChange={this.handleNickChange} />
					<input name="username" type="text" placeholder="Username" />
					<input name="channels" type="text" placeholder="Channels" />
					<input type="submit" value="Connect" />
				</form>
			</div>
		);
	}
});

module.exports = Connect;