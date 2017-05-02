import React, { PureComponent } from 'react';
import { List } from 'react-virtualized/dist/commonjs/List';
import { AutoSizer } from 'react-virtualized/dist/commonjs/AutoSizer';
import debounce from 'lodash/debounce';
import Message from './Message';
import { measureScrollBarWidth } from '../util';

const scrollBarWidth = measureScrollBarWidth();
const listStyle = { padding: '7px 0', boxSizing: 'content-box' };
const threshold = 100;

export default class MessageBox extends PureComponent {
  componentWillUpdate(nextProps) {
    if (nextProps.messages !== this.props.messages) {
      this.list.recomputeRowHeights();
    }

    if (nextProps.tab !== this.props.tab) {
      this.bottom = true;
    }

    if (nextProps.messages.get(0) !== this.props.messages.get(0)) {
      if (nextProps.tab === this.props.tab) {
        const addedMessages = nextProps.messages.size - this.props.messages.size;
        let addedHeight = 0;
        for (let i = 0; i < addedMessages; i++) {
          addedHeight += nextProps.messages.get(i).height;
        }

        this.nextScrollTop = addedHeight + this.container.scrollTop;
      }

      this.loading = false;
    }
  }

  componentDidUpdate() {
    if (this.nextScrollTop > 0) {
      this.container.scrollTop = this.nextScrollTop;
      this.nextScrollTop = 0;
    } else if (this.bottom) {
      this.container.scrollTop = this.container.scrollHeight;
    }

    this.updateWidth();
  }

  getRowHeight = ({ index }) => {
    if (index === 0) {
      if (this.props.hasMoreMessages) {
        return 100;
      }
      return 0;
    }
    return this.props.messages.get(index - 1).height;
  }

  listRef = el => {
    this.list = el;
    // eslint-disable-next-line no-underscore-dangle
    this.container = el.Grid._scrollingContainer;
  };

  updateWidth = (width) => {
    const { tab, setWrapWidth, updateMessageHeight } = this.props;
    let wrapWidth = width || this.width;

    if (width) {
      if (tab.isChannel() && window.innerWidth > 600) {
        wrapWidth += 200;
      }

      this.width = wrapWidth;
    }

    if (this.container.scrollHeight > this.container.clientHeight) {
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

  fetchMore = debounce(() => {
    this.loading = true;
    this.props.onFetchMore();
  }, 100);

  handleScroll = ({ scrollTop, clientHeight, scrollHeight }) => {
    if (this.props.hasMoreMessages &&
      scrollTop <= threshold &&
      scrollTop < this.prevScrollTop &&
      !this.loading) {
      if (this.mouseDown) {
        this.shouldFetch = true;
      } else {
        this.fetchMore();
      }
    }

    this.bottom = scrollTop + clientHeight >= scrollHeight;
    this.prevScrollTop = scrollTop;
  };

  handleMouseDown = () => {
    this.mouseDown = true;
  };

  handleMouseUp = () => {
    this.mouseDown = false;

    if (this.shouldFetch) {
      this.shouldFetch = false;
      this.loading = true;
      this.props.onFetchMore();
    }
  };

  renderMessage = ({ index, style }) => {
    if (index === 0) {
      if (this.props.hasMoreMessages) {
        return (
          <div
            key="top"
            className="messagebox-top-indicator"
            style={style}
          >
            Loading messages...
          </div>
        );
      }
      return null;
    }

    const { messages, onNickClick } = this.props;
    const message = messages.get(index - 1);

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
      <div
        className="messagebox"
        onMouseDown={this.handleMouseDown}
        onMouseUp={this.handleMouseUp}
      >
        <AutoSizer onResize={this.handleResize}>
          {({ width, height }) => (
            <List
              ref={this.listRef}
              width={width}
              height={height - 14}
              rowCount={this.props.messages.size + 1}
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
