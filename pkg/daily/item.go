package daily

import (
	"time"
)

type YAMLTime struct {
	time.Time
}

func (t *YAMLTime) MarshalYAML() (interface{}, error) {
	panic("marshal")
	//	return []byte(t.Format(time.RFC1123Z)), nil
}

func (t *YAMLTime) UnmarshalYAML(unmarshal func(interface{}) error) (time.Time, error) {
	panic("unmarshal")
	//	return time.Parse(time.RFC1123Z, string(b))
}

func NewYAMLTime(t time.Time) YAMLTime {
	return YAMLTime{t}
}

// FrontMatter is what gets stored in the header of an item (or in YAML sidecar)
type FrontMatter struct {
	// Kind is what kind of post this is. See ValidKinds (required)
	Kind string

	// Draft is if the page is a draft (do not publish)
	Draft bool

	// Posted is when was the content posted
	Posted YAMLTime

	// Title is a title of this post. (optional)
	Title string `yaml:",omitempty"`

	// Description is a short description of the post. (optional)
	Description string `yaml:",omitempty"`

	// Source is where the post content originated from
	Source string `yaml:",omitempty"`
}

// Item is a post to be rendered
type Item struct {
	FrontMatter FrontMatter
	// Content is inline content
	Content string
	// Path is where the input is saved on disk
	Path string
	// RelPath is where the item was found relative to the input directory
	RelPath string
	// Format is the format to store in
	Format string
}
