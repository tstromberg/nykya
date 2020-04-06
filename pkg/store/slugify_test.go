package store

import "testing"

func TestSlugify(t *testing.T) {
	var tests = []struct {
		in  string
		out string
	}{
		{"what?", "what"},
		{"Sometimes it works!", "sometimes-it-works"},
		{"what should we do", "what-should-we-do"},
		{"what should we do when it all comes down to you?", "what-should-we-do-when"},
		{"\r\n.....x\r\n", "x"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got := slugify(tt.in)
			if got != tt.out {
				t.Errorf("got %q, want %q", got, tt.out)
			}
		})
	}
}
