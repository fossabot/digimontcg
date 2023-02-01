package web

import (
	"net/url"
	"strconv"
	"strings"
)

func readString(qs url.Values, key string) *string {
	s := qs.Get(key)
	if s == "" {
		return nil
	}

	return &s
}

func readCSV(qs url.Values, key string) []string {
	csv := qs.Get(key)
	if csv == "" {
		return []string{}
	}

	return strings.Split(csv, ",")
}

func readInt(qs url.Values, key string, v *Validator) *int {
	s := qs.Get(key)
	if s == "" {
		return nil
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return nil
	}

	return &i
}
