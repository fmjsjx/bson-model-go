package model

import (
	mapset "github.com/deckarep/golang-set"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
)

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
	bson := bson.M{}
	valueType := smap.valueType
	for key, value := range smap.data {
		bson[key] = valueType.ToBson(value)
	}
	return bson
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
		dset, ok := updates["$set"].(bson.M)
		if !ok {
			dset = bson.M{}
		}
		valueType := smap.valueType
		for _, uk := range updatedKeys.ToSlice() {
			key := uk.(string)
			value := data[key]
			v := valueType.ToBson(value)
			dset[key] = v
		}
	}
	removedKeys := smap.removedKeys
	if removedKeys.Cardinality() > 0 {
		unset, ok := updates["$unset"].(bson.M)
		if !ok {
			unset = bson.M{}
		}
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

func NewStringSimpleMapModel(parent MongoModel, name string, valueType SimpleValueType) StringSimpleMapModel {
	mapModel := &stringSimpleMap{}
	mapModel.parent = parent
	mapModel.name = name
	mapModel.updatedKeys = mapset.NewThreadUnsafeSet()
	mapModel.removedKeys = mapset.NewThreadUnsafeSet()
	mapModel.valueType = valueType
	mapModel.data = make(map[string]interface{})
	return mapModel
}
