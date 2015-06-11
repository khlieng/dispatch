package irc

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

func (c *Client) Connect(address string) {
	if idx := strings.Index(address, ":"); idx < 0 {
		c.Host = address

		if c.TLS {
			address += ":6697"
		} else {
			address += ":6667"
		}
	} else {
		c.Host = address[:idx]
	}
	c.Server = address
	c.dialer = &net.Dialer{Timeout: 10 * time.Second}

	go c.run()
}

func (c *Client) Write(data string) {
	c.out <- data + "\r\n"
}

func (c *Client) Writef(format string, a ...interface{}) {
	c.out <- fmt.Sprintf(format+"\r\n", a...)
}

func (c *Client) write(data string) {
	c.conn.Write([]byte(data + "\r\n"))
}

func (c *Client) writef(format string, a ...interface{}) {
	fmt.Fprintf(c.conn, format+"\r\n", a...)
}

func (c *Client) connect() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.TLS {
		if c.TLSConfig == nil {
			c.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		}

		if conn, err := tls.DialWithDialer(c.dialer, "tcp", c.Server, c.TLSConfig); err != nil {
			return err
		} else {
			c.conn = conn
		}
	} else {
		if conn, err := c.dialer.Dial("tcp", c.Server); err != nil {
			return err
		} else {
			c.conn = conn
		}
	}

	c.connected = true
	c.reader = bufio.NewReader(c.conn)

	c.register()

	c.ready.Add(1)
	go c.send()
	go c.recv()

	return nil
}

func (c *Client) tryConnect() {
	// TODO: backoff
	for {
		select {
		case <-c.quit:
			return
		default:
		}

		err := c.connect()
		if err == nil {
			return
		}
	}
}

func (c *Client) run() {
	c.tryConnect()
	for {
		select {
		case <-c.quit:
			c.close()
			c.lock.Lock()
			c.connected = false
			c.lock.Unlock()
			return

		case <-c.reconnect:
			c.reconnect = make(chan struct{})
			c.once = sync.Once{}
			c.tryConnect()
		}
	}
}

func (c *Client) send() {
	c.ready.Wait()
	for {
		select {
		case <-c.quit:
			return

		case <-c.reconnect:
			return

		case msg := <-c.out:
			_, err := c.conn.Write([]byte(msg))
			if err != nil {
				return
			}
		}
	}
}

func (c *Client) recv() {
	defer c.conn.Close()
	for {
		line, err := c.reader.ReadString('\n')
		if err != nil {
			c.lock.Lock()
			c.connected = false
			c.lock.Unlock()
			c.once.Do(c.ready.Done)
			select {
			case <-c.quit:
				return
			default:
				close(c.reconnect)
				return
			}
		}

		msg := parseMessage(line)
		c.Messages <- msg

		switch msg.Command {
		case Ping:
			go c.write("PONG :" + msg.Trailing)

		case ReplyWelcome:
			c.once.Do(c.ready.Done)
		}
	}
}

func (c *Client) close() {
	if c.Connected() {
		c.once.Do(c.ready.Done)
	}
	close(c.out)
	close(c.Messages)
}
