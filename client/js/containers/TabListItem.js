import { createStructuredSelector } from 'reselect';
import get from 'lodash/get';
import TabListItem from 'components/TabListItem';
import connect from 'utils/connect';

const mapState = createStructuredSelector({
  error: (state, { network, target }) => {
    const messages = get(state, ['messages', network, target]);

    if (messages && messages.length > 0) {
      return messages[messages.length - 1].type === 'error';
    }
    return false;
  }
});

export default connect(mapState)(TabListItem);
