import { connect } from 'react-redux';
import Navicon from 'components/ui/Navicon';
import { toggleMenu } from 'state/ui';

const mapDispatch = {
  onClick: toggleMenu
};

export default connect(null, mapDispatch)(Navicon);
