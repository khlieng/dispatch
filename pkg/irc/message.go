package irc

import (
	"encoding/json"
	"fmt"
	"path"
	"strconv"
	"strings"
	"unicode"
)

type Message struct {
	Tags    map[string]string
	Prefix  string
	Nick    string
	Command string
	Params  []string
}

type DownloadProgress struct {
	File           string  `json:"file"`
	Error          error   `json:"error"`
	BytesCompleted string  `json:"bytes_completed"`
	BytesRemaining string  `json:"bytes_remaining"`
	PercCompletion float64 `json:"perc_completion"`
	AvgSpeed       string  `json:"avg_speed"`
	InstSpeed      string  `json:"speed"`
	SecondsElapsed int64   `json:"elapsed"`
	SecondsToGo    float64 `json:"eta"`
}

// CTCP is used to parse a message into a CTCP message
type CTCP struct {
	File   string `json:"file"`
	IP     string `json:"ip"`
	Port   uint16 `json:"port"`
	Length uint64 `json:"length"`
}

func (m *Message) LastParam() string {
	if len(m.Params) > 0 {
		return m.Params[len(m.Params)-1]
	}
	return ""
}

// ToCTCP tries to parse the message parameters into a CTCP message
func (m *Message) ToCTCP() *CTCP {
	params := strings.Join(m.Params, " ")
	if strings.Contains(params, "DCC SEND") {
		// to be extra sure that there are non-printable characters
		params = strings.TrimFunc(params, func(r rune) bool {
			return !unicode.IsPrint(r)
		})
		parts := strings.Split(params, " ")
		ip, err := strconv.Atoi(parts[4])
		port, err := strconv.Atoi(parts[5])
		length, err := strconv.Atoi(parts[6])

		if err != nil {
			return nil
		}

		ip3 := uint32ToIP(ip)

		filename := path.Base(parts[3])
		if filename == "/" || filename == "." {
			filename = ""
		}

		return &CTCP{
			File:   filename,
			IP:     ip3,
			Port:   uint16(port),
			Length: uint64(length),
		}
	}
	return nil
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

func uint32ToIP(n int) string {
	var byte1 = n & 255
	var byte2 = ((n >> 8) & 255)
	var byte3 = ((n >> 16) & 255)
	var byte4 = ((n >> 24) & 255)
	return fmt.Sprintf("%d.%d.%d.%d", byte4, byte3, byte2, byte1)
}

func (p DownloadProgress) ToJSON() string {
	progress, err := json.Marshal(p)
	if err != nil {
		return ""
	}
	return string(progress)
}
