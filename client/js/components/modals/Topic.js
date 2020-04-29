import React from 'react';
import { FiX } from 'react-icons/fi';
import Button from 'components/ui/Button';
import withModal from 'components/modals/withModal';
import { linkify } from 'utils';

const Topic = ({ payload: { topic, channel }, onClose }) => {
  return (
    <>
      <div className="modal-header">
        <h2>Topic in {channel}</h2>
        <Button icon={FiX} className="modal-close" onClick={onClose} />
      </div>
      <p className="modal-content">{linkify(topic)}</p>
    </>
  );
};

export default withModal({
  name: 'topic'
})(Topic);
