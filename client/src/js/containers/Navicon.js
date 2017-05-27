import { connect } from 'react-redux';
import Navicon from '../components/Navicon';
import { toggleMenu } from '../state/ui';

const mapDispatch = { toggleMenu };

export default connect(null, mapDispatch)(Navicon);
