import React, { Component } from 'react';
import { createSelector } from 'reselect';
import { Form, withFormik } from 'formik';
import Navicon from 'components/ui/Navicon';
import Button from 'components/ui/Button';
import Checkbox from 'components/ui/formik/Checkbox';
import TextInput from 'components/ui/TextInput';
import Error from 'components/ui/formik/Error';
import { isValidNick, isValidChannel, isValidUsername, isInt } from 'utils';
import { FiMoreHorizontal } from 'react-icons/fi';

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
    this.setState(prevState => ({ showOptionals: !prevState.showOptionals }));
  };

  renderOptionals = () => {
    const { hexIP } = this.props;

    return (
      <div>
        {!hexIP && <TextInput name="username" />}
        <TextInput name="password" type="password" noTrim />
        <TextInput name="realname" noTrim />
      </div>
    );
  };

  transformPort = port => {
    if (!port) {
      return this.props.values.tls ? 6697 : 6667;
    }
    return port;
  };

  transformChannels = channels => {
    const comma = channels[channels.length - 1] === ',';

    channels = channels
      .split(',')
      .map(channel => {
        channel = channel.trim();
        if (channel) {
          if (isValidChannel(channel, false) && channel[0] !== '#') {
            channel = `#${channel}`;
          }
        }
        return channel;
      })
      .filter(s => s)
      .join(',');

    return comma ? `${channels},` : channels;
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
          <TextInput name="nick" />
          <Button type="submit">Connect</Button>
        </Form>
      );
    } else {
      form = (
        <Form className="connect-form">
          <h1>Connect</h1>
          <TextInput name="name" autoCapitalize="words" noTrim />
          <div className="connect-form-address">
            <TextInput name="host" noError />
            <TextInput
              name="port"
              type="number"
              blurTransform={this.transformPort}
              noError
            />
            <Checkbox
              name="tls"
              label="SSL"
              topLabel
              onChange={this.handleSSLChange}
            />
          </div>
          <Error name="host" />
          <Error name="port" />
          <TextInput name="nick" />
          <TextInput name="channels" transform={this.transformChannels} />
          {this.state.showOptionals && this.renderOptionals()}
          <Button
            className="connect-form-button-optionals"
            icon={FiMoreHorizontal}
            onClick={this.handleShowClick}
          />
          <Button type="submit">Connect</Button>
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
  enableReinitialize: true,
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
      tls: defaults.ssl || false
    };
  },
  validate: values => {
    const errors = {};

    if (!values.host) {
      errors.host = 'Host is required';
    } else if (values.host.indexOf('.') < 1) {
      errors.host = 'Invalid host';
    }

    if (!isInt(values.port, 1, 65535)) {
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

    const channels = values.channels.split(',');
    for (let i = channels.length - 1; i >= 0; i--) {
      if (i === channels.length - 1 && channels[i] === '') {
        /* eslint-disable-next-line no-continue */
        continue;
      }

      if (!isValidChannel(channels[i])) {
        errors.channels = 'Invalid channel(s)';
        break;
      }
    }

    return errors;
  },
  handleSubmit: (values, { props }) => {
    const { connect, select, join } = props;
    const channels = values.channels ? values.channels.split(',') : [];
    delete values.channels;

    values.port = `${values.port}`;
    connect(values);
    select(values.host);

    if (channels.length > 0) {
      join(channels, values.host);
    }
  }
})(Connect);
