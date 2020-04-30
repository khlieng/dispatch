import { connect } from 'react-redux';

const strictEqual = (a, b) => a === b;

export default (mapState, mapDispatch) =>
  connect(mapState, mapDispatch, null, {
    areStatePropsEqual: strictEqual
  });
