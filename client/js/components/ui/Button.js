import React from 'react';

const Button = ({ children, category, ...props }) => (
  <button className={`button-${category}`} type="button" {...props}>
    {children}
  </button>
);

export default Button;
