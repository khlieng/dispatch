package irc

import (
	"strings"
)

type Message struct {
	Prefix   string
	Nick     string
	Command  string
	Params   []string
	Trailing string
}

func parseMessage(line string) *Message {
	line = strings.Trim(line, "\r\n")
	msg := Message{}
	cmdStart := 0
	cmdEnd := len(line)

	if strings.HasPrefix(line, ":") {
		cmdStart = strings.Index(line, " ") + 1
		msg.Prefix = line[1 : cmdStart-1]

		if i := strings.Index(msg.Prefix, "!"); i > 0 {
			msg.Nick = msg.Prefix[:i]
		} else if i := strings.Index(msg.Prefix, "@"); i > 0 {
			msg.Nick = msg.Prefix[:i]
		} else {
			msg.Nick = msg.Prefix
		}
	}

	if i := strings.Index(line, " :"); i > 0 {
		cmdEnd = i
		msg.Trailing = line[i+2:]
	}

	cmd := strings.Split(line[cmdStart:cmdEnd], " ")
	msg.Command = cmd[0]
	if len(cmd) > 1 {
		msg.Params = cmd[1:]
	}

	if msg.Trailing != "" {
		msg.Params = append(msg.Params, msg.Trailing)
	}

	return &msg
}
