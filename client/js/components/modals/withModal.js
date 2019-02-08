import React, { useCallback } from 'react';
import Modal from 'react-modal';
import { createStructuredSelector } from 'reselect';
import get from 'lodash/get';
import { getModals, closeModal } from 'state/modals';
import connect from 'utils/connect';
import { bindActionCreators } from 'redux';

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

    const mapState = createStructuredSelector({
      isOpen: state => get(getModals(state), [name, 'isOpen'], false),
      payload: state => get(getModals(state), [name, 'payload'], {}),
      ...modalProps.state
    });

    const mapDispatch = dispatch => {
      const actions = { onRequestClose: () => dispatch(closeModal(name)) };
      if (modalProps.actions) {
        return {
          ...actions,
          ...bindActionCreators(modalProps.actions, dispatch)
        };
      }
      return actions;
    };

    return connect(
      mapState,
      mapDispatch
    )(ReduxModal);
  };
}
