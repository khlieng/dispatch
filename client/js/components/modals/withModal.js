import React, { useCallback } from 'react';
import Modal from 'react-modal';
import { createSelector } from 'reselect';
import { getModals, closeModal } from 'state/modals';
import connect from 'utils/connect';

Modal.setAppElement('#root');

export default function withModal({ name, ...modalProps }) {
  modalProps = {
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
    closeTimeoutMS: 200,
    ...modalProps
  };

  return WrappedComponent => {
    const ReduxModal = ({ onRequestClose, ...props }) => {
      const handleRequestClose = useCallback(
        (dismissed = true) => {
          onRequestClose();

          if (dismissed && props.payload.onDismiss) {
            props.payload.onDismiss();
          }
        },
        [props.payload.onDismiss]
      );

      return (
        <Modal
          contentLabel={name}
          onRequestClose={handleRequestClose}
          {...modalProps}
          {...props}
        >
          <WrappedComponent onClose={handleRequestClose} {...props} />
        </Modal>
      );
    };

    const mapState = createSelector(
      getModals,
      modals => modals[name] || { payload: {} }
    );

    const mapDispatch = dispatch => ({
      onRequestClose: () => dispatch(closeModal(name))
    });

    return connect(
      mapState,
      mapDispatch
    )(ReduxModal);
  };
}
