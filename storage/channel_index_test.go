package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapChannelListIndex(t *testing.T) {
	i := NewMapChannelListIndex()
	i.Add(&ChannelListItem{
		Name:      "#apples",
		UserCount: 120,
	})
	i.Add(&ChannelListItem{
		Name:      "#cake",
		UserCount: 150,
	})
	i.Add(&ChannelListItem{
		Name:      "#beans",
		UserCount: 12,
	})
	i.Add(&ChannelListItem{
		Name:      "#pie",
		UserCount: 1200,
	})
	i.Add(&ChannelListItem{
		Name:      "#Pork",
		UserCount: 1200,
	})
	i.Finish()

	assert.Len(t, i.Search(""), 5)
	assert.Len(t, i.SearchN("", 0, 20), 5)
	assert.Len(t, i.SearchN("", 0, 1), 1)
	assert.Len(t, i.SearchN("", 1, 1), 1)
	assert.Len(t, i.SearchN("", 0, 0), 0)

	assert.Equal(t, "#pie", i.Search("")[0].Name)
	assert.Equal(t, "#Pork", i.Search("")[1].Name)

	assert.Len(t, i.Search("p"), 2)
	assert.Equal(t, "#Pork", i.Search("p")[1].Name)
}
