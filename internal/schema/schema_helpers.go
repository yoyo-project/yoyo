package schema

import (
	"regexp"
	"strings"
)

var splitter = regexp.MustCompile("[-_]")

func pascal(in string) string {
	ss := splitter.Split(in, -1)
	for i := range ss {
		ss[i] = strings.Title(ss[i])
	}

	return strings.Join(ss, "")
}
