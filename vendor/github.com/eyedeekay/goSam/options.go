package goSam

import (
	"fmt"
	"strconv"
	"strings"
)

//Option is a client Option
type Option func(*Client) error

//SetAddr sets a clients's address in the form host:port or host, port
func SetAddr(s ...string) func(*Client) error {
	return func(c *Client) error {
		if len(s) == 1 {
			split := strings.SplitN(s[0], ":", 2)
			if len(split) == 2 {
				if i, err := strconv.Atoi(split[1]); err == nil {
					if i < 65536 {
						c.host = split[0]
						c.port = split[1]
						return nil
					}
					return fmt.Errorf("Invalid port")
				}
				return fmt.Errorf("Invalid port; non-number")
			}
			return fmt.Errorf("Invalid address; use host:port %s", split)
		} else if len(s) == 2 {
			if i, err := strconv.Atoi(s[1]); err == nil {
				if i < 65536 {
					c.host = s[0]
					c.port = s[1]
					return nil
				}
				return fmt.Errorf("Invalid port")
			}
			return fmt.Errorf("Invalid port; non-number")
		} else {
			return fmt.Errorf("Invalid address")
		}
	}
}

//SetAddrMixed sets a clients's address in the form host, port(int)
func SetAddrMixed(s string, i int) func(*Client) error {
	return func(c *Client) error {
		if i < 65536 && i > 0 {
			c.host = s
			c.port = strconv.Itoa(i)
			return nil
		}
		return fmt.Errorf("Invalid port")
	}
}

//SetHost sets the host of the client's SAM bridge
func SetHost(s string) func(*Client) error {
	return func(c *Client) error {
		c.host = s
		return nil
	}
}

//SetLocalDestination sets the local destination of the tunnel from a private
//key
func SetLocalDestination(s string) func(*Client) error {
	return func(c *Client) error {
		c.destination = s
		return nil
	}
}

func setlastaddr(s string) func(*Client) error {
	return func(c *Client) error {
		c.lastaddr = s
		return nil
	}
}

func setid(s int32) func(*Client) error {
	return func(c *Client) error {
		c.id = s
		return nil
	}
}

//SetPort sets the port of the client's SAM bridge using a string
func SetPort(s string) func(*Client) error {
	return func(c *Client) error {
		port, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("Invalid port; non-number")
		}
		if port < 65536 && port > -1 {
			c.port = s
			return nil
		}
		return fmt.Errorf("Invalid port")
	}
}

//SetPortInt sets the port of the client's SAM bridge using a string
func SetPortInt(i int) func(*Client) error {
	return func(c *Client) error {
		if i < 65536 && i > -1 {
			c.port = strconv.Itoa(i)
			return nil
		}
		return fmt.Errorf("Invalid port")
	}
}

//SetFromPort sets the port of the client's SAM bridge using a string
func SetFromPort(s string) func(*Client) error {
	return func(c *Client) error {
		port, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("Invalid port; non-number")
		}
		if port < 65536 && port > -1 {
			c.fromport = s
			return nil
		}
		return fmt.Errorf("Invalid port")
	}
}

//SetFromPortInt sets the port of the client's SAM bridge using a string
func SetFromPortInt(i int) func(*Client) error {
	return func(c *Client) error {
		if i < 65536 && i > -1 {
			c.fromport = strconv.Itoa(i)
			return nil
		}
		return fmt.Errorf("Invalid port")
	}
}

//SetToPort sets the port of the client's SAM bridge using a string
func SetToPort(s string) func(*Client) error {
	return func(c *Client) error {
		port, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("Invalid port; non-number")
		}
		if port < 65536 && port > -1 {
			c.toport = s
			return nil
		}
		return fmt.Errorf("Invalid port")
	}
}

//SetToPortInt sets the port of the client's SAM bridge using a string
func SetToPortInt(i int) func(*Client) error {
	return func(c *Client) error {
		if i < 65536 && i > -1 {
			c.fromport = strconv.Itoa(i)
			return nil
		}
		return fmt.Errorf("Invalid port")
	}
}

//SetDebug enables debugging messages
func SetDebug(b bool) func(*Client) error {
	return func(c *Client) error {
		c.debug = b
		return nil
	}
}

//SetInLength sets the number of hops inbound
func SetInLength(u uint) func(*Client) error {
	return func(c *Client) error {
		if u < 7 {
			c.inLength = u
			return nil
		}
		return fmt.Errorf("Invalid inbound tunnel length")
	}
}

//SetOutLength sets the number of hops outbound
func SetOutLength(u uint) func(*Client) error {
	return func(c *Client) error {
		if u < 7 {
			c.outLength = u
			return nil
		}
		return fmt.Errorf("Invalid outbound tunnel length")
	}
}

