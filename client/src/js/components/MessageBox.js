import React, { Component } from 'react';
import { VirtualScroll } from 'react-virtualized';
import pure from 'pure-render-decorator';
import Message from './Message';

@pure
export default class MessageBox extends Component {
  state = {
    height: window.innerHeight - 100
  };

  componentDidMount() {
    this.updateWidth();
    window.addEventListener('resize', this.handleResize);
  }

  componentWillReceiveProps() {
    const el = this.list.refs.scrollingContainer;
    this.autoScroll = el.scrollTop + el.offsetHeight === el.scrollHeight;
  }

  componentWillUpdate(nextProps) {
    if (nextProps.messages !== this.props.messages) {
      this.list.recomputeRowHeights();
    }
  }

  componentDidUpdate() {
    this.updateWidth();

    if (this.autoScroll) {
      const el = this.list.refs.scrollingContainer;
      el.scrollTop = el.scrollHeight;
    }
  }

  componentWillUnmount() {
    window.removeEventListener('resize', this.handleResize);
  }

  getRowHeight = index => {
    const { messages } = this.props;

    if (index === 0 || index === messages.size + 1) {
      return 7;
    }

    return messages.get(index - 1).height;
  };

  updateWidth = resize => {
    const { isChannel, setWrapWidth, updateMessageHeight } = this.props;
    if (this.list) {
      let width = this.list.refs.scrollingContainer.clientWidth - 30;

      if (isChannel) {
        width += 200;
      }

      if (this.width !== width) {
        this.width = width;
        setWrapWidth(width);

        if (resize) {
          updateMessageHeight();
        }
      }
    }
  };

  handleResize = () => {
    this.updateWidth(true);
    this.setState({ height: window.innerHeight - 100 });
  };

  renderMessage = index => {
    const { messages } = this.props;

    if (index === 0 || index === messages.size + 1) {
      return <span style={{ height: '7px' }} />;
    }

    const { select, openPrivateChat } = this.props;
    const message = messages.get(index - 1);

    return (
      <Message
        message={message}
        select={select}
        openPrivateChat={openPrivateChat}
      />
    );
  };

  render() {
    return (
      <div className="messagebox">
        <VirtualScroll
          ref={el => { this.list = el; }}
          height={this.state.height}
          rowsCount={this.props.messages.size + 2}
          rowHeight={this.getRowHeight}
          rowRenderer={this.renderMessage}
        />
      </div>
    );
  }
}
