package main

import (
	"regexp"
)

var (
	_placeholderPattern = regexp.MustCompile(`^:\w+`)
	_verbatimPattern    = regexp.MustCompile(`^[^*:]+`)
)

// https://docs.netlify.com/routing/headers/#syntax-for-the-headers-file
func PatternToRegexp(pattern string) *regexp.Regexp {
	re := `^`
	for len(pattern) > 0 {
		if pattern[0] == '*' {
			re += ".*"
			pattern = pattern[1:]
			continue
		}
		if pattern[0] == ':' {
			placeholder := _placeholderPattern.FindString(pattern)
			if placeholder == "" {
				impossiblef("the impossible happened: placeholder not found at the beginning of %q", pattern)
			}
			re += "[^/]*"
			pattern = pattern[len(placeholder):]
			continue
		}
		literal := _verbatimPattern.FindString(pattern)
		if literal == "" {
			impossiblef("the impossible happened: literal string not found at the beginning of %q", pattern)
		}
		re += regexp.QuoteMeta(literal)
		pattern = pattern[len(literal):]
	}
	return regexp.MustCompile(re)
}
