package irc

import (
	"strings"
	"sync"

	"github.com/spf13/cast"
)

type Message struct {
	Tags    map[string]string
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

	if strings.HasPrefix(line, "@") {
		next := strings.Index(line, " ")
		if next == -1 {
			return nil
		}
		tags := strings.Split(line[1:next], ";")

		if len(tags) > 0 {
			msg.Tags = map[string]string{}
		}

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

		line = line[next+1:]
	}

	if strings.HasPrefix(line, ":") {
		next := strings.Index(line, " ")
		if next == -1 {
			return nil
		}
		msg.Prefix = line[1:next]

		if i := strings.Index(msg.Prefix, "!"); i > 0 {
			msg.Nick = msg.Prefix[:i]
		} else if i := strings.Index(msg.Prefix, "@"); i > 0 {
			msg.Nick = msg.Prefix[:i]
		} else {
			msg.Nick = msg.Prefix
		}

		line = line[next+1:]
	}

	cmdEnd := len(line)
	trailing := ""
	if i := strings.Index(line, " :"); i > 0 {
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

type iSupport struct {
	support map[string]string
	lock    sync.Mutex
}

func newISupport() *iSupport {
	return &iSupport{
		support: map[string]string{},
	}
}

func (i *iSupport) parse(params []string) {
	i.lock.Lock()
	for _, param := range params[1 : len(params)-1] {
		key, val := splitParam(param)
		if key == "" {
			continue
		}

		if key[0] == '-' {
			delete(i.support, key[1:])
		} else {
			i.support[key] = val
		}
	}
	i.lock.Unlock()
}

func (i *iSupport) Has(key string) bool {
	i.lock.Lock()
	_, has := i.support[key]
	i.lock.Unlock()
	return has
}

func (i *iSupport) Get(key string) string {
	i.lock.Lock()
	v := i.support[key]
	i.lock.Unlock()
	return v
}

func (i *iSupport) GetInt(key string) int {
	i.lock.Lock()
	v := cast.ToInt(i.support[key])
	i.lock.Unlock()
	return v
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
