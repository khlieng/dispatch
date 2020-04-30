import Navicon from 'components/ui/Navicon';
import { toggleMenu } from 'state/ui';
import connect from 'utils/connect';

const mapDispatch = {
  onClick: toggleMenu
};

export default connect(null, mapDispatch)(Navicon);
