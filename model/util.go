package model

import (
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
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

func valueTypeName(valueType jsoniter.ValueType) string {
	switch valueType {
	case jsoniter.InvalidValue:
		return "INVALID"
	case jsoniter.StringValue:
		return "STRING"
	case jsoniter.NumberValue:
		return "NUMBER"
	case jsoniter.NilValue:
		return "NULL"
	case jsoniter.BoolValue:
		return "BOOLEAN"
	case jsoniter.ArrayValue:
		return "ARRAY"
	case jsoniter.ObjectValue:
		return "OBJECT"
	default:
		return "UNKNOWN"
	}
}

func numberToDate(num int) time.Time {
	year := num / 10000
	month := num / 100 % 100
	day := num % 100
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
}

func dateToNumber(date time.Time) int {
	year, month, day := date.Date()
	return year*10000 + int(month)*100 + day
}
