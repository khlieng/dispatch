var React = require('react');
var _ = require('lodash');

var Navicon = require('./Navicon.jsx');
var serverActions = require('../actions/server');
var channelActions = require('../actions/channel');
var PureMixin = require('../mixins/pure');

var Connect = React.createClass({
	mixins: [PureMixin],

	getInitialState() {
		return {
			showOptionals: false
		};
	},

	handleSubmit(e) {
		e.preventDefault();

		var address = e.target.address.value.trim();
		var nick = e.target.nick.value.trim();
		var channels = _.filter(_.map(e.target.channels.value.split(','), _.trim));
		var opts = {
			name: e.target.name.value.trim(),
			tls: e.target.ssl.checked
		};

		if (this.state.showOptionals) {
			opts.realname = e.target.realname.value.trim();
			opts.username = e.target.username.value.trim();
			opts.password = e.target.password.value.trim();
		}

		if (address.indexOf('.') > 0 && nick) {
			serverActions.connect(address, nick, opts);

			if (channels.length > 0) {
				channelActions.join(channels, address);
			}
		}
	},

	handleShowClick: function() {
		this.setState({ showOptionals: !this.state.showOptionals});
	},

	render: function() {
		var optionals = null;

		if (this.state.showOptionals) {
			optionals = (
				<div>
					<input name="username" type="text" placeholder="Username" />
					<input name="password" type="text" placeholder="Password" />
					<input name="realname" type="text" placeholder="Realname" />
				</div>
			);
		}

		return (
			<div className="connect">
				<Navicon />
				<form ref="form" className="connect-form" onSubmit={this.handleSubmit}>
					<h1>Connect</h1>
					<input name="name" type="text" placeholder="Name" defaultValue="Freenode" />
					<input name="address" type="text" placeholder="Address" defaultValue="irc.freenode.net" />
					<input name="nick" type="text" placeholder="Nick" />
					<input name="channels" type="text" placeholder="Channels" />
					{optionals}
					<p>
						<label><input name="ssl" type="checkbox" />SSL</label>
						<i className="icon-ellipsis" onClick={this.handleShowClick}></i>
					</p>
					<input type="submit" value="Connect" />
				</form>
			</div>
		);
	}
});

module.exports = Connect;