import React from 'react';
import { ErrorMessage } from 'formik';

const Error = props => (
  <ErrorMessage component="div" className="form-error" {...props} />
);

export default Error;
