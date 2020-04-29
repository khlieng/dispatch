import React from 'react';
import cn from 'classnames';

const Button = ({ children, category, className, icon: Icon, ...props }) => (
  <button
    className={cn(
      {
        [`button-${category}`]: category,
        'icon-button': Icon && !children
      },
      className
    )}
    type="button"
    {...props}
  >
    {Icon && <Icon />}
    {children}
  </button>
);

export default Button;
