package irc

import (
	"strings"
)

type Message struct {
	Prefix  string
	Nick    string
	Command string
	Params  []string
}

func (m *Message) LastParam() string {
	if len(m.Params) > 0 {
		return m.Params[len(m.Params)-1]
	}
	return ""
}

func parseMessage(line string) *Message {
	line = strings.Trim(line, "\r\n ")
	msg := Message{}
	cmdStart := 0
	cmdEnd := len(line)

	if strings.HasPrefix(line, ":") {
		cmdStart = strings.Index(line, " ") + 1

		if cmdStart > 0 {
			msg.Prefix = line[1 : cmdStart-1]
		} else {
			return nil
		}

		if i := strings.Index(msg.Prefix, "!"); i > 0 {
			msg.Nick = msg.Prefix[:i]
		} else if i := strings.Index(msg.Prefix, "@"); i > 0 {
			msg.Nick = msg.Prefix[:i]
		} else {
			msg.Nick = msg.Prefix
		}
	}

	var trailing string

	if i := strings.Index(line, " :"); i > 0 {
		cmdEnd = i
		trailing = line[i+2:]
	}

	cmd := strings.Fields(line[cmdStart:cmdEnd])
	if len(cmd) == 0 {
		return nil
	}
	msg.Command = cmd[0]

	if len(cmd) > 1 {
		msg.Params = cmd[1:]
	}

	if cmdEnd != len(line) {
		msg.Params = append(msg.Params, trailing)
	}

	return &msg
}
