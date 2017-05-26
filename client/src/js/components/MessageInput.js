import React, { PureComponent } from 'react';

export default class MessageInput extends PureComponent {
  state = {
    value: ''
  };

  handleKey = e => {
    const { tab, runCommand, sendMessage,
      add, reset, increment, decrement, currentHistoryEntry } = this.props;

    if (e.key === 'Enter' && e.target.value) {
      if (e.target.value[0] === '/') {
        runCommand(e.target.value, tab.name, tab.server);
      } else if (tab.name) {
        sendMessage(e.target.value, tab.name, tab.server);
      }

      add(e.target.value);
      reset();
      this.setState({ value: '' });
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      increment();
    } else if (e.key === 'ArrowDown') {
      decrement();
    } else if (currentHistoryEntry) {
      this.setState({ value: e.target.value });
      reset();
    }
  };

  handleChange = e => {
    this.setState({ value: e.target.value });
  };

  render() {
    const { nick, currentHistoryEntry } = this.props;
    return (
      <div className="message-input-wrap">
        <span className="message-input-nick">{nick}</span>
        <input
          className="message-input"
          type="text"
          value={currentHistoryEntry || this.state.value}
          onKeyDown={this.handleKey}
          onChange={this.handleChange}
        />
      </div>
    );
  }
}
