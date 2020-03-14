package daily

import (
	"time"
)

// Item is a post
type Item struct {
	// Kind is what kind of post this is. See ValidKinds (required)
	Kind string

	// When was the content originally created
	Created *time.Time
	// When was the content posted
	Posted *time.Time
	// When was the content last updated
	Updated *time.Time

	// Depending on the kind, one of these will host the primary content of the post (required)
	Text string `yaml:",omitempty"`
	Data []byte `yaml:",omitempty"`
	URL  string `yaml:",omitempty"`

	// Title is a title of this post. (optional)
	Title string `yaml:",omitempty"`

	// Slug defines the URL name for this post (optional)
	Slug string `yaml:",omitempty"`

	// Description is a short description of the post. (optional)
	Description string `yaml:",omitempty"`

	// Hier is the subdirectory where this post lives
	Hier string `yaml:",omitempty"`

	// Tags are tags related to this pont (optional)
	Tags []string `yaml:",omitempty"`

	// Source is where the post content originated from
	Source string `yaml:",omitempty"`
}
