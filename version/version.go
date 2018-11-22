package version

import (
	"fmt"
	"strings"
	"time"
)

var (
	Tag    = "dev"
	Commit = "none"
	Date   = "unknown"
)

func init() {
	vParts := strings.Split(Tag, "-")

	if len(vParts) > 1 {
		Tag = fmt.Sprintf("%s + %s commits", vParts[0], vParts[1])
	}

	t, err := time.Parse(time.RFC3339, Date)
	if err == nil {
		Date = t.Format("02 Jan 2006, 15:04:05")
	}
}
