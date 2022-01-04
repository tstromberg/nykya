package nykya

import (
	"time"
)

// YAMLTime is time serializable to frontmatter
type YAMLTime struct {
	time.Time
}

var ymdFormat = "2006-01-02"

// MarshalYAML marshals time into RFC 1123
func (yt *YAMLTime) MarshalYAML() (interface{}, error) {
	return yt.Format(ymdFormat), nil
}

// UnmarshalYAML unmarshals RFC1123 or Y-M-D timestamps to time.Time
func (yt *YAMLTime) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	err := unmarshal(&s)
	if err != nil {
		return err
	}

	tp, err := time.Parse(ymdFormat, s)
	yt.Time = tp
	return err
}

// NewYAMLTime returns a populated YAMLTime object
func NewYAMLTime(t time.Time) YAMLTime {
	return YAMLTime{t}
}

// FrontMatter is metadata about an item saved to disk
type FrontMatter struct {
	// Kind is what kind of post this is. See ValidKinds (required)
	Kind string

	// Draft is if the page is a draft (do not publish)
	Draft bool

	// Date is when was the content posted
	Date YAMLTime `yaml:"date"`

	// Title is a title of this post. (optional)
	Title string `yaml:",omitempty"`

	// Description is a short description of the post. (optional)
	Description string `yaml:",omitempty"`

	// Origin is where the post content originated from
	Origin string `yaml:",omitempty"`
}

// RenderInput is ephemeral metadata for a post to be rendered
type RenderInput struct {
	FrontMatter FrontMatter

	// Inline is inline content
	Inline string

	// ContentPath is relative path to content (not the sidecar)
	ContentPath string

	// Format is the format of the content
	Format string
}
