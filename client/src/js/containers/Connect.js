import React, { Component } from 'react';
import { connect } from 'react-redux';
import pure from 'pure-render-decorator';
import Navicon from '../components/Navicon';
import * as serverActions from '../actions/server';
import { join } from '../actions/channel';
import { select } from '../actions/tab';

@pure
class Connect extends Component {
  state = {
    showOptionals: false
  };

  handleSubmit = (e) => {
    e.preventDefault();

    const { dispatch } = this.props;
    let address = e.target.address.value.trim();
    const nick = e.target.nick.value.trim();
    const channels = e.target.channels.value.split(',').map(s => s.trim()).filter(s => s);
    const opts = {
      name: e.target.name.value.trim(),
      tls: e.target.ssl.checked
    };

    if (this.state.showOptionals) {
      opts.realname = e.target.realname.value.trim();
      opts.username = e.target.username.value.trim();
      opts.password = e.target.password.value.trim();
    }

    if (address.indexOf('.') > 0 && nick) {
      dispatch(serverActions.connect(address, nick, opts));

      const i = address.indexOf(':');
      if (i > 0) {
        address = address.slice(0, i);
      }

      dispatch(select(address));

      if (channels.length > 0) {
        dispatch(join(channels, address));
      }
    }
  };

  handleShowClick = () => {
    this.setState({ showOptionals: !this.state.showOptionals });
  };

  render() {
    const defaults = window.__ENV__.defaults;
    let optionals = null;

    if (this.state.showOptionals) {
      optionals = (
        <div>
          <input name="username" type="text" placeholder="Username" />
          <input
            name="password"
            type="password"
            placeholder="Password"
            defaultValue={defaults.password}
          />
          <input name="realname" type="text" placeholder="Realname" />
        </div>
      );
    }

    return (
      <div className="connect">
        <Navicon />
        <form ref="form" className="connect-form" onSubmit={this.handleSubmit}>
          <h1>Connect</h1>
          <input name="name" type="text" placeholder="Name" defaultValue={defaults.name} />
          <input name="address" type="text" placeholder="Address" defaultValue={defaults.address} />
          <input name="nick" type="text" placeholder="Nick" />
          <input
            name="channels"
            type="text"
            placeholder="Channels"
            defaultValue={defaults.channels ? defaults.channels.join(',') : null}
          />
          {optionals}
          <p>
            <label><input name="ssl" type="checkbox" defaultChecked={defaults.ssl} />SSL</label>
            <i className="icon-ellipsis" onClick={this.handleShowClick}></i>
          </p>
          <input type="submit" value="Connect" />
        </form>
      </div>
    );
  }
}

export default connect()(Connect);
