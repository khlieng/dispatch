package irc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFeatures(t *testing.T) {
	s := NewFeatures()
	featureTransforms["CAKE"] = toInt
	s.Parse([]string{"bob", "CAKE=31", "PIE", ":durr"})
	assert.Equal(t, 31, s.Int("CAKE"))
	assert.Equal(t, "", s.String("CAKE"))
	assert.True(t, s.Has("CAKE"))
	assert.True(t, s.Has("PIE"))
	assert.False(t, s.Has("APPLES"))
	assert.Equal(t, "", s.String("APPLES"))
	assert.Equal(t, 0, s.Int("APPLES"))

	s.Parse([]string{"bob", "-PIE", ":hurr"})
	assert.False(t, s.Has("PIE"))

	s.Parse([]string{"bob", "CAKE=1337", ":durr"})
	assert.Equal(t, 1337, s.Int("CAKE"))

	s.Parse([]string{"bob", "CAKE=", ":durr"})
	assert.Equal(t, "", s.String("CAKE"))
	assert.True(t, s.Has("CAKE"))

	delete(featureTransforms, "CAKE")
	s.Parse([]string{"bob", "CAKE===", ":durr"})
	assert.Equal(t, "==", s.String("CAKE"))

	s.Parse([]string{"bob", "-CAKE=31", ":durr"})
	assert.False(t, s.Has("CAKE"))

	s.Parse([]string{"bob", "CHANLIMIT=#&:50", ":durr"})
	assert.Equal(t, map[string]int{"#": 50, "&": 50}, s.Get("CHANLIMIT"))

	s.Parse([]string{"bob", "CHANLIMIT=#:50,&:25", ":durr"})
	assert.Equal(t, map[string]int{"#": 50, "&": 25}, s.Get("CHANLIMIT"))

	s.Parse([]string{"bob", "CHANLIMIT=&:50,#:", ":durr"})
	assert.Equal(t, map[string]int{"#": 0, "&": 50}, s.Get("CHANLIMIT"))

	s.Parse([]string{"bob", "CHANTYPES=#&", ":durr"})
	assert.Equal(t, []string{"#", "&"}, s.Get("CHANTYPES"))
}
