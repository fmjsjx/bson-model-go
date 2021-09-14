package model

import (
	"strconv"
	"strings"
)

// DotNotation defines BSON dot notation.
type DotNotation interface {
	Value() string
	IsRoot() bool
	Resolve(name string) DotNotation
	ResolveIndex(index int) DotNotation
}

type path struct {
	value string
}

func (xpath *path) Value() string {
	return xpath.value
}

func (xpath *path) IsRoot() bool {
	return xpath.value == ""
}

func (xpath *path) Resolve(name string) DotNotation {
	if xpath.IsRoot() {
		return &path{name}
	}
	return &path{xpath.value + "." + name}
}

func (xpath *path) ResolveIndex(index int) DotNotation {
	name := strconv.Itoa(index)
	return xpath.Resolve(name)
}

func (xpath *path) String() string {
	return xpath.value
}

var root *path = &path{}

func RootPath() DotNotation {
	return root
}

func PathOfNames(base string, names ...string) DotNotation {
	if len(names) == 0 {
		return &path{base}
	}
	var builder strings.Builder
	builder.WriteString(base)
	for i := 0; i < len(names); i++ {
		builder.WriteString(".")
		builder.WriteString(names[i])
	}
	return &path{builder.String()}
}

func PathOf(base DotNotation, names ...string) DotNotation {
	return PathOfNames(base.Value(), names...)
}
