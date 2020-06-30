import React from 'react';
import Modal from 'react-modal';
import { useSelector } from 'react-redux';
import { FiX } from 'react-icons/fi';
import Text from 'components/Text';
import Button from 'components/ui/Button';
import useModal from 'components/modals/useModal';
import { getSelectedChannel } from 'state/channels';
import { linkify } from 'utils';

const Topic = () => {
  const [modal, channel, closeModal] = useModal('topic');

  const topic = useSelector(state => getSelectedChannel(state)?.topic);

  return (
    <Modal {...modal}>
      <div className="modal-header">
        <h2>Topic in {channel}</h2>
        <Button icon={FiX} className="modal-close" onClick={closeModal} />
      </div>
      <p className="modal-content">
        <Text>{linkify(topic)}</Text>
      </p>
    </Modal>
  );
};

export default Topic;
