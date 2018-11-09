import { createStructuredSelector } from 'reselect';
import connect from 'utils/connect';

const Route = ({ route, name, children }) => {
  if (route === name) {
    return children;
  }
  return null;
};

const getRoute = state => state.router.route;

const mapState = createStructuredSelector({
  route: getRoute
});

export default connect(mapState)(Route);
