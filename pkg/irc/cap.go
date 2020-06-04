package irc

import (
	"strings"
)

var clientWantedCaps = []string{"cap-notify"}

func (c *Client) GetCapability(name string) ([]string, bool) {
	c.lock.Lock()
	values, ok := c.enabledCapabilities[name]
	c.lock.Unlock()
	return values, ok
}

func (c *Client) HasCapability(name string, values ...string) bool {
	if capValues, ok := c.GetCapability(name); ok {
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

func (c *Client) beginCAP() {
	c.write("CAP LS 302")
}

func (c *Client) beginSASL() bool {
	if c.negotiating {
		if mechs, ok := c.GetCapability("sasl"); ok {
			if mechs != nil {
				c.filterSASLMechanisms(mechs)
			}
			c.tryNextSASL()
			return true
		}
	}
	return false
}

func (c *Client) finishCAP() {
	if c.negotiating {
		c.negotiating = false
		c.write("CAP END")
	}
}

func (c *Client) handleCAP(msg *Message) {
	if len(msg.Params) < 3 {
		c.write("CAP END")
		return
	}

	caps := parseCaps(msg.LastParam())

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

			c.negotiating = true

			reqCaps := []string{}
			for cap := range c.requestedCapabilities {
				reqCaps = append(reqCaps, cap)
			}

			c.write("CAP REQ :" + strings.Join(reqCaps, " "))
		}

	case "ACK":
		c.lock.Lock()
		for cap := range caps {
			if v, ok := c.requestedCapabilities[cap]; ok {
				c.enabledCapabilities[cap] = v
				delete(c.requestedCapabilities, cap)
			}
		}
		c.lock.Unlock()

		if len(c.requestedCapabilities) == 0 && !c.beginSASL() {
			c.finishCAP()
		}

	case "NAK":
		for cap := range caps {
			delete(c.requestedCapabilities, cap)
		}

		if len(c.requestedCapabilities) == 0 && !c.beginSASL() {
			c.finishCAP()
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
		c.lock.Lock()
		for cap := range caps {
			delete(c.enabledCapabilities, cap)
		}
		c.lock.Unlock()
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
