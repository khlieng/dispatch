import React from 'react';
import classnames from 'classnames';

const Checkbox = ({ name, label, topLabel, classNameLabel, ...props }) => (
  <label
    className={classnames('checkbox', classNameLabel, {
      'top-label': topLabel
    })}
    htmlFor={name}
  >
    {topLabel && label}
    <input type="checkbox" id={name} name={name} {...props} />
    <span />
    {!topLabel && label}
  </label>
);

export default Checkbox;
