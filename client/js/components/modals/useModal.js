import { useCallback } from 'react';
import Modal from 'react-modal';
import { useSelector, useDispatch } from 'react-redux';
import { closeModal } from 'state/modals';

Modal.setAppElement('#root');

const defaultPayload = {};

export default function useModal(name) {
  const isOpen = useSelector(state => state.modals[name]?.isOpen || false);
  const payload = useSelector(
    state => state.modals[name]?.payload || defaultPayload
  );
  const dispatch = useDispatch();

  const handleRequestClose = useCallback(
    (dismissed = true) => {
      dispatch(closeModal(name));

      if (dismissed && payload.onDismiss) {
        payload.onDismiss();
      }
    },
    [payload.onDismiss]
  );

  const modalProps = {
    isOpen,
    contentLabel: name,
    onRequestClose: handleRequestClose,
    className: {
      base: `modal modal-${name}`,
      afterOpen: 'modal-opening',
      beforeClose: 'modal-closing'
    },
    overlayClassName: {
      base: 'modal-overlay',
      afterOpen: 'modal-overlay-opening',
      beforeClose: 'modal-overlay-closing'
    },
    closeTimeoutMS: 200
  };

  return [modalProps, payload, handleRequestClose];
}
