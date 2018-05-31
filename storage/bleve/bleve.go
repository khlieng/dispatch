package bleve

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"

	"github.com/khlieng/dispatch/storage"
)

// Bleve implements storage.MessageSearchProvider
type Bleve struct {
	index bleve.Index
}

func New(path string) (*Bleve, error) {
	index, err := bleve.Open(path)
	if err == bleve.ErrorIndexPathDoesNotExist {
		keywordMapping := bleve.NewTextFieldMapping()
		keywordMapping.Analyzer = keyword.Name
		keywordMapping.Store = false
		keywordMapping.IncludeTermVectors = false
		keywordMapping.IncludeInAll = false

		contentMapping := bleve.NewTextFieldMapping()
		contentMapping.Analyzer = "en"
		contentMapping.Store = false
		contentMapping.IncludeTermVectors = false
		contentMapping.IncludeInAll = false

		messageMapping := bleve.NewDocumentMapping()
		messageMapping.StructTagKey = "bleve"
		messageMapping.AddFieldMappingsAt("server", keywordMapping)
		messageMapping.AddFieldMappingsAt("to", keywordMapping)
		messageMapping.AddFieldMappingsAt("content", contentMapping)

		mapping := bleve.NewIndexMapping()
		mapping.AddDocumentMapping("message", messageMapping)

		index, err = bleve.New(path, mapping)
	}
	if err != nil {
		return nil, err
	}
	return &Bleve{index: index}, nil
}

func (b *Bleve) Index(id string, message *storage.Message) error {
	return b.index.Index(id, message)
}

func (b *Bleve) SearchMessages(server, channel, q string) ([]string, error) {
	serverQuery := bleve.NewMatchQuery(server)
	serverQuery.SetField("server")
	channelQuery := bleve.NewMatchQuery(channel)
	channelQuery.SetField("to")
	contentQuery := bleve.NewMatchQuery(q)
	contentQuery.SetField("content")
	contentQuery.SetFuzziness(2)

	query := bleve.NewBooleanQuery()
	query.AddMust(serverQuery, channelQuery, contentQuery)

	search := bleve.NewSearchRequest(query)
	searchResults, err := b.index.Search(search)
	if err != nil {
		return nil, err
	}

	ids := make([]string, len(searchResults.Hits))
	for i, hit := range searchResults.Hits {
		ids[i] = hit.ID
	}

	return ids, nil
}

func (b *Bleve) Close() {
	b.index.Close()
}
