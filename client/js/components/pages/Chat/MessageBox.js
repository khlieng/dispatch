import React, { PureComponent, createRef } from 'react';
import { VariableSizeList as List } from 'react-window';
import AutoSizer from 'react-virtualized-auto-sizer';
import debounce from 'lodash/debounce';
import { formatDate, measureScrollBarWidth } from 'utils';
import { getScrollPos, saveScrollPos } from 'utils/scrollPosition';
import { windowHeight } from 'utils/size';
import Message from './Message';

const fetchThreshold = 600;
// The amount of time in ms that needs to pass without any
// scroll events happening before adding messages to the top,
// this is done to prevent the scroll from jumping all over the place
const scrollbackDebounce = 150;

const scrollBarWidth = `${measureScrollBarWidth()}px`;

const hasSameLastMessage = (m1, m2) => {
  if (m1.length === 0 || m2.length === 0) {
    if (m1.length === 0 && m2.length === 0) {
      return true;
    }
    return false;
  }

  return m1[m1.length - 1].id === m2[m2.length - 1].id;
};

export default class MessageBox extends PureComponent {
  state = { topDate: '' };

  list = createRef();

  outer = createRef();

  addMore = debounce(() => {
    const { tab, onAddMore } = this.props;
    this.ready = true;
    onAddMore(tab.server, tab.name);
  }, scrollbackDebounce);

  constructor(props) {
    super(props);

    this.loadScrollPos();
  }

  componentDidUpdate(prevProps) {
    const { messages } = this.props;

    if (prevProps.tab !== this.props.tab) {
      this.loadScrollPos(true);
    }

    if (this.nextScrollTop > 0) {
      this.list.current.scrollTo(this.nextScrollTop);
      this.nextScrollTop = 0;
    } else if (
      this.bottom &&
      !hasSameLastMessage(messages, prevProps.messages)
    ) {
      this.list.current.scrollToItem(messages.length + 1);
    }
  }

  componentWillUnmount() {
    this.saveScrollPos();
  }

  getSnapshotBeforeUpdate(prevProps) {
    if (prevProps.messages !== this.props.messages) {
      this.list.current.resetAfterIndex(0);
    }

    if (prevProps.tab !== this.props.tab) {
      this.saveScrollPos();
      this.bottom = false;
    }

    if (prevProps.messages[0] !== this.props.messages[0]) {
      const { messages, hasMoreMessages } = this.props;

      if (prevProps.tab === this.props.tab && prevProps.messages.length > 0) {
        const addedMessages = messages.length - prevProps.messages.length;
        let addedHeight = 0;
        for (let i = 0; i < addedMessages; i++) {
          addedHeight += messages[i].height;
        }

        this.nextScrollTop = addedHeight + this.outer.current.scrollTop;

        if (!hasMoreMessages) {
          this.nextScrollTop -= 93;
        }
      }

      this.loading = false;
      this.ready = false;
    }

    return null;
  }

  getRowHeight = index => {
    const { messages, hasMoreMessages } = this.props;

    if (index === 0) {
      if (hasMoreMessages) {
        return 100;
      }
      return 7;
    }
    if (index === messages.length + 1) {
      return 7;
    }
    return messages[index - 1].height;
  };

  getItemKey = index => {
    const { messages } = this.props;

    if (index === 0) {
      return 'top';
    }
    if (index === messages.length + 1) {
      return 'bottom';
    }
    return messages[index - 1].id;
  };

  updateScrollKey = () => {
    const { tab } = this.props;
    this.scrollKey = `msg:${tab.server}:${tab.name}`;
    return this.scrollKey;
  };

