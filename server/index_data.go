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
	Networks []*storage.Network
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
		Defaults: connectDefaults{
			Defaults:       &cfg.Defaults,
			ServerPassword: cfg.Defaults.ServerPassword != "",
		},
		HexIP: cfg.HexIP,
		Version: dispatchVersion{
			Tag:    version.Tag,
			Commit: version.Commit,
			Date:   version.Date,
		},
	}

	if state == nil {
		data.Settings = storage.DefaultClientSettings()
		return &data
	}

	data.Settings = state.user.ClientSettings()

	state.lock.Lock()
	for _, network := range state.networks {
		network = network.Copy()
		network.Password = ""
		network.Username = ""
		network.Realname = ""

		data.Networks = append(data.Networks, network)
		data.Channels = append(data.Channels, network.Channels()...)
	}
	state.lock.Unlock()

	openDMs, err := state.user.OpenDMs()
	if err == nil {
		data.OpenDMs = openDMs
	}

	tab, err := tabFromRequest(r)
	if err == nil && hasTab(data.Channels, openDMs, tab.Network, tab.Name) {
		data.addUsersAndMessages(tab.Network, tab.Name, state)
	}

	return &data
}

func (d *indexData) addUsersAndMessages(network, name string, state *State) {
	if i, ok := state.client(network); ok && isChannel(name) {
		if users := i.ChannelUsers(name); len(users) > 0 {
			d.Users = &Userlist{
				Network: network,
				Channel: name,
				Users:   users,
			}
		}
	}

	messages, hasMore, err := state.user.LastMessages(network, name, 50)
	if err == nil && len(messages) > 0 {
		m := Messages{
			Network:  network,
			To:       name,
			Messages: messages,
		}

		if hasMore {
			m.Next = messages[0].ID
		}

		d.Messages = &m
	}
}

func hasTab(channels []*storage.Channel, openDMs []storage.Tab, network, name string) bool {
	if name != "" {
		for _, ch := range channels {
			if network == ch.Network && name == ch.Name {
				return true
			}
		}

		for _, tab := range openDMs {
			if network == tab.Network && name == tab.Name {
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
			tab.Network = parts[0]
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

			tab.Network = parts[0]
		}
	}

	return tab, nil
}
