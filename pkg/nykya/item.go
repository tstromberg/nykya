package nykya

import (
	"time"
)

// YAMLTime is time serializable to frontmatter
type YAMLTime struct {
	time.Time
}

// MarshalYAML marshals time into RFC 1123
func (yt YAMLTime) MarshalYAML() (interface{}, error) {
	return yt.Format(time.RFC1123Z), nil
}

// UnmarshalYAML unmarshals RFC1123 or Y-M-D timestamps to time.Time
func (yt YAMLTime) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	err := unmarshal(&s)
	if err != nil {
		return err
	}
	tp, err := time.Parse(time.RFC1123Z, s)
	if err != nil {
		// Be forgiving
		tp, err = time.Parse("2006-01-02", s)
		if err != nil {
			return err
		}
	}
	yt.Time = tp
	return nil
}

// NewYAMLTime returns a populated YAMLTime object
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

// RawItem is a post to be rendered
type RawItem struct {
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
