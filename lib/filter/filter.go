package filter

import (
	"regexp"

	"github.com/masiulaniec/Dominator/lib/fsutil"
)

func load(filename string) (*Filter, error) {
	lines, err := fsutil.LoadLines(filename)
	if err != nil {
		return nil, err
	}
	return New(lines)
}

func newFilter(filterLines []string) (*Filter, error) {
	var filter Filter
	filter.FilterLines = make([]string, 0)
	for _, line := range filterLines {
		if line != "" {
			filter.FilterLines = append(filter.FilterLines, line)
		}
	}
	if err := filter.compile(); err != nil {
		return nil, err
	}
	return &filter, nil
}

func (filter *Filter) compile() error {
	filter.filterExpressions = make([]*regexp.Regexp, len(filter.FilterLines))
	for index, reEntry := range filter.FilterLines {
		if reEntry == "!" {
			filter.invertMatches = true
			continue
		}
		var err error
		filter.filterExpressions[index], err = regexp.Compile("^" + reEntry)
		if err != nil {
			return err
		}
	}
	return nil
}

func (filter *Filter) match(pathname string) bool {
	if len(filter.filterExpressions) != len(filter.FilterLines) {
		filter.compile()
	}
	defaultRetval := false
	matchRetval := true
	if filter.invertMatches {
		defaultRetval = true
		matchRetval = false
	}
	for _, regex := range filter.filterExpressions {
		if regex != nil && regex.MatchString(pathname) {
			return matchRetval
		}
	}
	return defaultRetval
}

func (filter *Filter) replaceStrings(replaceFunc func(string) string) {
	if filter != nil {
		for index, str := range filter.FilterLines {
			filter.FilterLines[index] = replaceFunc(str)
		}
	}
}
