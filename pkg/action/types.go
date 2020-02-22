package action

import (
	"time"

	"github.com/tstromberg/daily/pkg/parse"
)

// Stream is basically the entire blog.
type Stream struct {
	Posts []*RenderedPost
	Title string

	Timestamp time.Time
}

// RenderedPost is a post along with any dynamically generated data we found
type RenderedPost struct {
	Metadata *parse.Post

	ImageSrc string
	URL      string

	Thumbnails map[string]ThumbOpts
}
