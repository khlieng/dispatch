import React from 'react';
import Reflux from 'reflux';
import Infinite from 'react-infinite';
import MessageHeader from './MessageHeader.jsx';
import MessageLine from './MessageLine.jsx';
import messageLineStore from '../stores/messageLine';
import selectedTabStore from '../stores/selectedTab';
import messageActions from '../actions/message';
import PureMixin from '../mixins/pure';

export default React.createClass({
	mixins: [
		PureMixin,
		Reflux.connect(messageLineStore, 'messages'),
		Reflux.connect(selectedTabStore, 'selectedTab')
	],

	getInitialState() {
		return {
			height: window.innerHeight - 100
		};
	},

	componentDidMount() {
		this.updateWidth();
		window.addEventListener('resize', this.handleResize);
	},

	componentWillUnmount() {
		window.removeEventListener('resize', this.handleResize);
	},

	componentWillUpdate() {
		var el = this.refs.list.refs.scrollable;
		this.autoScroll = el.scrollTop + el.offsetHeight === el.scrollHeight;
	},

	componentDidUpdate() {
		setTimeout(this.updateWidth, 0);

		if (this.autoScroll) {
			var el = this.refs.list.refs.scrollable;
			el.scrollTop = el.scrollHeight;
		}
	},

	handleResize() {
		this.updateWidth();
		this.setState({ height: window.innerHeight - 100 });
	},

	updateWidth() {
		const { list } = this.refs;
		if (list) {
			const width = list.refs.scrollable.offsetWidth - 30;
			if (this.width !== width) {
				this.width = width;
				messageActions.setWrapWidth(width);
			}
		}
	},

	render() {
		const tab = this.state.selectedTab;
		const dest = tab.channel || tab.server;
		const lines = [];

		this.state.messages.forEach((message, j) => {
			const key = message.server + dest + j;

			lines.push(<MessageHeader key={key} message={message} />);

			for (let i = 1; i < message.lines.length; i++) {
				lines.push(
					<MessageLine key={key + '-' + i} type={message.type} line={message.lines[i]} />
				);
			}
		});

		return (
			<div className="messagebox">
				<Infinite
					ref="list"
					className="messagebox-scrollable"
					containerHeight={this.state.height}
					elementHeight={24}
					displayBottomUpwards={false}>
					{lines}
				</Infinite>
			</div>
		);
	}
});
