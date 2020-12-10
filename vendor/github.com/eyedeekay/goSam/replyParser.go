package goSam

import (
	"fmt"
	"strings"
)

// The Possible Results send by SAM
const (
	ResultOk             = "OK"              //Operation completed successfully
	ResultCantReachPeer  = "CANT_REACH_PEER" //The peer exists, but cannot be reached
	ResultDuplicatedID   = "DUPLICATED_ID"   //If the nickname is already associated with a session :
	ResultDuplicatedDest = "DUPLICATED_DEST" //The specified Destination is already in use
	ResultI2PError       = "I2P_ERROR"       //A generic I2P error (e.g. I2CP disconnection, etc.)
	ResultInvalidKey     = "INVALID_KEY"     //The specified key is not valid (bad format, etc.)
	ResultKeyNotFound    = "KEY_NOT_FOUND"   //The naming system can't resolve the given name
	ResultPeerNotFound   = "PEER_NOT_FOUND"  //The peer cannot be found on the network
	ResultTimeout        = "TIMEOUT"         // Timeout while waiting for an event (e.g. peer answer)
)

// A ReplyError is a custom error type, containing the Result and full Reply
type ReplyError struct {
	Result string
	Reply  *Reply
}

func (r ReplyError) Error() string {
	return fmt.Sprintf("ReplyError: Result:%s - Reply:%+v", r.Result, r.Reply)
}

// Reply is the parsed result of a SAM command, containing a map of all the key-value pairs
type Reply struct {
	Topic string
	Type  string
	From  string
	To    string

	Pairs map[string]string
}

func parseReply(line string) (*Reply, error) {
	line = strings.TrimSpace(line)
	parts := strings.Split(line, " ")
	if len(parts) < 3 {
		return nil, fmt.Errorf("Malformed Reply.\n%s\n", line)
	}
	preParseReply := func() []string {
		val := ""
		quote := false
		for _, v := range parts {
			if strings.Contains(v, "=\"") {
				quote = true
			}
			if strings.Contains(v, "\"\n") || strings.Contains(v, "\" ") {
				quote = false
			}
			if quote {
				val += v + "_"
			} else {
				val += v + " "
			}
		}
		return strings.Split(strings.TrimSuffix(strings.TrimSpace(val), "_"), " ")
	}
	parts = preParseReply()

	r := &Reply{
		Topic: parts[0],
		Type:  parts[1],
		Pairs: make(map[string]string, len(parts)-2),
	}

	for _, v := range parts[2:] {
		if strings.Contains(v, "FROM_PORT") {
			if v != "FROM_PORT=0" {
				r.From = v
			}
		} else if strings.Contains(v, "TO_PORT") {
			if v != "TO_PORT=0" {
				r.To = v
			}
		} else {
			kvPair := strings.SplitN(v, "=", 2)
			if kvPair != nil {
				if len(kvPair) != 2 {
					return nil, fmt.Errorf("Malformed key-value-pair len != 2.\n%s\n", kvPair)
				}
			}
			r.Pairs[kvPair[0]] = kvPair[len(kvPair)-1]
		}
	}

	return r, nil
}
