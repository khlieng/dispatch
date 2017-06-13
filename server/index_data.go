package server

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/spf13/viper"

	"github.com/khlieng/dispatch/storage"
)

type connectDefaults struct {
	Name     string   `json:"name"`
	Address  string   `json:"address"`
	Channels []string `json:"channels"`
	Password bool     `json:"password"`
	SSL      bool     `json:"ssl"`
}

type indexData struct {
	Defaults connectDefaults   `json:"defaults"`
	Servers  []storage.Server  `json:"servers,omitempty"`
	Channels []storage.Channel `json:"channels,omitempty"`

	// Users in the selected channel
	Users *Userlist `json:"users,omitempty"`

	// Last messages in the selected channel
	Messages *Messages `json:"messages,omitempty"`
}

func (d *indexData) addUsersAndMessages(server, channel string, session *Session) {
	users := channelStore.GetUsers(server, channel)
	if len(users) > 0 {
		d.Users = &Userlist{
			Server:  server,
			Channel: channel,
			Users:   users,
		}
	}

	messages, hasMore, err := session.user.GetLastMessages(server, channel, 50)
	if err == nil && len(messages) > 0 {
		m := Messages{
			Server:   server,
			To:       channel,
			Messages: messages,
		}

		if hasMore {
			m.Next = messages[0].ID
		}

		d.Messages = &m
	}
}

func getIndexData(r *http.Request, session *Session) *indexData {
	servers := session.user.GetServers()
	connections := session.getConnectionStates()
	for i, server := range servers {
		servers[i].Connected = connections[server.Host]
		servers[i].Port = ""
		servers[i].TLS = false
		servers[i].Password = ""
		servers[i].Username = ""
		servers[i].Realname = ""
	}

	channels := session.user.GetChannels()
	for i, channel := range channels {
		channels[i].Topic = channelStore.GetTopic(channel.Server, channel.Name)
	}

	data := indexData{
		Defaults: connectDefaults{
			Name:     viper.GetString("defaults.name"),
			Address:  viper.GetString("defaults.address"),
			Channels: viper.GetStringSlice("defaults.channels"),
			Password: viper.GetString("defaults.password") != "",
			SSL:      viper.GetBool("defaults.ssl"),
		},
		Servers:  servers,
		Channels: channels,
	}

	server, channel := getTabFromPath(r.URL.EscapedPath())
	if channel != "" {
		data.addUsersAndMessages(server, channel, session)
		return &data
	}

	server, channel = parseTabCookie(r, r.URL.Path)
	if channel != "" {
		for _, ch := range channels {
			if server == ch.Server && channel == ch.Name {
				data.addUsersAndMessages(server, channel, session)
				break
			}
		}
	}

	return &data
}

func getTabFromPath(rawPath string) (string, string) {
	path := strings.Split(strings.Trim(rawPath, "/"), "/")
	if len(path) == 2 {
		name, err := url.PathUnescape(path[1])
		if err == nil && isChannel(name) {
			return path[0], name
		}
	}
	return "", ""
}

func parseTabCookie(r *http.Request, path string) (string, string) {
	if path == "/" {
		cookie, err := r.Cookie("tab")
		if err == nil {
			tab := strings.Split(cookie.Value, "-")

			if len(tab) == 2 && isChannel(tab[1]) {
				return tab[0], tab[1]
			}
		}
	}
	return "", ""
}
