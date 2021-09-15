package model

import (
	"strconv"

	mapset "github.com/deckarep/golang-set"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
)

type IntObjectMapValueModel interface {
	MapValueModel
	setKey(key int)
	Key() int
}

type IntObjectMapModel interface {
	mapModel
	Keys() []int
	Get(key int) IntObjectMapValueModel
	Put(key int, value IntObjectMapValueModel) IntObjectMapValueModel
	Remove(key int) bool
}

type IntObjectMapValueFactory func() IntObjectMapValueModel

type intObjectMap struct {
	baseMap
	valueFactory IntObjectMapValueFactory
	data         map[int]IntObjectMapValueModel
}

func (imap *intObjectMap) Size() int {
	return len(imap.data)
}

func (imap *intObjectMap) Clear() {
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

func (imap *intObjectMap) Keys() []int {
	data := imap.data
	keys := make([]int, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	return keys
}

func (imap *intObjectMap) Get(key int) IntObjectMapValueModel {
	value, ok := imap.data[key]
	if ok {
		return value
	}
	return nil
}

func (imap *intObjectMap) Put(key int, value IntObjectMapValueModel) IntObjectMapValueModel {
	data := imap.data
	old, ok := data[key]
	if ok {
		if old != value {
			data[key] = value
			imap.updatedKeys.Add(key)
			value.setKey(key)
			value.setParent(imap)
			value.SetFullyUpdate(true)
			old.unbind()
		}
		return old
	}
	data[key] = value
	imap.updatedKeys.Add(key)
	imap.removedKeys.Remove(key)
	value.setKey(key)
	value.setParent(imap)
	value.SetFullyUpdate(true)
	return nil
}

func (imap *intObjectMap) Remove(key int) bool {
	data := imap.data
	old, ok := data[key]
	if ok {
		delete(data, key)
		imap.updatedKeys.Remove(key)
		imap.removedKeys.Add(key)
		old.unbind()
		return true
	}
	return false
}

func (imap *intObjectMap) ToBson() interface{} {
	bson := bson.M{}
	for key, value := range imap.data {
		bson[strconv.Itoa(key)] = value.ToBson()
	}
	return bson
}

func (imap *intObjectMap) ToData() interface{} {
	data := make(map[int]interface{})
	for key, value := range imap.data {
		data[key] = value.ToData()
	}
	return data
}

func (imap *intObjectMap) Reset() {
	imap.updatedKeys.Clear()
	imap.removedKeys.Clear()
	for _, v := range imap.data {
		v.Reset()
	}
}

func (imap *intObjectMap) LoadJsoniter(any jsoniter.Any) error {
	imap.Reset()
	data := imap.data
	for k, v := range data {
		v.unbind()
		delete(data, k)
	}
	if any.ValueType() == jsoniter.ObjectValue {
		valueFactory := imap.valueFactory
		keys := any.Keys()
		for _, key := range keys {
			k, err := strconv.Atoi(key)
			if err != nil {
				// skip key that not be an int
				continue
			}
			value := valueFactory()
			err = value.LoadJsoniter(any.Get(key))
			if err != nil {
				return err
			}
			data[k] = value
			value.setParent(imap)
			value.setKey(k)
		}
	}
	return nil
}

func (imap *intObjectMap) AppendUpdates(updates bson.M) bson.M {
	data := imap.data
	updatedKeys := imap.updatedKeys
	if updatedKeys.Cardinality() > 0 {
		dset, ok := updates["$set"].(bson.M)
		if !ok {
			dset = bson.M{}
		}
		for _, uk := range updatedKeys.ToSlice() {
			key := uk.(int)
			value := data[key]
			if value.FullyUpdate() {
				dset[value.XPath().Value()] = value.ToBson()
			} else {
				value.AppendUpdates(updates)
			}
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

func (imap *intObjectMap) ToDocument() bson.M {
	doc := bson.M{}
	for k, v := range imap.data {
		key := strconv.Itoa(k)
		value := v.ToBson()
		doc[key] = value
	}
	return doc
}

func (imap *intObjectMap) LoadDocument(document bson.M) error {
	imap.Reset()
	data := imap.data
	for k := range data {
		delete(data, k)
	}
	valueFactory := imap.valueFactory
	for key, value := range document {
		k, err := strconv.Atoi(key)
		if err != nil {
			// skip key that not be an int
			continue
		}
		obj, ok := value.(bson.M)
		if !ok {
			// skip value that not be a bson.M
			continue
		}
		v := valueFactory()
		err = v.LoadDocument(obj)
		if err != nil {
			return err
		}
		data[k] = v
	}
	return nil
}

func NewIntObjectMapModel(parent BsonModel, name string, valueFactory IntObjectMapValueFactory) IntObjectMapModel {
	mapModel := &intObjectMap{}
	mapModel.parent = parent
	mapModel.name = name
	mapModel.updatedKeys = mapset.NewThreadUnsafeSet()
	mapModel.removedKeys = mapset.NewThreadUnsafeSet()
	mapModel.valueFactory = valueFactory
	mapModel.data = make(map[int]IntObjectMapValueModel)
	return mapModel
}

type StringObjectMapValueModel interface {
	MapValueModel
	setKey(key string)
	Key() string
}

type StringObjectMapModel interface {
	mapModel
	Keys() []string
	Get(key string) StringObjectMapValueModel
	Put(key string, value StringObjectMapValueModel) StringObjectMapValueModel
	Remove(key string) bool
}

type StringObjectMapValueFactory func() StringObjectMapValueModel

type stringObjectMap struct {
	baseMap
	valueFactory StringObjectMapValueFactory
	data         map[string]StringObjectMapValueModel
}

func (imap *stringObjectMap) Size() int {
	return len(imap.data)
}

func (imap *stringObjectMap) Clear() {
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

func (imap *stringObjectMap) Keys() []string {
	data := imap.data
	keys := make([]string, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	return keys
}

func (imap *stringObjectMap) Get(key string) StringObjectMapValueModel {
	value, ok := imap.data[key]
	if ok {
		return value
	}
	return nil
}

func (imap *stringObjectMap) Put(key string, value StringObjectMapValueModel) StringObjectMapValueModel {
	data := imap.data
	old, ok := data[key]
	if ok {
		if old != value {
			data[key] = value
			imap.updatedKeys.Add(key)
			value.setKey(key)
			value.setParent(imap)
			value.SetFullyUpdate(true)
			old.unbind()
		}
		return old
	}
	data[key] = value
	imap.updatedKeys.Add(key)
	imap.removedKeys.Remove(key)
	value.setKey(key)
	value.setParent(imap)
	value.SetFullyUpdate(true)
	return nil
}

func (imap *stringObjectMap) Remove(key string) bool {
	data := imap.data
	old, ok := data[key]
	if ok {
		delete(data, key)
		imap.updatedKeys.Remove(key)
		imap.removedKeys.Add(key)
		old.unbind()
		return true
	}
	return false
}

func (imap *stringObjectMap) ToBson() interface{} {
	bson := bson.M{}
	for key, value := range imap.data {
		bson[key] = value.ToBson()
	}
	return bson
}

func (imap *stringObjectMap) ToData() interface{} {
	data := make(map[string]interface{})
	for key, value := range imap.data {
		data[key] = value.ToData()
	}
	return data
}

func (imap *stringObjectMap) Reset() {
	imap.updatedKeys.Clear()
	imap.removedKeys.Clear()
	for _, v := range imap.data {
		v.Reset()
	}
}

func (imap *stringObjectMap) LoadJsoniter(any jsoniter.Any) error {
	imap.Reset()
	data := imap.data
	for k, v := range data {
		v.unbind()
		delete(data, k)
	}
	if any.ValueType() == jsoniter.ObjectValue {
		valueFactory := imap.valueFactory
		keys := any.Keys()
		for _, key := range keys {
			value := valueFactory()
			err := value.LoadJsoniter(any.Get(key))
			if err != nil {
				return err
			}
			data[key] = value
			value.setParent(imap)
			value.setKey(key)
		}
	}
	return nil
}

func (imap *stringObjectMap) AppendUpdates(updates bson.M) bson.M {
	data := imap.data
	updatedKeys := imap.updatedKeys
	if updatedKeys.Cardinality() > 0 {
		dset, ok := updates["$set"].(bson.M)
		if !ok {
			dset = bson.M{}
		}
		for _, uk := range updatedKeys.ToSlice() {
			key := uk.(string)
			value := data[key]
			if value.FullyUpdate() {
				dset[value.XPath().Value()] = value.ToBson()
			} else {
				value.AppendUpdates(updates)
			}
		}
	}
	removedKeys := imap.removedKeys
	if removedKeys.Cardinality() > 0 {
		unset, ok := updates["$unset"].(bson.M)
		if !ok {
			unset = bson.M{}
		}
		for _, uk := range removedKeys.ToSlice() {
			key := uk.(string)
			name := imap.XPath().Resolve(key)
			unset[name.Value()] = 1
		}
	}
	return updates
}

func (imap *stringObjectMap) ToDocument() bson.M {
	doc := bson.M{}
	for k, v := range imap.data {
		value := v.ToBson()
		doc[k] = value
	}
	return doc
}

func (imap *stringObjectMap) LoadDocument(document bson.M) error {
	imap.Reset()
	data := imap.data
	for k := range data {
		delete(data, k)
	}
	valueFactory := imap.valueFactory
	for key, value := range document {
		obj, ok := value.(bson.M)
		if !ok {
			// skip value that not be a bson.M
			continue
		}
		v := valueFactory()
		err := v.LoadDocument(obj)
		if err != nil {
			return err
		}
		data[key] = v
	}
	return nil
}

func NewStringObjectMapModel(parent BsonModel, name string, valueFactory StringObjectMapValueFactory) StringObjectMapModel {
	mapModel := &stringObjectMap{}
	mapModel.parent = parent
	mapModel.name = name
	mapModel.updatedKeys = mapset.NewThreadUnsafeSet()
	mapModel.removedKeys = mapset.NewThreadUnsafeSet()
	mapModel.valueFactory = valueFactory
	mapModel.data = make(map[string]StringObjectMapValueModel)
	return mapModel
}
