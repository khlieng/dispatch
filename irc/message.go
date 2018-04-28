package irc

import (
	"strings"
	"sync"

	"github.com/spf13/cast"
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
	for _, param := range params[1 : len(params)-1] {
		parts := strings.SplitN(param, "=", 2)
		i.lock.Lock()
		if parts[0][0] == '-' {
			delete(i.support, parts[0][1:])
		} else if len(parts) == 2 {
			i.support[parts[0]] = parts[1]
		} else {
			i.support[param] = ""
		}
		i.lock.Unlock()
	}
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
