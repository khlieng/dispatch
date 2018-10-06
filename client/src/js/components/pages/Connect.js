import React, { Component } from 'react';
import { createSelector } from 'reselect';
import { Form, withFormik } from 'formik';
import Navicon from 'containers/Navicon';
import Checkbox from 'components/ui/Checkbox';
import TextInput from 'components/ui/TextInput';
import Error from 'components/ui/formik/Error';
import { isValidNick, isValidChannel, isValidUsername, isInt } from 'utils';

const getSortedDefaultChannels = createSelector(
  defaults => defaults.channels,
  channels => channels.split(',').sort()
);

class Connect extends Component {
  state = {
    showOptionals: false
  };

  handleSSLChange = e => {
    const { values, setFieldValue } = this.props;
    if (e.target.checked && values.port === 6667) {
      setFieldValue('port', 6697, false);
    } else if (!e.target.checked && values.port === 6697) {
      setFieldValue('port', 6667, false);
    }
  };

  handleShowClick = () => {
    this.setState({ showOptionals: !this.state.showOptionals });
  };

  renderOptionals = () => {
    const { hexIP } = this.props;

    return (
      <div>
        {!hexIP && [
          <TextInput name="username" placeholder="Username" />,
          <Error name="username" />
        ]}
        <TextInput type="password" name="password" placeholder="Password" />
        <TextInput name="realname" placeholder="Realname" />
      </div>
    );
  };

  render() {
    const { defaults, values } = this.props;
    const { readOnly, showDetails } = defaults;
    let form;

    if (readOnly) {
      form = (
        <Form className="connect-form">
          <h1>Connect</h1>
          {showDetails && (
            <div className="connect-details">
              <h2>
                {values.host}:{values.port}
              </h2>
              {getSortedDefaultChannels(values).map(channel => (
                <p>{channel}</p>
              ))}
            </div>
          )}
          <TextInput name="nick" placeholder="Nick" />
          <Error name="nick" />
          <button>Connect</button>
        </Form>
      );
    } else {
      form = (
        <Form className="connect-form">
          <h1>Connect</h1>
          <TextInput name="name" placeholder="Name" autoCapitalize="words" />
          <div className="connect-form-address">
            <TextInput name="host" placeholder="Host" />
            <TextInput name="port" type="number" placeholder="Port" />
            <Checkbox name="tls" label="SSL" onChange={this.handleSSLChange} />
          </div>
          <Error name="host" />
          <Error name="port" />
          <TextInput name="nick" placeholder="Nick" />
          <Error name="nick" />
          <TextInput name="channels" placeholder="Channels" />
          <Error name="channels" />
          {this.state.showOptionals && this.renderOptionals()}
          <i className="icon-ellipsis" onClick={this.handleShowClick} />
          <button>Connect</button>
        </Form>
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

export default withFormik({
  mapPropsToValues: ({ defaults }) => {
    let port = 6667;
    if (defaults.port) {
      ({ port } = defaults);
    } else if (defaults.ssl) {
      port = 6697;
    }

    return {
      name: defaults.name,
      host: defaults.host,
      port,
      nick: '',
      channels: defaults.channels.join(','),
      username: '',
      password: defaults.password ? '      ' : '',
      realname: '',
      tls: defaults.ssl
    };
  },
  validate: values => {
    Object.keys(values).forEach(k => {
      if (typeof values[k] === 'string') {
        values[k] = values[k].trim();
      }
    });

    const errors = {};

    if (!values.host) {
      errors.host = 'Host is required';
    } else if (values.host.indexOf('.') < 1) {
      errors.host = 'Invalid host';
    }

    if (!values.port) {
      values.port = values.tls ? 6697 : 6667;
    } else if (!isInt(values.port, 1, 65535)) {
      errors.port = 'Invalid port';
    }

    if (!values.nick) {
      errors.nick = 'Nick is required';
    } else if (!isValidNick(values.nick)) {
      errors.nick = 'Invalid nick';
    }

    if (values.username && !isValidUsername(values.username)) {
      errors.username = 'Invalid username';
    }

    values.channels = values.channels
      .split(',')
      .map(channel => {
        channel = channel.trim();
        if (channel) {
          if (isValidChannel(channel, false)) {
            if (channel[0] !== '#') {
              channel = `#${channel}`;
            }
          } else {
            errors.channels = 'Invalid channel(s)';
          }
        }
        return channel;
      })
      .filter(s => s)
      .join(',');

    return errors;
  },
  handleSubmit: (values, { props }) => {
    const { connect, select, join } = props;
    const channels = values.channels.split(',');
    delete values.channels;

    values.port = `${values.port}`;
    connect(values);
    select(values.host);

    if (channels.length > 0) {
      join(channels, values.host);
    }
  }
})(Connect);
