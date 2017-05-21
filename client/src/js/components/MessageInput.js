import React, { PureComponent } from 'react';

export default class MessageInput extends PureComponent {
  state = {
    value: ''
  };

  handleKey = e => {
    const { tab, runCommand, sendMessage, addInputHistory, incrementInputHistory,
      decrementInputHistory, resetInputHistory, history } = this.props;

    if (e.key === 'Enter' && e.target.value) {
      if (e.target.value[0] === '/') {
        runCommand(e.target.value, tab.name, tab.server);
      } else if (tab.name) {
        sendMessage(e.target.value, tab.name, tab.server);
      }

      addInputHistory(e.target.value);
      resetInputHistory();
      this.setState({ value: '' });
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      incrementInputHistory();
    } else if (e.key === 'ArrowDown') {
      decrementInputHistory();
    } else if (history) {
      this.setState({ value: e.target.value });
      resetInputHistory();
    }
  };

  handleChange = e => {
    this.setState({ value: e.target.value });
  };

  render() {
    const { nick } = this.props;
    return (
      <div className="message-input-wrap">
        <span className="message-input-nick">{nick}</span>
        <input
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
