import React, { Component } from 'react';
import Infinite from 'react-infinite';
import pure from 'pure-render-decorator';
import MessageHeader from './MessageHeader';
import MessageLine from './MessageLine';

@pure
export default class MessageBox extends Component {
  state = {
    height: window.innerHeight - 100
  };

  componentDidMount() {
    this.updateWidth();
    window.addEventListener('resize', this.handleResize);
  }

  componentWillUpdate() {
    const el = this.refs.list.refs.scrollable;
    this.autoScroll = el.scrollTop + el.offsetHeight === el.scrollHeight;
  }

  componentDidUpdate() {
    setTimeout(this.updateWidth, 0);

    if (this.autoScroll) {
      const el = this.refs.list.refs.scrollable;
      el.scrollTop = el.scrollHeight;
    }
  }

  componentWillUnmount() {
    window.removeEventListener('resize', this.handleResize);
  }

  updateWidth = () => {
    const { setWrapWidth } = this.props;
    const { list } = this.refs;
    if (list) {
      const width = list.refs.scrollable.offsetWidth - 30;
      if (this.width !== width) {
        this.width = width;
        setWrapWidth(width);
      }
    }
  };

  handleResize = () => {
    this.updateWidth();
    this.setState({ height: window.innerHeight - 100 });
  };

  render() {
    const { tab, messages, select, openPrivateChat } = this.props;
    const dest = tab.channel || tab.user || tab.server;
    const lines = [];

    messages.forEach((message, j) => {
      const key = message.server + dest + j;
      lines.push(
        <MessageHeader
          key={key}
          message={message}
          select={select}
          openPrivateChat={openPrivateChat}
        />
      );

      for (let i = 1; i < message.lines.length; i++) {
        lines.push(
          <MessageLine key={`${key}-${i}`} type={message.type} line={message.lines[i]} />
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
          displayBottomUpwards={false}
        >
          {lines}
        </Infinite>
      </div>
    );
  }
}
