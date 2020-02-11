package parse

import (
	"time"
)

type Post struct {
	// Kind is what kind of post this is. See ValidKinds (required)
	Kind string

	// When was the content originally created
	Created time.Time
	// When was the content posted
	Posted time.Time
	// When was the content last updated
	Updated time.Time

	// Depending on the kind, one of these will host the primary content of the post (required)
	Text string
	Data []byte
	URL  string

	// Title is a title of this post. (optional)
	Title string

	// Slug defines the URL name for this post (optional)
	Slug string

	// Description is a short description of the post. (optional)
	Description string

	// Hier is the subdirectory where this post lives
	Hier string

	// Tags are tags related to this pont (optional)
	Tags []string

	// Source is where the post content originated from
	Source string
}
