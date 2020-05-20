package irc

import (
	"fmt"
	"strings"
	"time"
)

// ClientInfo is the CTCP messages this client implements
const ClientInfo = "ACTION CLIENTINFO DCC FINGER PING SOURCE TIME VERSION USERINFO"

type CTCP struct {
	Command string
	Params  string
}

func DecodeCTCP(str string) *CTCP {
	if len(str) > 1 && str[0] == 0x01 {
		parts := strings.SplitN(strings.Trim(str, "\x01"), " ", 2)
		ctcp := CTCP{}

		if parts[0] != "" {
			ctcp.Command = parts[0]
		} else {
			return nil
		}

		if len(parts) == 2 {
			ctcp.Params = parts[1]
		}

		return &ctcp
	}

	return nil
}

func EncodeCTCP(ctcp *CTCP) string {
	if ctcp == nil || ctcp.Command == "" {
		return ""
	}
	return fmt.Sprintf("\x01%s %s\x01", ctcp.Command, ctcp.Params)
}

func (c *Client) handleCTCP(ctcp *CTCP, msg *Message) {
	switch ctcp.Command {
	case "CLIENTINFO":
		c.ReplyCTCP(msg.Nick, ctcp.Command, ClientInfo)

	case "DCC":
		if strings.HasPrefix(ctcp.Params, "SEND") {
			if dccSend := ParseDCCSend(ctcp); dccSend != nil {
				go c.Download(dccSend)
			}
		}

	case "FINGER", "VERSION":
		if c.Version != "" {
			c.ReplyCTCP(msg.Nick, ctcp.Command, c.Version)
		}

	case "PING":
		c.ReplyCTCP(msg.Nick, ctcp.Command, ctcp.Params)

	case "SOURCE":
		if c.Source != "" {
			c.ReplyCTCP(msg.Nick, ctcp.Command, c.Source)
		}

	case "TIME":
		c.ReplyCTCP(msg.Nick, ctcp.Command, time.Now().UTC().Format(time.RFC3339))

	case "USERINFO":
		c.ReplyCTCP(msg.Nick, ctcp.Command, fmt.Sprintf("%s (%s)", c.GetNick(), c.Realname))
	}
}