  loadScrollPos = scroll => {
    const pos = getScrollPos(this.updateScrollKey());
    const { messages } = this.props;

    if (pos >= 0) {
      this.bottom = false;
      if (scroll) {
        this.list.current.scrollTo(pos);
      } else {
        this.initialScrollTop = pos;
      }
    } else {
      this.bottom = true;
      if (scroll) {
        this.list.current.scrollToItem(messages.length + 1);
      } else if (messages.length > 0) {
        let totalHeight = 14;
        for (let i = 0; i < messages.length; i++) {
          totalHeight += messages[i].height;
        }

        const messageBoxHeight = windowHeight() - 100;
        if (totalHeight > messageBoxHeight) {
          this.initialScrollTop = totalHeight;
        }
      }
    }
  };

  saveScrollPos = () => {
    if (this.bottom) {
      saveScrollPos(this.scrollKey, -1);
    } else {
      saveScrollPos(this.scrollKey, this.scrollTop);
    }
  };

  fetchMore = () => {
    this.loading = true;
    this.props.onFetchMore();
  };

  handleScroll = ({ scrollOffset, scrollDirection }) => {
    if (
      !this.loading &&
      this.props.hasMoreMessages &&
      scrollOffset <= fetchThreshold &&
      scrollDirection === 'backward'
    ) {
      this.fetchMore();
    }

    if (this.loading && !this.ready) {
      if (this.mouseDown) {
        this.ready = true;
        this.shouldAdd = true;
      } else {
        this.addMore();
      }
    }

    const { clientHeight, scrollHeight } = this.outer.current;

    this.scrollTop = scrollOffset;
    this.bottom = scrollOffset + clientHeight >= scrollHeight - 20;
  };

  handleItemsRendered = ({ visibleStartIndex }) => {
    const startIndex = visibleStartIndex === 0 ? 0 : visibleStartIndex - 1;
    const firstVisibleMessage = this.props.messages[startIndex];

    if (firstVisibleMessage && firstVisibleMessage.date) {
      this.setState({ topDate: formatDate(firstVisibleMessage.date) });
    } else {
      this.setState({ topDate: '' });
    }
  };

  handleMouseDown = () => {
    this.mouseDown = true;
  };

  handleMouseUp = () => {
    this.mouseDown = false;

    if (this.shouldAdd) {
      const { tab, onAddMore } = this.props;
      this.shouldAdd = false;
      onAddMore(tab.server, tab.name);
    }
  };

  renderMessage = ({ index, style }) => {
    const { messages } = this.props;

    if (index === 0) {
      if (this.props.hasMoreMessages) {
        return (
          <div className="messagebox-top-indicator" style={style}>
            Loading messages...
          </div>
        );
      }
      return null;
    }
    if (index === messages.length + 1) {
      return null;
    }

    const { coloredNicks, onNickClick } = this.props;
    const message = messages[index - 1];

    return (
      <div style={style}>
        <Message
          message={message}
          coloredNick={coloredNicks}
          onNickClick={onNickClick}
        />
      </div>
    );
  };

  render() {
    const { messages, hideTopDate } = this.props;
    const { topDate } = this.state;

    const dateContainerStyle = {
      right: scrollBarWidth
    };

    return (
      <div
        className="messagebox"
        onMouseDown={this.handleMouseDown}
        onMouseUp={this.handleMouseUp}
      >
        <div
          className="messagebox-topdate-container"
          style={dateContainerStyle}
        >
          {!hideTopDate && topDate && (
            <span className="messagebox-topdate">{topDate}</span>
          )}
        </div>
        <AutoSizer>
          {({ width, height }) => (
            <List
              ref={this.list}
              outerRef={this.outer}
              width={width}
              height={height}
              itemCount={messages.length + 2}
              itemKey={this.getItemKey}
              itemSize={this.getRowHeight}
              estimatedItemSize={32}
              initialScrollOffset={this.initialScrollTop}
              onScroll={this.handleScroll}
              onItemsRendered={this.handleItemsRendered}
              className="messagebox-window"
              overscanCount={5}
            >
              {this.renderMessage}
            </List>
          )}
        </AutoSizer>
      </div>
    );
  }
}
