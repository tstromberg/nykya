package tree

import (
	"time"
)

type Stream struct {
	Posts []*Post
	Title string

	Timestamp time.Time
}

type Post struct {
	Kind string
	Source string

	Title string
	Content string
	URL string

	Tags []string

	Timestamp time.Time

	Thumbnails map[string]string
}
