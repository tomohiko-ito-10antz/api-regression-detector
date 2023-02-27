package errors

import (
	"fmt"
	"strings"

	"golang.org/x/exp/maps"
)

type Info map[string]any

func (info Info) String() string {
	values := []string{}
	for k, v := range info {
		values = append(values, fmt.Sprintf(`%s=%#v`, k, v))
	}
	return strings.Join(values, ", ")
}

func (info Info) With(key string, val any) Info {
	m := maps.Clone(info)
	m[key] = val
	return m
}

func (info Info) AppendTo(message string) string {
	return fmt.Sprintf("%s (%s)", message, info.String())
}
