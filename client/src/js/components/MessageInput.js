import React, { Component } from 'react';
import pure from 'pure-render-decorator';

@pure
export default class MessageInput extends Component {
  state = {
    value: ''
  };

  handleKey = e => {
    const { tab, runCommand, sendMessage, addInputHistory, incrementInputHistory,
      decrementInputHistory, resetInputHistory } = this.props;

    if (e.which === 13 && e.target.value) {
      if (e.target.value[0] === '/') {
        runCommand(e.target.value, tab.channel || tab.user, tab.server);
      } else if (tab.channel) {
        sendMessage(e.target.value, tab.channel, tab.server);
      } else if (tab.user) {
        sendMessage(e.target.value, tab.user, tab.server);
      }

      addInputHistory(e.target.value);
      resetInputHistory();
      this.setState({ value: '' });
    } else if (e.which === 38) {
      e.preventDefault();
      incrementInputHistory();
    } else if (e.which === 40) {
      decrementInputHistory();
    } else if (e.key === 'Backspace' || e.key === 'Delete') {
      resetInputHistory();
    } else if (e.key === 'Unidentified') {
      this.setState({ value: e.target.value });
      resetInputHistory();
    }
  };

  handleChange = e => {
    this.setState({ value: e.target.value });
  };

  render() {
    return (
      <div className="message-input-wrap">
        <input
          ref="input"
          className="message-input"
          type="text"
          value={this.props.history || this.state.value}
          onKeyDown={this.handleKey}
          onChange={this.handleChange}
        />
      </div>
    );
  }
}
