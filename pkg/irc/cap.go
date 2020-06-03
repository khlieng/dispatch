package irc

import (
	"strings"
)

func (c *Client) HasCapability(name string, values ...string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	if capValues, ok := c.enabledCapabilities[name]; ok {
		if len(values) == 0 || capValues == nil {
			return true
		}

		for _, v := range values {
			for _, vCap := range capValues {
				if v == vCap {
					return true
				}
			}
		}
	}

	return false
}

var clientWantedCaps = []string{}

func (c *Client) writeCAP() {
	c.write("CAP LS 302")
}

func (c *Client) handleCAP(msg *Message) {
	if len(msg.Params) < 3 {
		c.write("CAP END")
		return
	}

	caps := parseCaps(msg.LastParam())

	c.lock.Lock()
	defer c.lock.Unlock()

	switch msg.Params[1] {
	case "LS":
		for cap, values := range caps {
			for _, wanted := range c.wantedCapabilities {
				if cap == wanted {
					c.requestedCapabilities[cap] = values
				}
			}
		}

		if len(msg.Params) == 3 {
			if len(c.requestedCapabilities) == 0 {
				c.write("CAP END")
				return
			}

			reqCaps := []string{}
			for cap := range c.requestedCapabilities {
				reqCaps = append(reqCaps, cap)
			}

			c.write("CAP REQ :" + strings.Join(reqCaps, " "))
		}

	case "ACK":
		for cap := range caps {
			if v, ok := c.requestedCapabilities[cap]; ok {
				c.enabledCapabilities[cap] = v
				delete(c.requestedCapabilities, cap)
			}
		}

		if len(c.requestedCapabilities) == 0 {
			if c.Config.SASL != nil && c.HasCapability("sasl", c.Config.SASL.Name()) {
				c.write("AUTHENTICATE " + c.Config.SASL.Name())
			} else {
				c.write("CAP END")
			}
		}

	case "NAK":
		for cap := range caps {
			delete(c.requestedCapabilities, cap)
		}

		if len(c.requestedCapabilities) == 0 {
			c.write("CAP END")
		}

	case "NEW":
		reqCaps := []string{}
		for cap, values := range caps {
			for _, wanted := range c.wantedCapabilities {
				if cap == wanted && !c.HasCapability(cap) {
					c.requestedCapabilities[cap] = values
					reqCaps = append(reqCaps, cap)
				}
			}
		}

		if len(reqCaps) > 0 {
			c.write("CAP REQ :" + strings.Join(reqCaps, " "))
		}

	case "DEL":
		for cap := range caps {
			delete(c.enabledCapabilities, cap)
		}
	}
}

func parseCaps(caps string) map[string][]string {
	result := map[string][]string{}

	parts := strings.Split(caps, " ")
	for _, part := range parts {
		capParts := strings.Split(part, "=")
		name := capParts[0]

		if len(capParts) > 1 {
			result[name] = strings.Split(capParts[1], ",")
		} else {
			result[name] = nil
		}
	}

	return result
}
