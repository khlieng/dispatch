import React from 'react';
import { FiMenu } from 'react-icons/fi';
import { useDispatch } from 'react-redux';
import Button from 'components/ui/Button';
import { toggleMenu } from 'state/ui';

const Navicon = () => {
  const dispatch = useDispatch();

  return (
    <Button
      className="navicon"
      icon={FiMenu}
      onClick={() => dispatch(toggleMenu())}
    />
  );
};

export default Navicon;
