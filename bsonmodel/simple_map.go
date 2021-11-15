package bsonmodel

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
	"unsafe"

	mapset "github.com/deckarep/golang-set"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
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
		return NumberToDate(int(value.(int32))), nil
	case int64:
		return NumberToDate(int(value.(int64))), nil
	case int:
		return NumberToDate(value.(int)), nil
	case float64:
		return NumberToDate(int(value.(float64))), nil
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
		return NumberToDate(value.ToInt()), nil
	}
	e := errors.New(fmt.Sprintf("The value is not a NUMBER (%s)", valueTypeName(value.ValueType())))
	return nil, e
}

func (valueType *dateValueType) ToBson(value interface{}) interface{} {
	return DateToNumber(value.(time.Time))
}

func (valueType *dateValueType) ToData(value interface{}) interface{} {
	return DateToNumber(value.(time.Time))
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

type IntSimpleMapModel interface {
	mapModel
	Keys() []int
	Get(key int) interface{}
	Put(key int, value interface{}) interface{}
	Remove(key int) bool
}

type intSimpleMap struct {
	baseSimpleMap
	data map[int]interface{}
}

func (imap *intSimpleMap) Size() int {
	return len(imap.data)
}

func (imap *intSimpleMap) Clear() {
	imap.updatedKeys.Clear()
	removedKeys := imap.removedKeys
	data := imap.data
	for k := range data {
		removedKeys.Add(k)
	}
	for k := range data {
		delete(data, k)
	}
}

func (imap *intSimpleMap) Keys() []int {
	data := imap.data
	keys := make([]int, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	return keys
}

func (imap *intSimpleMap) Get(key int) interface{} {
	value, ok := imap.data[key]
	if ok {
		return value
	}
	return nil
}

func (imap *intSimpleMap) Put(key int, value interface{}) interface{} {
	data := imap.data
	old, ok := data[key]
	if ok {
		if old != value {
			data[key] = value
			imap.updatedKeys.Add(key)
		}
		return old
	}
	data[key] = value
	imap.updatedKeys.Add(key)
	imap.removedKeys.Remove(key)
	return nil
}

func (imap *intSimpleMap) Remove(key int) bool {
	data := imap.data
	_, ok := data[key]
	if ok {
		delete(data, key)
		imap.updatedKeys.Remove(key)
		imap.removedKeys.Add(key)
		return true
	}
	return false
}

func (imap *intSimpleMap) ToBson() interface{} {
	return imap.ToDocument()
}

func (imap *intSimpleMap) ToData() interface{} {
	data := make(map[int]interface{})
	valueType := imap.valueType
	for key, value := range imap.data {
		data[key] = valueType.ToData(value)
	}
	return data
}

func (imap *intSimpleMap) Reset() {
	imap.updatedKeys.Clear()
	imap.removedKeys.Clear()
}

func (imap *intSimpleMap) LoadJsoniter(any jsoniter.Any) error {
	imap.Reset()
	data := imap.data
	for k := range data {
		delete(data, k)
	}
	if any.ValueType() == jsoniter.ObjectValue {
		valueType := imap.valueType
		keys := any.Keys()
		for _, key := range keys {
			k, err := strconv.Atoi(key)
			if err != nil {
				// skip key that not be an int
				continue
			}
			v, err := valueType.ParseJsoniter(any.Get(key))
			if err != nil {
				return err
			}
			data[k] = v
		}
	}
	return nil
}

func (imap *intSimpleMap) AppendUpdates(updates bson.M) bson.M {
	data := imap.data
	updatedKeys := imap.updatedKeys
	if updatedKeys.Cardinality() > 0 {
		dset := FixedEmbedded(updates, "$set")
		valueType := imap.valueType
		for _, uk := range updatedKeys.ToSlice() {
			key := uk.(int)
			value := data[key]
			k := strconv.Itoa(key)
			name := imap.XPath().Resolve(k)
			v := valueType.ToBson(value)
			dset[name.Value()] = v
		}
	}
	removedKeys := imap.removedKeys
	if removedKeys.Cardinality() > 0 {
		unset := FixedEmbedded(updates, "$unset")
		for _, uk := range removedKeys.ToSlice() {
			key := uk.(int)
			k := strconv.Itoa(key)
			name := imap.XPath().Resolve(k)
			unset[name.Value()] = 1
		}
	}
	return updates
}

func (imap *intSimpleMap) ToDocument() bson.M {
	doc := bson.M{}
	valueType := imap.valueType
	for k, v := range imap.data {
		key := strconv.Itoa(k)
		value := valueType.ToBson(v)
		doc[key] = value
	}
	return doc
}

func (imap *intSimpleMap) LoadDocument(document bson.M) error {
	imap.Reset()
	data := imap.data
	for k := range data {
		delete(data, k)
	}
	valueType := imap.valueType
	for key, value := range document {
		k, err := strconv.Atoi(key)
		if err != nil {
			// skip key that not be an int
			continue
		}
		v, err := valueType.Parse(value)
		if err != nil {
			return err
		}
		data[k] = v
	}
	return nil
}

func (imap *intSimpleMap) ToSync() interface{} {
	sync := make(map[string]interface{})
	updatedKeys := imap.updatedKeys
	data := imap.data
	if updatedKeys.Cardinality() > 0 {
		valueType := imap.valueType
		for _, uk := range updatedKeys.ToSlice() {
			key := uk.(int)
			value := data[key]
			sync[strconv.Itoa(key)] = valueType.ToData(value)
		}
	}
	return sync
}

func (imap *intSimpleMap) ToDelete() interface{} {
	delete := make(map[string]int)
	removedKeys := imap.removedKeys
	if removedKeys.Cardinality() > 0 {
		for _, uk := range removedKeys.ToSlice() {
			key := uk.(int)
			delete[strconv.Itoa(key)] = 1
		}
	}
	return delete
}

func (imap *intSimpleMap) ToDataJson() (string, error) {
	return jsoniter.MarshalToString(imap.ToData())
}

func (imap *intSimpleMap) ToSyncJson() (string, error) {
	return jsoniter.MarshalToString(imap.ToSync())
}

func (imap *intSimpleMap) ToDeleteJson() (string, error) {
	return jsoniter.MarshalToString(imap.ToDelete())
}

func NewIntSimpleMapModel(parent BsonModel, name string, valueType SimpleValueType) IntSimpleMapModel {
	mapModel := &intSimpleMap{}
	mapModel.parent = parent
	mapModel.name = name
	mapModel.updatedKeys = mapset.NewThreadUnsafeSet()
	mapModel.removedKeys = mapset.NewThreadUnsafeSet()
	mapModel.valueType = valueType
	mapModel.data = make(map[int]interface{})
	return mapModel
}

type StringSimpleMapModel interface {
	mapModel
	Keys() []string
	Get(key string) interface{}
	Put(key string, value interface{}) interface{}
	Remove(key string) bool
}

type stringSimpleMap struct {
	baseSimpleMap
	data map[string]interface{}
}

func (smap *stringSimpleMap) Size() int {
	return len(smap.data)
}

func (smap *stringSimpleMap) Clear() {
	smap.updatedKeys.Clear()
	removedKeys := smap.removedKeys
	data := smap.data
	for k := range data {
		removedKeys.Add(k)
	}
	for k := range data {
		delete(data, k)
	}
}

func (smap *stringSimpleMap) Keys() []string {
	data := smap.data
	keys := make([]string, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	return keys
}

func (smap *stringSimpleMap) Get(key string) interface{} {
	value, ok := smap.data[key]
	if ok {
		return value
	}
	return nil
}

func (smap *stringSimpleMap) Put(key string, value interface{}) interface{} {
	data := smap.data
	old, ok := data[key]
	if ok {
		if old != value {
			data[key] = value
			smap.updatedKeys.Add(key)
		}
		return old
	}
	data[key] = value
	smap.updatedKeys.Add(key)
	smap.removedKeys.Remove(key)
	return nil
}

func (smap *stringSimpleMap) Remove(key string) bool {
	data := smap.data
	_, ok := data[key]
	if ok {
		delete(data, key)
		smap.updatedKeys.Remove(key)
		smap.removedKeys.Add(key)
		return true
	}
	return false
}

func (smap *stringSimpleMap) ToBson() interface{} {
	return smap.ToDocument()
}

func (smap *stringSimpleMap) ToData() interface{} {
	data := make(map[string]interface{})
	valueType := smap.valueType
	for key, value := range smap.data {
		data[key] = valueType.ToData(value)
	}
	return data
}

func (smap *stringSimpleMap) Reset() {
	smap.updatedKeys.Clear()
	smap.removedKeys.Clear()
}

func (smap *stringSimpleMap) LoadJsoniter(any jsoniter.Any) error {
	smap.Reset()
	data := smap.data
	for k := range data {
		delete(data, k)
	}
	if any.ValueType() == jsoniter.ObjectValue {
		valueType := smap.valueType
		keys := any.Keys()
		for _, key := range keys {
			v, err := valueType.ParseJsoniter(any.Get(key))
			if err != nil {
				return err
			}
			data[key] = v
		}
	}
	return nil
}

func (smap *stringSimpleMap) AppendUpdates(updates bson.M) bson.M {
	data := smap.data
	updatedKeys := smap.updatedKeys
	if updatedKeys.Cardinality() > 0 {
		dset := FixedEmbedded(updates, "$set")
		valueType := smap.valueType
		for _, uk := range updatedKeys.ToSlice() {
			key := uk.(string)
			name := smap.XPath().Resolve(key)
			value := data[key]
			v := valueType.ToBson(value)
			dset[name.Value()] = v
		}
	}
	removedKeys := smap.removedKeys
	if removedKeys.Cardinality() > 0 {
		unset := FixedEmbedded(updates, "$unset")
		for _, uk := range removedKeys.ToSlice() {
			key := uk.(string)
			name := smap.XPath().Resolve(key)
			unset[name.Value()] = 1
		}
	}
	return updates
}

func (smap *stringSimpleMap) ToDocument() bson.M {
	doc := bson.M{}
	valueType := smap.valueType
	for key, v := range smap.data {
		value := valueType.ToBson(v)
		doc[key] = value
	}
	return doc
}

func (smap *stringSimpleMap) LoadDocument(document bson.M) error {
	smap.Reset()
	data := smap.data
	for k := range data {
		delete(data, k)
	}
	valueType := smap.valueType
	for key, value := range document {
		v, err := valueType.Parse(value)
		if err != nil {
			return err
		}
		data[key] = v
	}
	return nil
}

func (smap *stringSimpleMap) ToSync() interface{} {
	sync := make(map[string]interface{})
	updatedKeys := smap.updatedKeys
	data := smap.data
	if updatedKeys.Cardinality() > 0 {
		valueType := smap.valueType
		for _, uk := range updatedKeys.ToSlice() {
			key := uk.(string)
			value := data[key]
			sync[key] = valueType.ToData(value)
		}
	}
	return sync
}

func (smap *stringSimpleMap) ToDelete() interface{} {
	delete := make(map[string]int)
	removedKeys := smap.removedKeys
	if removedKeys.Cardinality() > 0 {
		for _, uk := range removedKeys.ToSlice() {
			key := uk.(string)
			delete[key] = 1
		}
	}
	return delete
}

func (smap *stringSimpleMap) ToDataJson() (string, error) {
	return jsoniter.MarshalToString(smap.ToData())
}

func (smap *stringSimpleMap) ToSyncJson() (string, error) {
	return jsoniter.MarshalToString(smap.ToSync())
}

func (smap *stringSimpleMap) ToDeleteJson() (string, error) {
	return jsoniter.MarshalToString(smap.ToDelete())
}

func NewStringSimpleMapModel(parent BsonModel, name string, valueType SimpleValueType) StringSimpleMapModel {
	mapModel := &stringSimpleMap{}
	mapModel.parent = parent
	mapModel.name = name
	mapModel.updatedKeys = mapset.NewThreadUnsafeSet()
	mapModel.removedKeys = mapset.NewThreadUnsafeSet()
	mapModel.valueType = valueType
	mapModel.data = make(map[string]interface{})
	return mapModel
}

type intSimpleMapEncoder struct{}

func (codec *intSimpleMapEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	imap := ((*intSimpleMap)(ptr))
	return len(imap.data) == 0
}

func (codec *intSimpleMapEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	imap := ((*intSimpleMap)(ptr))
	stream.WriteVal(imap.data)
}

type stringSimpleMapEncoder struct{}

func (codec *stringSimpleMapEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	smap := ((*stringSimpleMap)(ptr))
	return len(smap.data) == 0
}

func (codec *stringSimpleMapEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	smap := ((*stringSimpleMap)(ptr))
	stream.WriteVal(smap.data)
}

func init() {
	jsoniter.RegisterTypeEncoder("bsonmodel.intSimpleMap", &intSimpleMapEncoder{})
	jsoniter.RegisterTypeEncoder("bsonmodel.stringSimpleMap", &stringSimpleMapEncoder{})
}
