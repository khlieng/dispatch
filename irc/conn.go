package irc

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"time"
)

func (c *Client) Connect(address string) {
	c.ConnectionChanged <- false

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

func (c *Client) run() {
	c.tryConnect()

	for {
		select {
		case <-c.quit:
			if c.Connected() {
				c.disconnect()
			}

			c.sendRecv.Wait()
			close(c.Messages)
			return

		case <-c.reconnect:
			c.disconnect()

			c.sendRecv.Wait()
			c.reconnect = make(chan struct{})
			c.once.Reset()

			c.tryConnect()
		}
	}
}

func (c *Client) disconnect() {
	c.ConnectionChanged <- false
	c.lock.Lock()
	c.connected = false
	c.lock.Unlock()

	c.once.Do(c.ready.Done)
	c.conn.Close()
}

func (c *Client) tryConnect() {
	for {
		select {
		case <-c.quit:
			return

		default:
		}

		err := c.connect()
		if err == nil {
			c.backoff.Reset()

			c.flushChannels()
			return
		}

		time.Sleep(c.backoff.Duration())
	}
}

func (c *Client) connect() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.TLS {
		conn, err := tls.DialWithDialer(c.dialer, "tcp", c.Server, c.TLSConfig)
		if err != nil {
			return err
		}

		c.conn = conn
	} else {
		conn, err := c.dialer.Dial("tcp", c.Server)
		if err != nil {
			return err
		}

		c.conn = conn
	}

	c.connected = true
	c.ConnectionChanged <- true
	c.reader = bufio.NewReader(c.conn)

	c.register()

	c.ready.Add(1)
	c.sendRecv.Add(2)
	go c.send()
	go c.recv()

	return nil
}

func (c *Client) send() {
	defer c.sendRecv.Done()

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
	defer c.sendRecv.Done()

	for {
		line, err := c.reader.ReadString('\n')
		if err != nil {
			select {
			case <-c.quit:
				return

			default:
				close(c.reconnect)
				return
			}
		}

		msg := parseMessage(line)

		switch msg.Command {
		case Ping:
			go c.write("PONG :" + msg.LastParam())

		case Join:
			if msg.Nick == c.GetNick() {
				c.addChannel(msg.Params[0])
			}

		case Nick:
			if msg.Nick == c.GetNick() {
				c.setNick(msg.LastParam())
			}

		case ReplyWelcome:
			c.once.Do(c.ready.Done)

		case ErrNicknameInUse:
			if c.HandleNickInUse != nil {
				newNick := c.HandleNickInUse(msg.Params[1])
				// Set the nick here aswell incase this happens during registration
				// since there will be no NICK message to confirm it then
				c.setNick(newNick)
				go c.writeNick(newNick)
			}
		}

		c.Messages <- msg
	}
}
