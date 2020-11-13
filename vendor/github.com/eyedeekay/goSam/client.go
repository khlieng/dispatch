package goSam

import (
	"bufio"
	"crypto/sha256"
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"net"
	"strings"
	"sync"
)

// A Client represents a single Connection to the SAM bridge
type Client struct {
	host     string
	port     string
	fromport string
	toport   string

	SamConn net.Conn
	rd      *bufio.Reader

	sigType     string
	destination string

	inLength   uint
	inVariance int
	inQuantity uint
	inBackups  uint

	outLength   uint
	outVariance int
	outQuantity uint
	outBackups  uint

	dontPublishLease bool
	encryptLease     bool
	leaseSetEncType  string

	reduceIdle         bool
	reduceIdleTime     uint
	reduceIdleQuantity uint

	closeIdle     bool
	closeIdleTime uint

	compression bool

	debug bool
	//NEVER, EVER modify lastaddr or id yourself. They are used internally only.
	lastaddr string
	id       int32
	ml       sync.Mutex
}

var SAMsigTypes = []string{
	"SIGNATURE_TYPE=DSA_SHA1",
	"SIGNATURE_TYPE=ECDSA_SHA256_P256",
	"SIGNATURE_TYPE=ECDSA_SHA384_P384",
	"SIGNATURE_TYPE=ECDSA_SHA512_P521",
	"SIGNATURE_TYPE=EdDSA_SHA512_Ed25519",
}

var (
	i2pB64enc *base64.Encoding = base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-~")
	i2pB32enc *base32.Encoding = base32.NewEncoding("abcdefghijklmnopqrstuvwxyz234567")
)

// NewDefaultClient creates a new client, connecting to the default host:port at localhost:7656
func NewDefaultClient() (*Client, error) {
	return NewClient("localhost:7656")
}

// NewClient creates a new client, connecting to a specified port
func NewClient(addr string) (*Client, error) {
	return NewClientFromOptions(SetAddr(addr))
}

// NewID generates a random number to use as an tunnel name
func (c *Client) NewID() int32 {
	return rand.Int31n(math.MaxInt32)
}

// Destination returns the full destination of the local tunnel
func (c *Client) Destination() string {
	return c.destination
}

// Base32 returns the base32 of the local tunnel
func (c *Client) Base32() string {
	//	hash := sha256.New()
	b64, err := i2pB64enc.DecodeString(c.Base64())
	if err != nil {
		return ""
	}
	//hash.Write([]byte(b64))
	var s []byte
	for _, e := range sha256.Sum256(b64) {
		s = append(s, e)
	}
	return strings.ToLower(strings.Replace(i2pB32enc.EncodeToString(s), "=", "", -1))
}

func (c *Client) base64() []byte {
	if c.destination != "" {
		s, _ := i2pB64enc.DecodeString(c.destination)
		alen := binary.BigEndian.Uint16(s[385:387])
		return s[:387+alen]
	}
	return []byte("")
}

// Base64 returns the base64 of the local tunnel
func (c *Client) Base64() string {
	return i2pB64enc.EncodeToString(c.base64())
}

// NewClientFromOptions creates a new client, connecting to a specified port
func NewClientFromOptions(opts ...func(*Client) error) (*Client, error) {
	var c Client
	c.host = "127.0.0.1"
	c.port = "7656"
	c.inLength = 3
	c.inVariance = 0
	c.inQuantity = 1
	c.inBackups = 1
	c.outLength = 3
	c.outVariance = 0
	c.outQuantity = 1
	c.outBackups = 1
	c.dontPublishLease = true
	c.encryptLease = false
	c.reduceIdle = false
	c.reduceIdleTime = 300000
	c.reduceIdleQuantity = 1
	c.closeIdle = true
	c.closeIdleTime = 600000
	c.debug = true
	c.sigType = SAMsigTypes[4]
	c.id = 0
	c.lastaddr = "invalid"
	c.destination = ""
	c.leaseSetEncType = "4,0"
	c.fromport = ""
	c.toport = ""
	for _, o := range opts {
		if err := o(&c); err != nil {
			return nil, err
		}
	}
	conn, err := net.Dial("tcp", c.samaddr())
	if err != nil {
		return nil, err
	}
	if c.debug {
		conn = WrapConn(conn)
	}
	c.SamConn = conn
	c.rd = bufio.NewReader(conn)
	return &c, c.hello()
}

func (p *Client) ID() string {
	return fmt.Sprintf("%d", p.id)
}

func (p *Client) Addr() net.Addr {
	return nil
}

//return the combined host:port of the SAM bridge
func (c *Client) samaddr() string {
	return fmt.Sprintf("%s:%s", c.host, c.port)
}

// send the initial handshake command and check that the reply is ok
func (c *Client) hello() error {
	r, err := c.sendCmd("HELLO VERSION MIN=3.0 MAX=3.2\n")
	if err != nil {
		return err
	}

	if r.Topic != "HELLO" {
		return fmt.Errorf("Client Hello Unknown Reply: %+v\n", r)
	}

	if r.Pairs["RESULT"] != "OK" {
		return fmt.Errorf("Handshake did not succeed\nReply:%+v\n", r)
	}

	return nil
}

// helper to send one command and parse the reply by sam
func (c *Client) sendCmd(str string, args ...interface{}) (*Reply, error) {
	if _, err := fmt.Fprintf(c.SamConn, str, args...); err != nil {
		return nil, err
	}

	line, err := c.rd.ReadString('\n')
	if err != nil {
		return nil, err
	}

	return parseReply(line)
}

// Close the underlying socket to SAM
func (c *Client) Close() error {
	c.rd = nil
	return c.SamConn.Close()
}

// NewClient generates an exact copy of the client with the same options
func (c *Client) NewClient() (*Client, error) {
	return NewClientFromOptions(
		SetHost(c.host),
		SetPort(c.port),
		SetDebug(c.debug),
		SetInLength(c.inLength),
		SetOutLength(c.outLength),
		SetInVariance(c.inVariance),
		SetOutVariance(c.outVariance),
		SetInQuantity(c.inQuantity),
		SetOutQuantity(c.outQuantity),
		SetInBackups(c.inBackups),
		SetOutBackups(c.outBackups),
		SetUnpublished(c.dontPublishLease),
		SetEncrypt(c.encryptLease),
		SetReduceIdle(c.reduceIdle),
		SetReduceIdleTime(c.reduceIdleTime),
		SetReduceIdleQuantity(c.reduceIdleQuantity),
		SetCloseIdle(c.closeIdle),
		SetCloseIdleTime(c.closeIdleTime),
		SetCompression(c.compression),
		setlastaddr(c.lastaddr),
		setid(c.id),
	)
}
