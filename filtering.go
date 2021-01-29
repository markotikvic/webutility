package webutility

import (
	"net/http"
	"net/url"
	"strings"
)

// Filter ...
type Filter map[string][]string

// Get ...
func (f Filter) Get(key string) (values []string, ok bool) {
	values, ok = f[key]
	return values, ok
}

// Count ...
func (f Filter) Count() int {
	return len(f)
}

// Add ...
func (f Filter) Add(key, val string) {
	f[key] = append(f[key], val)
}

// ValueAt ...
func (f Filter) ValueAt(val string, index int) string {
	if filter, ok := f[val]; ok {
		if len(filter) > index {
			return filter[index]
		}
	}

	return ""
}

func (f Filter) Validate(validFilters []string) (Filter, bool) {
	goodFilters := make(Filter)
	cnt, len := 0, 0
	for fi := range f {
		len++
		for _, v := range validFilters {
			if fi == v {
				cnt++
				goodFilters[fi] = f[fi]
			}
		}
	}

	result := true
	if len > 0 && cnt == 0 {
		// if no valid filters are found declare filtering request as invalid
		result = false
	}

	return goodFilters, result
}

// ParseFilters requires input in format: "param1::value1|param2::value2..."
func ParseFilters(req *http.Request, header string) (filters Filter) {
	q := req.FormValue(header)
	q = strings.Trim(q, "\"")
	kvp := strings.Split(q, "|")
	filters = make(Filter, len(kvp))

	for i := range kvp {
		kv := strings.Split(kvp[i], "::")
		if len(kv) == 2 {
			key, _ := url.QueryUnescape(kv[0])

			// get values (if more than 1)
			vals := strings.Split(kv[1], ",")
			for _, v := range vals {
				u, _ := url.QueryUnescape(v)
				filters[key] = append(filters[key], u)
			}
		}
	}

	return filters
}
