package model

import (
	"strconv"

	mapset "github.com/deckarep/golang-set"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
)

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
	bson := bson.M{}
	valueType := imap.valueType
	for key, value := range imap.data {
		bson[strconv.Itoa(key)] = valueType.ToBson(value)
	}
	return bson
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
		dset, ok := updates["$set"].(bson.M)
		if !ok {
			dset = bson.M{}
		}
		valueType := imap.valueType
		for _, uk := range updatedKeys.ToSlice() {
			key := uk.(int)
			value := data[key]
			k := strconv.Itoa(key)
			v := valueType.ToBson(value)
			dset[k] = v
		}
	}
	removedKeys := imap.removedKeys
	if removedKeys.Cardinality() > 0 {
		unset, ok := updates["$unset"].(bson.M)
		if !ok {
			unset = bson.M{}
		}
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

func NewIntSimpleMapModel(parent MongoModel, name string, valueType SimpleValueType) IntSimpleMapModel {
	mapModel := &intSimpleMap{}
	mapModel.parent = parent
	mapModel.name = name
	mapModel.updatedKeys = mapset.NewThreadUnsafeSet()
	mapModel.removedKeys = mapset.NewThreadUnsafeSet()
	mapModel.valueType = valueType
	mapModel.data = make(map[int]interface{})
	return mapModel
}
