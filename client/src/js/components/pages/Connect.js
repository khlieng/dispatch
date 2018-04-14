import React, { Component } from 'react';
import Navicon from 'containers/Navicon';

export default class Connect extends Component {
  state = {
    showOptionals: false,
    passwordTouched: false
  };

  handleSubmit = e => {
    const { connect, select, join, defaults } = this.props;

    e.preventDefault();

    const nick = e.target.nick.value.trim();
    let { address, channels } = defaults;
    const opts = {
      name: defaults.name
    };

    if (!defaults.readonly) {
      address = e.target.address.value.trim();
      channels = e.target.channels.value
        .split(',')
        .map(s => s.trim())
        .filter(s => s);
      opts.name = e.target.name.value.trim();
      opts.tls = e.target.ssl.checked;

      if (this.state.showOptionals) {
        opts.realname = e.target.realname.value.trim();
        opts.username = e.target.username.value.trim();

        if (this.state.passwordTouched) {
          opts.password = e.target.password.value.trim();
        }
      }
    }

    if (address.indexOf('.') > 0 && nick) {
      connect(address, nick, opts);

      const i = address.indexOf(':');
      if (i > 0) {
        address = address.slice(0, i);
      }

      select(address);

      if (channels.length > 0) {
        join(channels, address);
      }
    }
  };

  handleShowClick = () => {
    this.setState({ showOptionals: !this.state.showOptionals });
  };

  handlePasswordChange = () => {
    this.setState({ passwordTouched: true });
  };

  render() {
    const { defaults } = this.props;
    let optionals = null;

    if (this.state.showOptionals) {
      optionals = (
        <div>
          <input name="username" type="text" placeholder="Username" />
          <input
            name="password"
            type="password"
            placeholder="Password"
            defaultValue={defaults.password ? '      ' : null}
            onChange={this.handlePasswordChange}
          />
          <input name="realname" type="text" placeholder="Realname" />
        </div>
      );
    }

    let form;

    if (defaults.readonly) {
      form = (
        <form className="connect-form" onSubmit={this.handleSubmit}>
          <h1>Connect</h1>
          {defaults.showDetails && (
            <div className="connect-details">
              <h2>{defaults.address}</h2>
              {defaults.channels.sort().map(channel => <p>{channel}</p>)}
            </div>
          )}
          <input name="nick" type="text" placeholder="Nick" />
          <input type="submit" value="Connect" />
        </form>
      );
    } else {
      form = (
        <form className="connect-form" onSubmit={this.handleSubmit}>
          <h1>Connect</h1>
          <input
            name="name"
            type="text"
            placeholder="Name"
            defaultValue={defaults.name}
          />
          <input
            name="address"
            type="text"
            placeholder="Address"
            defaultValue={defaults.address}
          />
          <input name="nick" type="text" placeholder="Nick" />
          <input
            name="channels"
            type="text"
            placeholder="Channels"
            defaultValue={
              defaults.channels ? defaults.channels.join(',') : null
            }
          />
          {optionals}
          <p>
            <label htmlFor="ssl">
              <input name="ssl" type="checkbox" defaultChecked={defaults.ssl} />SSL
            </label>
            <i className="icon-ellipsis" onClick={this.handleShowClick} />
          </p>
          <input type="submit" value="Connect" />
        </form>
      );
    }

    return (
      <div className="connect">
        <Navicon />
        {form}
      </div>
    );
  }
}
