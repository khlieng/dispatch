package server

import (
	"net/http"
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
	Messages []storage.Message `json:"messages,omitempty"`
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

	params := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(params) == 2 && isChannel(params[1]) {
		users := channelStore.GetUsers(params[0], params[1])
		if len(users) > 0 {
			data.Users = &Userlist{
				Server:  params[0],
				Channel: params[1],
				Users:   users,
			}
		}
	}

	return &data
}
