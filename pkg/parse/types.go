package parse

import (
	"time"
)

type Post struct {
	// Kind is what kind of post this is. See ValidKinds (required)
	Kind string

	// Timestamps
	Created time.Time
	Updated time.Time

	// Depending on the kind, one of these will host the primary content of the post (required)
	Text string
	Data []byte
	URL  string

	// Title is a title of this post. (optional)
	Title string

	// Description is a short description of the post. (optional)
	Description string

	// Tags are tags related to this pont (optional)
	Tags []string
}
