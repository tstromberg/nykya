package store

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/secure/precis"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"k8s.io/klog"
)

// slugRemove is chars to be removed from slug calculation
var slugRemove = regexp.MustCompile(`[^a-zA-Z0-9 -]`)

// slugReplace are chars that become spaces
var slugReplace = regexp.MustCompile(`\s+`)

// maxSlugWords is how many words to consider when sluggin'
var maxSlugWords = 5

func slugify(in string) string {
	klog.Infof("slugify: %q", in)

	// TODO: does not work for umlauts yet?
	loose := precis.NewIdentifier(
		precis.AdditionalMapping(func() transform.Transformer {
			return transform.Chain(norm.NFD, transform.RemoveFunc(func(r rune) bool {
				return unicode.Is(unicode.Mn, r)
			}))
		}),
		precis.Norm(norm.NFC),
	)

	p, err := loose.String(in)
	if err != nil {
		klog.Warningf("loose string reported error: %v", err)
		p = in
	}

	p = slugRemove.ReplaceAllString(p, "")
	p = slugReplace.ReplaceAllString(p, " ")

	words := strings.Split(strings.ToLower(p), " ")
	slug := strings.Join(words[0:len(words)], "-")
	if len(words) > maxSlugWords {
		slug = strings.Join(words[0:maxSlugWords], "-")
	}
	if slug == "" {
		return "untitled"
	}
	return slug
}
