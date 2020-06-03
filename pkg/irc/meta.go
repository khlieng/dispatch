package irc

// GetNickChannels returns the channels the client has in common with
// the user that changed nick
func GetNickChannels(msg *Message) []string {
	return stringListMeta(msg)
}

// GetQuitChannels returns the channels the client has in common with
// the user that quit
func GetQuitChannels(msg *Message) []string {
	return stringListMeta(msg)
}

func GetMode(msg *Message) *Mode {
	if mode, ok := msg.meta.(*Mode); ok {
		return mode
	}
	return nil
}

// GetNamreplyUsers returns all RPL_NAMREPLY users
// when passed a RPL_ENDOFNAMES message
func GetNamreplyUsers(msg *Message) []string {
	return stringListMeta(msg)
}

func stringListMeta(msg *Message) []string {
	if list, ok := msg.meta.([]string); ok {
		return list
	}
	return nil
}
