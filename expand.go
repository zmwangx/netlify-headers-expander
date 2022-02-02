package main

import (
	"strings"
)

func ExpandRules(root string, rules []Rule) (expanded []Rule, err error) {
	paths, err := ListPaths(root)
	if err != nil {
		return
	}
	for _, rule := range rules {
		re := PatternToRegexp(rule.Path)
		var matching []string
		if !rule.NoExpand {
			for _, path := range paths {
				if strings.HasSuffix(path, "/") && rule.ExcludeDirectories {
					continue
				}
				if re.MatchString(path) {
					matching = append(matching, path)
				}
			}
		}
		if len(matching) == 0 {
			matching = append(matching, rule.Path)
		}
		for _, path := range matching {
			expanded = append(expanded, Rule{Path: path, Headers: rule.Headers})
		}
	}
	return
}
