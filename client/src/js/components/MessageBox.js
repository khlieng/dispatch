import React, { PureComponent } from 'react';
import { List } from 'react-virtualized/dist/commonjs/List';
import { AutoSizer } from 'react-virtualized/dist/commonjs/AutoSizer';
import Message from './Message';
import { measureScrollBarWidth } from '../util';

const scrollBarWidth = measureScrollBarWidth();
const listStyle = { padding: '7px 0', boxSizing: 'content-box' };

export default class MessageBox extends PureComponent {
  componentWillUpdate(nextProps) {
    if (nextProps.tab !== this.props.tab) {
      this.bottom = true;
    }

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
    const { tab, setWrapWidth, updateMessageHeight } = this.props;
    let wrapWidth = width || this.width;

    if (width) {
      if (tab.isChannel() && window.innerWidth > 600) {
        wrapWidth += 200;
      }

      this.width = wrapWidth;
    }

    // eslint-disable-next-line no-underscore-dangle
    const container = this.list.Grid._scrollingContainer;
    if (container.scrollHeight > container.clientHeight) {
      wrapWidth -= scrollBarWidth;
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

  renderMessage = ({ index, style }) => {
    const { messages, onNickClick } = this.props;
    const message = messages.get(index);

    return (
      <Message
        key={message.id}
        message={message}
        style={style}
        onNickClick={onNickClick}
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
