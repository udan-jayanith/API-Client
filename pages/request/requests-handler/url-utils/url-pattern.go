package url_utils

import (
	attr "Zbolt/pages/request/requests-handler/attributes"
	"strings"
)

type Pattern struct {
	raw_path string
	Values   map[string]int
	List     []attr.Attribute
}

func (r *Pattern) Path() string {
	p := r.raw_path
	for _, attr := range r.List {
		p = strings.Replace(p, "{"+attr.Key+"}", attr.Value, 1)
	}
	return p
}

func (r *Pattern) Set(key, value string) {
	r.List[r.Values[key]].Value = value
}

func ParsePattern(pattern string) (Pattern, error) {
	list := make([]attr.Attribute, 0, 4)
	value := make(map[string]int, 4)

	i := -1
	var idx int
	var attr attr.Attribute

	for j, char := range pattern {
		if char == '{' && i == -1 {
			i = j
		} else if char == '}' && i > -1 {
			attr.Key = pattern[i+1 : j]
			value[attr.Key] = idx
			list = append(list, attr)
			idx++
			i = -1
		}
	}

	return Pattern{
		raw_path: pattern,
		Values:   value,
		List:     list,
	}, nil
}
