var shallowEqual = require('react-pure-render/shallowEqual');

module.exports = {
	shouldComponentUpdate(nextProps, nextState) {
		if (this.context.router) {
			var changed = this.pureComponentLastPath !== this.context.router.getCurrentPath();
			this.pureComponentLastPath = this.context.router.getCurrentPath();

			if (changed) {
				return true;
			}
		}

		return !shallowEqual(this.props, nextProps) ||
			!shallowEqual(this.state, nextState);
	}
};