import React, { useCallback } from 'react';
import withModal from 'components/modals/withModal';
import Button from 'components/ui/Button';

const Confirm = ({
  payload: { question, confirmation, onConfirm },
  onClose
}) => {
  const handleConfirm = useCallback(() => {
    onClose(false);
    onConfirm();
  }, []);

  return (
    <>
      <p>{question}</p>
      <Button onClick={handleConfirm}>{confirmation || 'OK'}</Button>
      <Button category="normal" onClick={onClose}>
        Cancel
      </Button>
    </>
  );
};

export default withModal({
  name: 'confirm'
})(Confirm);
