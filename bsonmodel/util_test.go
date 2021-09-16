package bsonmodel

import (
	"testing"
	"time"

	jsoniter "github.com/json-iterator/go"
)

func TestValue(t *testing.T) {
	xpath := &path{"a.b.c"}
	if xpath.Value() != "a.b.c" {
		t.Errorf("The value expected \"a.b.c\" but was \"%s\"", xpath.Value())
	}
}

func TestIsRoot(t *testing.T) {
	xpath := &path{""}
	if !xpath.IsRoot() {
		t.Error("The value expected true but was false")
	}
	xpath = &path{"any.path.here"}
	if xpath.IsRoot() {
		t.Error("The value expected false but was true")
	}
}

func TestResolve(t *testing.T) {
	var xpath DotNotation = &path{""}
	xpath = xpath.Resolve("a")
	if xpath.Value() != "a" {
		t.Errorf("The value expected \"a\" but was \"%s\"", xpath.Value())
	}
	xpath = xpath.Resolve("b")
	if xpath.Value() != "a.b" {
		t.Errorf("The value expected \"a.b\" but was \"%s\"", xpath.Value())
	}
	xpath = xpath.Resolve("c")
	if xpath.Value() != "a.b.c" {
		t.Errorf("The value expected \"a.b.c\" but was \"%s\"", xpath.Value())
	}
}

func TestResolveIndex(t *testing.T) {
	var xpath DotNotation = &path{"a"}
	xpath = xpath.ResolveIndex(1)
	if xpath.Value() != "a.1" {
		t.Errorf("The value expected \"a.1\" but was \"%s\"", xpath.Value())
	}
	xpath = xpath.ResolveIndex(2)
	if xpath.Value() != "a.1.2" {
		t.Errorf("The value expected \"a.1.2\" but was \"%s\"", xpath.Value())
	}
}

func TestRootPath(t *testing.T) {
	xpath := RootPath()
	if xpath != root {
		t.Error("The value expected root but not")
	}
}

func TestPathOfNames(t *testing.T) {
	xpath := PathOfNames("root")
	if xpath.Value() != "root" {
		t.Errorf("The value expected \"root\" but was \"%s\"", xpath.Value())
	}
	xpath = PathOfNames("root", "a")
	if xpath.Value() != "root.a" {
		t.Errorf("The value expected \"root.a\" but was \"%s\"", xpath.Value())
	}
	xpath = PathOfNames("root", "a", "b")
	if xpath.Value() != "root.a.b" {
		t.Errorf("The value expected \"root.a.b\" but was \"%s\"", xpath.Value())
	}
}

func TestPathOf(t *testing.T) {
	var base DotNotation = &path{"a"}
	xpath := PathOf(base, "b", "c")
	if xpath.Value() != "a.b.c" {
		t.Errorf("The value expected \"a.b.c\" but was \"%s\"", xpath.Value())
	}
}

func TestValueTypeName(t *testing.T) {
	if valueTypeName(jsoniter.InvalidValue) != "INVALID" {
		t.Errorf("The value expected \"INVALID\" nut was \"%s\"", valueTypeName(jsoniter.InvalidValue))
	}
	if valueTypeName(jsoniter.StringValue) != "STRING" {
		t.Errorf("The value expected \"STRING\" nut was \"%s\"", valueTypeName(jsoniter.StringValue))
	}
	if valueTypeName(jsoniter.NumberValue) != "NUMBER" {
		t.Errorf("The value expected \"NUMBER\" nut was \"%s\"", valueTypeName(jsoniter.NumberValue))
	}
	if valueTypeName(jsoniter.NilValue) != "NULL" {
		t.Errorf("The value expected \"NULL\" nut was \"%s\"", valueTypeName(jsoniter.NilValue))
	}
	if valueTypeName(jsoniter.BoolValue) != "BOOLEAN" {
		t.Errorf("The value expected \"BOOLEAN\" nut was \"%s\"", valueTypeName(jsoniter.BoolValue))
	}
	if valueTypeName(jsoniter.ArrayValue) != "ARRAY" {
		t.Errorf("The value expected \"ARRAY\" nut was \"%s\"", valueTypeName(jsoniter.ArrayValue))
	}
	if valueTypeName(jsoniter.ObjectValue) != "OBJECT" {
		t.Errorf("The value expected \"OBJECT\" nut was \"%s\"", valueTypeName(jsoniter.ObjectValue))
	}
}

func TestNumberToDate(t *testing.T) {
	date := numberToDate(20210916)
	year, month, day := date.Date()
	if year != 2021 {
		t.Errorf("The value expected 2021 but was %d", year)
	}
	if month != time.September {
		t.Errorf("The value expected September but was %v", month)
	}
	if day != 16 {
		t.Errorf("The value expected 16 but was %d", day)
	}
}

func TestDateToNumber(t *testing.T) {
	date := time.Now()
	duration := time.Duration(date.Hour())*time.Hour + time.Duration(date.Minute())*time.Minute + time.Duration(date.Second())*time.Second + time.Duration(date.Nanosecond())
	date = date.Add(-duration)
	actual := dateToNumber(date)
	year, month, day := date.Date()
	expected := year*10000 + int(month)*100 + day
	if actual != expected {
		t.Errorf("The valued expected %d but was %d", expected, actual)
	}
}
