package irc

import (
	"strconv"
	"strings"
	"sync"
)

type Features struct {
	m    map[string]interface{}
	lock sync.Mutex
}

func NewFeatures() *Features {
	return &Features{
		m: map[string]interface{}{},
	}
}

func (f *Features) Map() map[string]interface{} {
	m := map[string]interface{}{}
	f.lock.Lock()
	for k, v := range f.m {
		m[k] = v
	}
	f.lock.Unlock()
	return m
}

func (f *Features) Parse(params []string) {
	f.lock.Lock()
	for _, param := range params[1 : len(params)-1] {
		key, val := splitParam(param)
		if key == "" {
			continue
		}

		if key[0] == '-' {
			delete(f.m, key[1:])
		} else {
			if t, ok := featureTransforms[key]; ok {
				f.m[key] = t(val)
			} else {
				f.m[key] = val
			}
		}
	}

	f.lock.Unlock()
}

func (f *Features) Has(key string) bool {
	f.lock.Lock()
	_, has := f.m[key]
	f.lock.Unlock()
	return has
}

func (f *Features) Get(key string) interface{} {
	f.lock.Lock()
	v := f.m[key]
	f.lock.Unlock()
	return v
}

func (f *Features) String(key string) string {
	if v, ok := f.Get(key).(string); ok {
		return v
	}
	return ""
}

func (f *Features) Int(key string) int {
	if v, ok := f.Get(key).(int); ok {
		return v
	}
	return 0
}

type featureTransform func(interface{}) interface{}

func toInt(v interface{}) interface{} {
	s := v.(string)
	if s == "" {
		return 0
	}

	i, _ := strconv.Atoi(s)
	return i
}

func toCharList(v interface{}) interface{} {
	s := v.(string)
	list := make([]string, len(s))
	for i := range s {
		list[i] = s[i : i+1]
	}
	return list
}

func parseChanlimit(v interface{}) interface{} {
	limits := map[string]int{}

	pairs := strings.Split(v.(string), ",")
	for _, p := range pairs {
		pair := strings.Split(p, ":")

		if len(pair) == 2 {
			prefixes := pair[0]
			limit, _ := strconv.Atoi(pair[1])

			for i := range prefixes {
				limits[prefixes[i:i+1]] = limit
			}
		}
	}

	return limits
}

var featureTransforms = map[string]featureTransform{
	"AWAYLEN":     toInt,
	"CHANLIMIT":   parseChanlimit,
	"CHANNELLEN":  toInt,
	"CHANTYPES":   toCharList,
	"HOSTLEN":     toInt,
	"KICKLEN":     toInt,
	"MAXCHANNELS": toInt,
	"MAXTARGETS":  toInt,
	"MODES":       toInt,
	"NICKLEN":     toInt,
	"TOPICLEN":    toInt,
	"USERLEN":     toInt,
}
