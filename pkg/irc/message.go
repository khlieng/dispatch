package irc

import (
	"strings"
)

type Message struct {
	Tags    map[string]string
	Sender  string
	Ident   string
	Host    string
	Command string
	Params  []string

	meta interface{}
}

func (m *Message) LastParam() string {
	if len(m.Params) > 0 {
		return m.Params[len(m.Params)-1]
	}
	return ""
}

func (m *Message) IsFromServer() bool {
	return m.Sender == "" || strings.Contains(m.Sender, ".")
}

func (m *Message) ToCTCP() *CTCP {
	return DecodeCTCP(m.LastParam())
}

func ParseMessage(line string) *Message {
	msg := Message{}

	if strings.HasPrefix(line, "@") {
		next := strings.Index(line, " ")
		if next == -1 {
			return nil
		}
		tags := strings.Split(line[1:next], ";")

		if len(tags) > 0 {
			msg.Tags = map[string]string{}

			for _, tag := range tags {
				key, val := splitParam(tag)
				if key == "" {
					continue
				}

				if val != "" {
					msg.Tags[key] = unescapeTag(val)
				} else {
					msg.Tags[key] = ""
				}
			}
		}

		for line[next+1] == ' ' {
			next++
		}
		line = line[next+1:]
	}

	if strings.HasPrefix(line, ":") {
		next := strings.Index(line, " ")
		if next == -1 {
			return nil
		}
		prefix := line[1:next]

		if i := strings.Index(prefix, "!"); i > 0 {
			msg.Sender = prefix[:i]
			prefix = prefix[i+1:]

			if i = strings.Index(prefix, "@"); i > 0 {
				msg.Ident = prefix[:i]
				msg.Host = prefix[i+1:]
			} else {
				msg.Ident = prefix
			}
		} else if i = strings.Index(prefix, "@"); i > 0 {
			msg.Sender = prefix[:i]
			msg.Host = prefix[i+1:]
		} else {
			msg.Sender = prefix
		}

		line = line[next+1:]
	}

	cmdEnd := len(line)
	trailing := ""
	if i := strings.Index(line, " :"); i >= 0 {
		cmdEnd = i
		trailing = line[i+2:]
	}

	cmd := strings.Fields(line[:cmdEnd])
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

func splitParam(param string) (string, string) {
	parts := strings.SplitN(param, "=", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return parts[0], ""
}

var unescapeTagReplacer = strings.NewReplacer(
	"\\:", ";",
	"\\s", " ",
	"\\\\", "\\",
	"\\r", "\r",
	"\\n", "\n",
)

func unescapeTag(s string) string {
	return unescapeTagReplacer.Replace(s)
}
