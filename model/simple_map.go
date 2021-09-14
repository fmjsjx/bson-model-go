package model

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

type identityValueType struct {
}

func (valueType *identityValueType) ToBson(value interface{}) interface{} {
	return value
}

func (valueType *identityValueType) ToData(value interface{}) interface{} {
	return value
}

type intValeType struct {
	identityValueType
}

func (valueType *intValeType) Parse(value interface{}) (interface{}, error) {
	switch value.(type) {
	case int32:
		return int(value.(int32)), nil
	case int64:
		return int(value.(int64)), nil
	case int:
		return value.(int), nil
	case float64:
		return int(value.(float64)), nil
	default:
		e := errors.New(fmt.Sprintf("Type %v can not be cast to type int", reflect.TypeOf(value)))
		return nil, e
	}
}

func (valueType *intValeType) ParseJsoniter(value jsoniter.Any) (interface{}, error) {
	if value.ValueType() == jsoniter.NumberValue {
		return value.ToInt(), nil
	}
	e := errors.New(fmt.Sprintf("The value is not a BOOLEAN (%s)", valueTypeName(value.ValueType())))
	return nil, e
}

type stringValeType struct {
	identityValueType
}

func (valueType *stringValeType) Parse(value interface{}) (interface{}, error) {
	switch value.(type) {
	case string:
		return value.(string), nil
	default:
		e := errors.New(fmt.Sprintf("Type %v can not be cast to type string", reflect.TypeOf(value)))
		return nil, e
	}
}

func (valueType *stringValeType) ParseJsoniter(value jsoniter.Any) (interface{}, error) {
	if value.ValueType() == jsoniter.StringValue {
		return value.ToString(), nil
	}
	e := errors.New(fmt.Sprintf("The value is not a STRING (%s)", valueTypeName(value.ValueType())))
	return nil, e
}

type float64ValeType struct {
	identityValueType
}

func (valueType *float64ValeType) Parse(value interface{}) (interface{}, error) {
	switch value.(type) {
	case int32:
		return float64(value.(int32)), nil
	case int64:
		return float64(value.(int64)), nil
	case int:
		return float64(value.(int)), nil
	case float32:
		return value.(float32), nil
	case float64:
		return value.(float64), nil
	default:
		e := errors.New(fmt.Sprintf("Type %v can not be cast to type float64", reflect.TypeOf(value)))
		return nil, e
	}
}

func (valueType *float64ValeType) ParseJsoniter(value jsoniter.Any) (interface{}, error) {
	if value.ValueType() == jsoniter.NumberValue {
		return value.ToFloat64(), nil
	}
	e := errors.New(fmt.Sprintf("The value is not a NUBER (%s)", valueTypeName(value.ValueType())))
	return nil, e
}

type boolValeType struct {
	identityValueType
}

func (valueType *boolValeType) Parse(value interface{}) (interface{}, error) {
	switch value.(type) {
	case bool:
		return value.(bool), nil
	default:
		e := errors.New(fmt.Sprintf("Type %v can not be cast to type bool", reflect.TypeOf(value)))
		return nil, e
	}
}

func (valueType *boolValeType) ParseJsoniter(value jsoniter.Any) (interface{}, error) {
	if value.ValueType() == jsoniter.BoolValue {
		return value.ToBool(), nil
	}
	e := errors.New(fmt.Sprintf("The value is not a BOOLEAN (%s)", valueTypeName(value.ValueType())))
	return nil, e
}

type datetimeValueType struct {
}

func (valueType *datetimeValueType) Parse(value interface{}) (interface{}, error) {
	if value == nil {
		return nil, nil
	}
	switch value.(type) {
	case primitive.DateTime:
		return value.(primitive.DateTime).Time(), nil
	case time.Time:
		return value.(time.Time), nil
	case primitive.Timestamp:
		t := value.(primitive.Timestamp)
		return time.Unix(int64(t.T), 0), nil
	default:
		e := errors.New(fmt.Sprintf("Type %v can not be cast to type time.Time", reflect.TypeOf(value)))
		return nil, e
	}
}

func (valueType *datetimeValueType) ParseJsoniter(value jsoniter.Any) (interface{}, error) {
	if value.ValueType() == jsoniter.NilValue {
		return nil, nil
	}
	if value.ValueType() == jsoniter.NumberValue {
		return time.UnixMilli(value.ToInt64()), nil
	}
	e := errors.New(fmt.Sprintf("The value is not a NUMBER (%s)", valueTypeName(value.ValueType())))
	return nil, e
}

func (valueType *datetimeValueType) ToBson(value interface{}) interface{} {
	return primitive.NewDateTimeFromTime(value.(time.Time))
}

func (valueType *datetimeValueType) ToData(value interface{}) interface{} {
	return value.(time.Time).UnixMilli()
}

type dateValueType struct {
}

func (valueType *dateValueType) Parse(value interface{}) (interface{}, error) {
	if value == nil {
		return nil, nil
	}
	switch value.(type) {
	case int32:
		return numberToDate(int(value.(int32))), nil
	case int64:
		return numberToDate(int(value.(int64))), nil
	case int:
		return numberToDate(value.(int)), nil
	case float64:
		return numberToDate(int(value.(float64))), nil
	default:
		e := errors.New(fmt.Sprintf("Type %v can not be cast to Date", reflect.TypeOf(value)))
		return nil, e
	}
}

func (valueType *dateValueType) ParseJsoniter(value jsoniter.Any) (interface{}, error) {
	if value.ValueType() == jsoniter.NilValue {
		return nil, nil
	}
	if value.ValueType() == jsoniter.NumberValue {
		return numberToDate(value.ToInt()), nil
	}
	e := errors.New(fmt.Sprintf("The value is not a NUMBER (%s)", valueTypeName(value.ValueType())))
	return nil, e
}

func (valueType *dateValueType) ToBson(value interface{}) interface{} {
	return dateToNumber(value.(time.Time))
}

func (valueType *dateValueType) ToData(value interface{}) interface{} {
	return dateToNumber(value.(time.Time))
}

var intSimpleValueType *intValeType = &intValeType{}
var stringSimpleValueType *stringValeType = &stringValeType{}
var float64SimpleValueType *float64ValeType = &float64ValeType{}
var boolSimpleValueType *boolValeType = &boolValeType{}
var datetimeSimpleValueType *datetimeValueType = &datetimeValueType{}
var dateSimpleValueType *dateValueType = &dateValueType{}

func IntValueType() SimpleValueType {
	return intSimpleValueType
}

func StringValeType() SimpleValueType {
	return stringSimpleValueType
}

func Float64ValueType() SimpleValueType {
	return float64SimpleValueType
}

func BoolValueType() SimpleValueType {
	return boolSimpleValueType
}

func DateTimeSimpleValueType() SimpleValueType {
	return datetimeSimpleValueType
}

func DateSimpleValueType() SimpleValueType {
	return dateSimpleValueType
}
