package bsonmodel

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func FixedEmbedded(m bson.M, name string) bson.M {
	doc, ok := m[name].(bson.M)
	if !ok {
		doc = bson.M{}
		m[name] = doc
	}
	return doc
}

func IntValue(m bson.M, name string, def int) (int, error) {
	v := m[name]
	if v == nil {
		return def, nil
	}
	switch v.(type) {
	case int32:
		return int(v.(int32)), nil
	case int64:
		return int(v.(int64)), nil
	case float64:
		return int(v.(float64)), nil
	case int:
		return v.(int), nil
	default:
		return def, errors.New(fmt.Sprintf("Type %v can not be cast to type int", reflect.TypeOf(v)))
	}
}

func Float64Value(m bson.M, name string, def float64) (float64, error) {
	v := m[name]
	if v == nil {
		return def, nil
	}
	switch v.(type) {
	case int32:
		return float64(v.(int32)), nil
	case int64:
		return float64(v.(int64)), nil
	case float64:
		return v.(float64), nil
	case int:
		return float64(v.(int)), nil
	default:
		return def, errors.New(fmt.Sprintf("Type %v can not be cast to type int", reflect.TypeOf(v)))
	}
}

func StringValue(m bson.M, name string, def string) (string, error) {
	v := m[name]
	if v == nil {
		return def, nil
	}
	switch v.(type) {
	case string:
		return v.(string), nil
	default:
		return def, errors.New(fmt.Sprintf("Type %v can not be cast to type string", reflect.TypeOf(v)))
	}
}

func DateTimeValue(m bson.M, name string) (t time.Time, err error) {
	v := m[name]
	if v == nil {
		return
	}
	switch v.(type) {
	case primitive.DateTime:
		t = v.(primitive.DateTime).Time()
	case primitive.Timestamp:
		t = time.Unix(int64(v.(primitive.Timestamp).T), 0)
	default:
		err = errors.New(fmt.Sprintf("Type %v can not be cast to type time.Time", reflect.TypeOf(v)))
	}
	return
}

func DateValue(m bson.M, name string) (t time.Time, err error) {
	v := m[name]
	if v == nil {
		return
	}
	switch v.(type) {
	case int32:
		t = numberToDate(int(v.(int32)))
	case int64:
		t = numberToDate(int(v.(int64)))
	case float64:
		t = numberToDate(int(v.(float64)))
	case int:
		t = numberToDate(v.(int))
	default:
		err = errors.New(fmt.Sprintf("Type %v can not be cast to type int", reflect.TypeOf(v)))
	}
	return
}

func EmbeddedValue(m bson.M, name string) (bson.M, error) {
	v := m[name]
	if v == nil {
		return nil, nil
	}
	switch v.(type) {
	case bson.M:
		return v.(bson.M), nil
	default:
		return nil, errors.New(fmt.Sprintf("Type %v can not be cast to type bson.M", reflect.TypeOf(v)))
	}
}

func AnyIntValue(any jsoniter.Any, def int) (int, error) {
	switch any.ValueType() {
	case jsoniter.NilValue:
		return def, nil
	case jsoniter.InvalidValue:
		return def, nil
	case jsoniter.NumberValue:
		return any.ToInt(), nil
	default:
		return def, errors.New(fmt.Sprintf("The value is not a NUMBER (%s)", valueTypeName(any.ValueType())))
	}
}

func AnyFloat64Value(any jsoniter.Any, def float64) (float64, error) {
	switch any.ValueType() {
	case jsoniter.NilValue:
		return def, nil
	case jsoniter.InvalidValue:
		return def, nil
	case jsoniter.NumberValue:
		return any.ToFloat64(), nil
	default:
		return def, errors.New(fmt.Sprintf("The value is not a NUMBER (%s)", valueTypeName(any.ValueType())))
	}
}

func AnyStringValue(any jsoniter.Any, def string) (string, error) {
	switch any.ValueType() {
	case jsoniter.NilValue:
		return def, nil
	case jsoniter.InvalidValue:
		return def, nil
	case jsoniter.StringValue:
		return any.ToString(), nil
	default:
		return def, errors.New(fmt.Sprintf("The value is not a STRING (%s)", valueTypeName(any.ValueType())))
	}
}

func AnyDateTimeValue(any jsoniter.Any) (t time.Time, err error) {
	switch any.ValueType() {
	case jsoniter.NilValue:
		return
	case jsoniter.InvalidValue:
		return
	case jsoniter.NumberValue:
		t = time.UnixMilli(any.ToInt64())
	default:
		err = errors.New(fmt.Sprintf("The value is not a NUMBER (%s)", valueTypeName(any.ValueType())))
	}
	return
}

func AnyDateValue(any jsoniter.Any) (t time.Time, err error) {
	switch any.ValueType() {
	case jsoniter.NilValue:
		return
	case jsoniter.InvalidValue:
		return
	case jsoniter.NumberValue:
		t = numberToDate(any.ToInt())
	default:
		err = errors.New(fmt.Sprintf("The value is not a NUMBER (%s)", valueTypeName(any.ValueType())))
	}
	return
}
