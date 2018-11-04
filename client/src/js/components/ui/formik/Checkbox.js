import React, { memo } from 'react';
import { FastField } from 'formik';
import Checkbox from 'components/ui/Checkbox';

const FormikCheckbox = ({ name, onChange, ...props }) => (
  <FastField
    name={name}
    render={({ field, form }) => {
      return (
        <Checkbox
          name={name}
          checked={field.value}
          onChange={e => {
            form.setFieldTouched(name, true);
            field.onChange(e);
            if (onChange) {
              onChange(e);
            }
          }}
          {...props}
        />
      );
    }}
  />
);

export default memo(FormikCheckbox);