//SetInVariance sets the variance of a number of hops inbound
func SetInVariance(i int) func(*Client) error {
	return func(c *Client) error {
		if i < 7 && i > -7 {
			c.inVariance = i
			return nil
		}
		return fmt.Errorf("Invalid inbound tunnel length")
	}
}

//SetOutVariance sets the variance of a number of hops outbound
func SetOutVariance(i int) func(*Client) error {
	return func(c *Client) error {
		if i < 7 && i > -7 {
			c.outVariance = i
			return nil
		}
		return fmt.Errorf("Invalid outbound tunnel variance")
	}
}

//SetInQuantity sets the inbound tunnel quantity
func SetInQuantity(u uint) func(*Client) error {
	return func(c *Client) error {
		if u <= 16 {
			c.inQuantity = u
			return nil
		}
		return fmt.Errorf("Invalid inbound tunnel quantity")
	}
}

//SetOutQuantity sets the outbound tunnel quantity
func SetOutQuantity(u uint) func(*Client) error {
	return func(c *Client) error {
		if u <= 16 {
			c.outQuantity = u
			return nil
		}
		return fmt.Errorf("Invalid outbound tunnel quantity")
	}
}

//SetInBackups sets the inbound tunnel backups
func SetInBackups(u uint) func(*Client) error {
	return func(c *Client) error {
		if u < 6 {
			c.inBackups = u
			return nil
		}
		return fmt.Errorf("Invalid inbound tunnel backup quantity")
	}
}

//SetOutBackups sets the inbound tunnel backups
func SetOutBackups(u uint) func(*Client) error {
	return func(c *Client) error {
		if u < 6 {
			c.outBackups = u
			return nil
		}
		return fmt.Errorf("Invalid outbound tunnel backup quantity")
	}
}

//SetUnpublished tells the router to not publish the client leaseset
func SetUnpublished(b bool) func(*Client) error {
	return func(c *Client) error {
		c.dontPublishLease = b
		return nil
	}
}

//SetEncrypt tells the router to use an encrypted leaseset
func SetEncrypt(b bool) func(*Client) error {
	return func(c *Client) error {
		c.encryptLease = b
		return nil
	}
}

//SetLeaseSetEncType tells the router to use an encrypted leaseset of a specific type.
//defaults to 4,0
func SetLeaseSetEncType(b string) func(*Client) error {
	return func(c *Client) error {
		c.leaseSetEncType = b
		return nil
	}
}

//SetReduceIdle sets the created tunnels to be reduced during extended idle time to avoid excessive resource usage
func SetReduceIdle(b bool) func(*Client) error {
	return func(c *Client) error {
		c.reduceIdle = b
		return nil
	}
}

//SetReduceIdleTime sets time to wait before the tunnel quantity is reduced
func SetReduceIdleTime(u uint) func(*Client) error {
	return func(c *Client) error {
		if u > 299999 {
			c.reduceIdleTime = u
			return nil
		}
		return fmt.Errorf("Invalid reduce idle time %v", u)
	}
}

//SetReduceIdleQuantity sets number of tunnels to keep alive during an extended idle period
func SetReduceIdleQuantity(u uint) func(*Client) error {
	return func(c *Client) error {
		if u < 5 {
			c.reduceIdleQuantity = u
			return nil
		}
		return fmt.Errorf("Invalid reduced tunnel quantity %v", u)
	}
}

//SetCloseIdle sets the tunnels to close after a specific amount of time
func SetCloseIdle(b bool) func(*Client) error {
	return func(c *Client) error {
		c.closeIdle = b
		return nil
	}
}

//SetCloseIdleTime sets the time in milliseconds to wait before closing tunnels
func SetCloseIdleTime(u uint) func(*Client) error {
	return func(c *Client) error {
		if u > 299999 {
			c.closeIdleTime = u
			return nil
		}
		return fmt.Errorf("Invalid close idle time %v", u)
	}
}

//SetCompression sets the tunnels to close after a specific amount of time
func SetCompression(b bool) func(*Client) error {
	return func(c *Client) error {
		c.compression = b
		return nil
	}
}

/* SAM v 3.1 Options*/

//SetSignatureType tells gosam to pass SAM a signature_type parameter with one
// of the following values:
//    "SIGNATURE_TYPE=DSA_SHA1",
//    "SIGNATURE_TYPE=ECDSA_SHA256_P256",
//    "SIGNATURE_TYPE=ECDSA_SHA384_P384",
//    "SIGNATURE_TYPE=ECDSA_SHA512_P521",
//    "SIGNATURE_TYPE=EdDSA_SHA512_Ed25519",
// or an empty string
func SetSignatureType(s string) func(*Client) error {
	return func(c *Client) error {
		if s == "" {
			c.sigType = ""
			return nil
		}
		for _, valid := range SAMsigTypes {
			if s == valid {
				c.sigType = valid
				return nil
			}
		}
		return fmt.Errorf("Invalid signature type specified at construction time")
	}
}

