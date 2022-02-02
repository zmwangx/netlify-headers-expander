package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
)

type Annotations struct {
	NoExpand           bool
	ExcludeDirectories bool
}

type Rule struct {
	Path    string
	Headers []string
	Annotations
}

func (r Rule) String() string {
	s := r.Path + "\n"
	for _, header := range r.Headers {
		s += "  " + header + "\n"
	}
	return s
}

func ParseHeadersFile(r io.Reader) (rules []Rule, err error) {
	var annotations Annotations
	var rule Rule

	registerCurrentRule := func() {
		if rule.Path == "" {
			return
		}
		if len(rule.Headers) == 0 {
			log.Printf("warning: no headers found for %q", rule.Path)
			return
		}
		rules = append(rules, rule)
	}

	s := bufio.NewScanner(r)
	linePattern := regexp.MustCompile(`^(\s*)(.*)$`)
	for s.Scan() {
		line := s.Text()
		matches := linePattern.FindStringSubmatch(line)
		if matches == nil {
			impossiblef("line %q does not match %q", line, linePattern.String())
		}
		indent := matches[1]
		content := strings.TrimSpace(matches[2])
		if content == "" {
			continue
		}
		if content[0] == '#' {
			content = strings.TrimSpace(content[1:])
			if content == "no-expand" {
				annotations.NoExpand = true
			} else if content == "exclude-directories" {
				annotations.ExcludeDirectories = true
			}
			continue
		}
		if indent == "" {
			// Path, begin a new rule.
			if rule.Path != "" {
				registerCurrentRule()
			}
			rule = Rule{Path: content, Annotations: annotations}
			annotations = Annotations{}
		} else {
			// Header.
			if rule.Path == "" {
				err = fmt.Errorf("header %q found without associated path", content)
				return
			}
			rule.Headers = append(rule.Headers, content)
		}
	}
	registerCurrentRule()
	return
}
