//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package bleve

import (
	"expvar"
	"io/ioutil"
	"log"
	"time"

	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/analyzers/custom_analyzer"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/analyzers/keyword_analyzer"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/analyzers/simple_analyzer"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/analyzers/standard_analyzer"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/byte_array_converters/ignore" // token filters

	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/byte_array_converters/json"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/byte_array_converters/string"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/char_filters/html_char_filter"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/char_filters/regexp_char_filter"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/char_filters/zero_width_non_joiner"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/datetime_parsers/datetime_optional" // analyzers
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/datetime_parsers/flexible_go"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/language/ar" // languages
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/language/bg"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/language/ca"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/language/cjk"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/language/ckb"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/language/cs"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/language/da"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/language/de"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/language/el"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/language/en"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/language/es"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/language/eu"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/language/fa"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/language/fi"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/language/fr"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/language/ga"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/language/gl"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/language/hi"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/language/hu"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/language/hy"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/language/id"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/language/in"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/language/it"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/language/nl"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/language/no"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/language/pt"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/language/ro"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/language/ru"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/language/sv"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/language/th"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/language/tr"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/token_filters/apostrophe_filter" // kv stores
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/token_filters/compound"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/token_filters/edge_ngram_filter"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/token_filters/elision_filter"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/token_filters/keyword_marker_filter"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/token_filters/length_filter"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/token_filters/lower_case_filter"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/token_filters/ngram_filter"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/token_filters/shingle"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/token_filters/stop_tokens_filter"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/token_filters/truncate_token_filter"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/token_filters/unicode_normalize"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/token_map"            // tokenizers
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/tokenizers/exception" // fragment formatters
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/tokenizers/regexp_tokenizer"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/tokenizers/single_token"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/tokenizers/unicode"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/analysis/tokenizers/whitespace_tokenizer"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/index/store/boltdb" // date time parsers
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/index/store/goleveldb"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/index/store/gtreap"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/index/store/inmem"
	"github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/index/upside_down" // byte array converters
	"github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/registry"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/search/highlight/fragment_formatters/ansi"

	// token maps
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/search/highlight/fragment_formatters/html"
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/search/highlight/fragmenters/simple"  // fragmenters
	_ "github.com/khlieng/name_pending/Godeps/_workspace/src/github.com/blevesearch/bleve/search/highlight/highlighters/simple" // highlighters
)

var bleveExpVar = expvar.NewMap("bleve")

type configuration struct // char filters
{
	Cache                  *registry.Cache
	DefaultHighlighter     string
	DefaultKVStore         string
	SlowSearchLogThreshold time.Duration
	analysisQueue          *upside_down.AnalysisQueue
}

func newConfiguration() *configuration {
	return &configuration{
		Cache:         registry.NewCache(),
		analysisQueue: upside_down.NewAnalysisQueue(4),
	}
}

// Config contains library level configuration
var Config *configuration

func init() {
	bootStart := time.Now()

	// build the default configuration
	Config = newConfiguration()

	_, err := Config.Cache.DefineFragmentFormatter("highlightSpanHTML",
		map[string]interface{}{
			"type":   "html",
			"before": `<span class="highlight">`,
			"after":  `</span>`,
		})
	if err != nil {
		panic(err)
	}

	_, err = Config.Cache.DefineHighlighter("html",
		map[string]interface{}{
			"type":       "simple",
			"fragmenter": "simple",
			"formatter":  "highlightSpanHTML",
		})
	if err != nil {
		panic(err)
	}

	_, err = Config.Cache.DefineHighlighter("ansi",
		map[string]interface{}{
			"type":       "simple",
			"fragmenter": "simple",
			"formatter":  "ansi",
		})
	if err != nil {
		panic(err)
	}

	// set the default highlighter
	Config.DefaultHighlighter = "html"

	// default kv store
	Config.DefaultKVStore = "boltdb"

	bootDuration := time.Since(bootStart)
	bleveExpVar.Add("bootDuration", int64(bootDuration))
}

var logger = log.New(ioutil.Discard, "bleve", log.LstdFlags)

// SetLog sets the logger used for logging
// by default log messages are sent to ioutil.Discard
func SetLog(l *log.Logger) {
	logger = l
}
