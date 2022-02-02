package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	var inputFile string
	flag.StringVar(&inputFile, "i", "", "input file (defaults to _headers in the root directory, which will be overwritten)")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), `Usage:
  %s [options] <root>

root is the directory to be published to Netlify. Expanded rules are written to _headers within this root directory.

Options:
`, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatalf("expected 1 argument, got %d; see --help", flag.NArg())
	}
	root := flag.Arg(0)
	if root == "" {
		log.Fatal("root directory is empty string")
	}
	headersFile := filepath.Join(root, "_headers")
	if inputFile == "" {
		inputFile = headersFile
	}

	rules, err := extractRules(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	expanded, err := ExpandRules(root, rules)
	if err != nil {
		log.Fatal(err)
	}
	if err := writeRules(headersFile, expanded); err != nil {
		log.Fatal(err)
	}
}

func extractRules(path string) (rules []Rule, err error) {
	var f *os.File
	f, err = os.Open(path)
	if err != nil {
		return
	}
	defer func() {
		cerr := f.Close()
		if err == nil {
			err = cerr
		}
	}()
	rules, err = ParseHeadersFile(f)
	if err != nil {
		return
	}
	return
}

func writeRules(path string, rules []Rule) (err error) {
	var f *os.File
	f, err = os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		cerr := f.Close()
		if err == nil {
			err = cerr
		}
	}()
	for i, rule := range rules {
		s := rule.String()
		if i != len(rules)-1 {
			s += "\n"
		}
		_, err = f.WriteString(s)
		if err != nil {
			return
		}
	}
	return
}
