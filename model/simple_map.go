package model

import (
	jsoniter "github.com/json-iterator/go"
)

type SimpleValueType interface {
	Parse(value interface{}) (interface{}, error)
	ParseJsoniter(value jsoniter.Any) (interface{}, error)
	ToBson(value interface{}) interface{}
	ToData(value interface{}) interface{}
}

type baseSimpleMap struct {
	baseMap
	valueType SimpleValueType
}
