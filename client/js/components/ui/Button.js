import React from 'react';
import cn from 'classnames';

const Button = ({ children, category, className, ...props }) => (
  <button
    className={cn(`button-${category}`, className)}
    type="button"
    {...props}
  >
    {children}
  </button>
);

export default Button;
