import React from 'react';
import Modal from 'react-modal';
import useModal from 'components/modals/useModal';
import Button from 'components/ui/Button';

const Confirm = () => {
  const [modal, payload, closeModal] = useModal('confirm');
  const { question, confirmation, onConfirm } = payload;

  const handleConfirm = () => {
    closeModal(false);
    onConfirm();
  };

  return (
    <Modal {...modal}>
      <p>{question}</p>
      <Button onClick={handleConfirm}>{confirmation || 'OK'}</Button>
      <Button category="normal" onClick={closeModal}>
        Cancel
      </Button>
    </Modal>
  );
};

export default Confirm;
