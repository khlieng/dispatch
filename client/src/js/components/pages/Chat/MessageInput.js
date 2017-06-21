import React, { PureComponent } from 'react';
import Editable from 'components/ui/Editable';

export default class MessageInput extends PureComponent {
  state = {
    value: ''
  };

  handleKey = e => {
    const { tab, onCommand, onMessage,
      add, reset, increment, decrement, currentHistoryEntry } = this.props;

    if (e.key === 'Enter' && e.target.value) {
      if (e.target.value[0] === '/') {
        onCommand(e.target.value, tab.name, tab.server);
      } else if (tab.name) {
        onMessage(e.target.value, tab.name, tab.server);
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
    const { nick, currentHistoryEntry, onNickChange, onNickEditDone } = this.props;
    return (
      <div className="message-input-wrap">
        <Editable
          className="message-input-nick"
          value={nick}
          onBlur={onNickEditDone}
          onChange={onNickChange}
        >
          <span className="message-input-nick">{nick}</span>
        </Editable>
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
