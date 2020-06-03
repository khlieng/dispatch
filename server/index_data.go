package server

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/khlieng/dispatch/config"
	"github.com/khlieng/dispatch/storage"
	"github.com/khlieng/dispatch/version"
)

type connectDefaults struct {
	*config.Defaults
	ServerPassword bool
}

type dispatchVersion struct {
	Tag    string
	Commit string
	Date   string
}

type indexData struct {
	Defaults connectDefaults
	Servers  []Server
	Channels []*storage.Channel
	OpenDMs  []storage.Tab
	HexIP    bool
	Version  dispatchVersion

	Settings *storage.ClientSettings

	// Users in the selected channel
	Users *Userlist

	// Last messages in the selected channel
	Messages *Messages
}

func (d *Dispatch) getIndexData(r *http.Request, state *State) *indexData {
	cfg := d.Config()

	data := indexData{
		Defaults: connectDefaults{Defaults: &cfg.Defaults},
		HexIP:    cfg.HexIP,
		Version: dispatchVersion{
			Tag:    version.Tag,
			Commit: version.Commit,
			Date:   version.Date,
		},
	}

	data.Defaults.ServerPassword = cfg.Defaults.ServerPassword != ""

	if state == nil {
		data.Settings = storage.DefaultClientSettings()
		return &data
	}

	data.Settings = state.user.GetClientSettings()

	servers, err := state.user.GetServers()
	if err != nil {
		return nil
	}
	connections := state.getConnectionStates()
	for _, server := range servers {
		server.Password = ""
		server.Username = ""
		server.Realname = ""

		s := Server{
			Server: server,
			Status: newConnectionUpdate(server.Host, connections[server.Host]),
		}

		if i, ok := state.irc[server.Host]; ok {
			s.Features = i.Features.Map()
		}

		data.Servers = append(data.Servers, s)
	}

	channels, err := state.user.GetChannels()
	if err != nil {
		return nil
	}
	for i, channel := range channels {
		if client, ok := state.getIRC(channel.Server); ok {
			channels[i].Topic = client.ChannelTopic(channel.Name)
		}
	}
	data.Channels = channels

	openDMs, err := state.user.GetOpenDMs()
	if err != nil {
		return nil
	}
	data.OpenDMs = openDMs

	tab, err := tabFromRequest(r)
	if err == nil && hasTab(channels, openDMs, tab.Server, tab.Name) {
		data.addUsersAndMessages(tab.Server, tab.Name, state)
	}

	return &data
}

func (d *indexData) addUsersAndMessages(server, name string, state *State) {
	if i, ok := state.getIRC(server); ok && isChannel(name) {
		if users := i.ChannelUsers(name); len(users) > 0 {
			d.Users = &Userlist{
				Server:  server,
				Channel: name,
				Users:   users,
			}
		}
	}

	messages, hasMore, err := state.user.GetLastMessages(server, name, 50)
	if err == nil && len(messages) > 0 {
		m := Messages{
			Server:   server,
			To:       name,
			Messages: messages,
		}

		if hasMore {
			m.Next = messages[0].ID
		}

		d.Messages = &m
	}
}

func hasTab(channels []*storage.Channel, openDMs []storage.Tab, server, name string) bool {
	if name != "" {
		for _, ch := range channels {
			if server == ch.Server && name == ch.Name {
				return true
			}
		}

		for _, tab := range openDMs {
			if server == tab.Server && name == tab.Name {
				return true
			}
		}
	}
	return false
}

func tabFromRequest(r *http.Request) (Tab, error) {
	tab := Tab{}

	var path string
	if strings.HasPrefix(r.URL.Path, "/ws") {
		path = r.URL.EscapedPath()[3:]
	} else {
		referer, err := url.Parse(r.Referer())
		if err != nil {
			return tab, err
		}

		path = referer.EscapedPath()
	}

	if path == "/" {
		cookie, err := r.Cookie("tab")
		if err != nil {
			return tab, err
		}

		v, err := url.PathUnescape(cookie.Value)
		if err != nil {
			return tab, err
		}

		parts := strings.SplitN(v, ";", 2)
		if len(parts) == 2 {
			tab.Server = parts[0]
			tab.Name = parts[1]
		}
	} else {
		parts := strings.Split(strings.Trim(path, "/"), "/")
		if len(parts) > 0 && len(parts) < 3 {
			if len(parts) == 2 {
				name, err := url.PathUnescape(parts[1])
				if err != nil {
					return tab, err
				}

				tab.Name = name
			}

			tab.Server = parts[0]
		}
	}

	return tab, nil
}
