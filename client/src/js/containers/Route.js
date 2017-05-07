import { PureComponent } from 'react';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';

class Route extends PureComponent {
  render() {
    if (this.props.route === this.props.name) {
      return this.props.children;
    }
    return null;
  }
}

const getRoute = state => state.router.route;

const mapStateToProps = createStructuredSelector({
  route: getRoute
});

export default connect(mapStateToProps)(Route);