//return the from port as a string.
func (c *Client) from() string {
	if c.fromport == "FROM_PORT=0" {
		return ""
	}
	if c.fromport == "0" {
		return ""
	}
	if c.fromport == "" {
		return ""
	}
	return fmt.Sprintf(" FROM_PORT=%v ", c.fromport)
}

//return the to port as a string.
func (c *Client) to() string {
	if c.fromport == "TO_PORT=0" {
		return ""
	}
	if c.fromport == "0" {
		return ""
	}
	if c.toport == "" {
		return ""
	}
	return fmt.Sprintf(" TO_PORT=%v ", c.toport)
}

//return the signature type as a string.
func (c *Client) sigtype() string {
	return fmt.Sprintf(" %s ", c.sigType)
}

//return the inbound length as a string.
func (c *Client) inlength() string {
	return fmt.Sprintf(" inbound.length=%d ", c.inLength)
}

//return the outbound length as a string.
func (c *Client) outlength() string {
	return fmt.Sprintf(" outbound.length=%d ", c.outLength)
}

//return the inbound length variance as a string.
func (c *Client) invariance() string {
	return fmt.Sprintf(" inbound.lengthVariance=%d ", c.inVariance)
}

//return the outbound length variance as a string.
func (c *Client) outvariance() string {
	return fmt.Sprintf(" outbound.lengthVariance=%d ", c.outVariance)
}

//return the inbound tunnel quantity as a string.
func (c *Client) inquantity() string {
	return fmt.Sprintf(" inbound.quantity=%d ", c.inQuantity)
}

//return the outbound tunnel quantity as a string.
func (c *Client) outquantity() string {
	return fmt.Sprintf(" outbound.quantity=%d ", c.outQuantity)
}

//return the inbound tunnel quantity as a string.
func (c *Client) inbackups() string {
	return fmt.Sprintf(" inbound.backupQuantity=%d ", c.inQuantity)
}

//return the outbound tunnel quantity as a string.
func (c *Client) outbackups() string {
	return fmt.Sprintf(" outbound.backupQuantity=%d ", c.outQuantity)
}

func (c *Client) encryptlease() string {
	if c.encryptLease {
		return " i2cp.encryptLeaseSet=true "
	}
	return " i2cp.encryptLeaseSet=false "
}

func (c *Client) leasesetenctype() string {
	if c.encryptLease {
		return fmt.Sprintf(" i2cp.leaseSetEncType=%s ", c.leaseSetEncType)
	}
	return " i2cp.leaseSetEncType=4,0 "
}

func (c *Client) dontpublishlease() string {
	if c.dontPublishLease {
		return " i2cp.dontPublishLeaseSet=true "
	}
	return " i2cp.dontPublishLeaseSet=false "
}

func (c *Client) closeonidle() string {
	if c.closeIdle {
		return " i2cp.closeOnIdle=true "
	}
	return " i2cp.closeOnIdle=false "
}

func (c *Client) closeidletime() string {
	return fmt.Sprintf(" i2cp.closeIdleTime=%d ", c.closeIdleTime)
}

func (c *Client) reduceonidle() string {
	if c.reduceIdle {
		return " i2cp.reduceOnIdle=true "
	}
	return " i2cp.reduceOnIdle=false "
}

func (c *Client) reduceidletime() string {
	return fmt.Sprintf(" i2cp.reduceIdleTime=%d ", c.reduceIdleTime)
}

func (c *Client) reduceidlecount() string {
	return fmt.Sprintf(" i2cp.reduceIdleQuantity=%d ", c.reduceIdleQuantity)
}

func (c *Client) compresion() string {
	if c.compression {
		return " i2cp.gzip=true "
	}
	return " i2cp.gzip=false "
}

//return all options as string ready for passing to sendcmd
func (c *Client) allOptions() string {
	return c.inlength() +
		c.outlength() +
		c.invariance() +
		c.outvariance() +
		c.inquantity() +
		c.outquantity() +
		c.inbackups() +
		c.outbackups() +
		c.dontpublishlease() +
		c.encryptlease() +
		c.leasesetenctype() +
		c.reduceonidle() +
		c.reduceidletime() +
		c.reduceidlecount() +
		c.closeonidle() +
		c.closeidletime() +
		c.compresion()
}

//Print return all options as string
func (c *Client) Print() string {
	return c.inlength() +
		c.outlength() +
		c.invariance() +
		c.outvariance() +
		c.inquantity() +
		c.outquantity() +
		c.inbackups() +
		c.outbackups() +
		c.dontpublishlease() +
		c.encryptlease() +
		c.leasesetenctype() +
		c.reduceonidle() +
		c.reduceidletime() +
		c.reduceidlecount() +
		c.closeonidle() +
		c.closeidletime() +
		c.compresion()
}
