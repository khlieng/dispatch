This is Go package to generate guid (globally unique id) with good properties.

Usage:
```go
import "github.com/kjk/betterguid"

id := betterguid.New()
fmt.Printf("guid: '%s'\n", id)
```

Generated guids have good properties:
* they're 20 character strings, safe for inclusion in urls (don't require escaping)
* they're based on timestamp; they sort **after** any existing ids
* they contain 72-bits of random data after the timestamp so that IDs won't
  collide with other IDs
* they sort **lexicographically** (the timestamp is converted to a string
  that will sort correctly)
* they're monotonically increasing. Even if you generate more than one in the
  same timestamp, the latter ones will sort after the former ones. We do this
  by using the previous random bits but "incrementing" them by 1 (only in the
  case of a timestamp collision).

Read https://www.firebase.com/blog/2015-02-11-firebase-unique-identifiers.html
for more info.

Based on https://gist.github.com/mikelehen/3596a30bd69384624c11

You can read [Generating good, random and unique ids in Go](https://blog.kowalczyk.info/article/JyRZ/generating-good-random-and-unique-ids-in-go.html) to see how it compares to other similar libraries.
