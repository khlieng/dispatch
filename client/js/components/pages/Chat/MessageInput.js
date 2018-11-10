import React, { memo, useState } from 'react';
import classnames from 'classnames';
import Editable from 'components/ui/Editable';
import { isValidNick } from 'utils';

const MessageInput = ({
  nick,
  currentHistoryEntry,
  onNickChange,
  onNickEditDone,
  tab,
  onCommand,
  onMessage,
  add,
  reset,
  increment,
  decrement
}) => {
  const [value, setValue] = useState('');

  const handleKey = e => {
    if (e.key === 'Enter' && e.target.value) {
      if (e.target.value[0] === '/') {
        onCommand(e.target.value, tab.name, tab.server);
      } else if (tab.name) {
        onMessage(e.target.value, tab.name, tab.server);
      }

      add(e.target.value);
      reset();
      setValue('');
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      increment();
    } else if (e.key === 'ArrowDown') {
      decrement();
    } else if (currentHistoryEntry) {
      setValue(e.target.value);
      reset();
    }
  };

  const handleChange = e => setValue(e.target.value);

  return (
    <div className="message-input-wrap">
      <Editable
        className={classnames('message-input-nick', {
          invalid: !isValidNick(nick)
        })}
        value={nick}
        onBlur={onNickEditDone}
        onChange={onNickChange}
      >
        <span className="message-input-nick">{nick}</span>
      </Editable>
      <input
        className="message-input"
        type="text"
        value={currentHistoryEntry || value}
        onKeyDown={handleKey}
        onChange={handleChange}
      />
    </div>
  );
};

export default memo(MessageInput);
