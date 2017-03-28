import React, { PureComponent } from 'react';
import { List } from 'react-virtualized/dist/commonjs/List';
import { AutoSizer } from 'react-virtualized/dist/commonjs/AutoSizer';
import Message from './Message';
import { scrollBarWidth } from '../util';

const sbWidth = scrollBarWidth();
const listStyle = { padding: '7px 0', boxSizing: 'content-box' };

export default class MessageBox extends PureComponent {
  componentWillReceiveProps(nextProps) {
    if (nextProps.tab !== this.props.tab) {
      this.bottom = true;
    }
  }

  componentWillUpdate(nextProps) {
    if (nextProps.messages !== this.props.messages) {
      this.list.recomputeRowHeights();
    }
  }

  componentDidUpdate() {
    if (this.bottom) {
      this.list.scrollToRow(this.props.messages.size);
    }

    this.updateWidth();
  }

  getRowHeight = ({ index }) => this.props.messages.get(index).height;

  listRef = el => { this.list = el; };

  updateWidth = (width) => {
    const { isChannel, setWrapWidth, updateMessageHeight } = this.props;
    let wrapWidth = width || this.width;

    if (width) {
      if (isChannel && window.innerWidth > 600) {
        wrapWidth += 200;
      }

      this.width = wrapWidth;
    }

    // eslint-disable-next-line no-underscore-dangle
    const c = this.list.Grid._scrollingContainer;
    if (c.scrollHeight > c.clientHeight) {
      wrapWidth -= sbWidth;
    }

    if (this.wrapWidth !== wrapWidth) {
      this.wrapWidth = wrapWidth;
      setWrapWidth(wrapWidth);
      updateMessageHeight();
    }
  };

  handleResize = size => {
    this.updateWidth(size.width - 30);
  };

  handleScroll = ({ scrollTop, clientHeight, scrollHeight }) => {
    this.bottom = scrollTop + clientHeight >= scrollHeight;
  };

  renderMessage = ({ index, style, key }) => {
    const { messages, select, openPrivateChat } = this.props;

    return (
      <Message
        key={key}
        message={messages.get(index)}
        select={select}
        openPrivateChat={openPrivateChat}
        style={style}
      />
    );
  };

  render() {
    return (
      <div className="messagebox">
        <AutoSizer onResize={this.handleResize}>
          {({ width, height }) => (
            <List
              ref={this.listRef}
              width={width}
              height={height - 14}
              rowCount={this.props.messages.size}
              rowHeight={this.getRowHeight}
              rowRenderer={this.renderMessage}
              onScroll={this.handleScroll}
              style={listStyle}
            />
          )}
        </AutoSizer>
      </div>
    );
  }
}
