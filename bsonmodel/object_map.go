package bsonmodel

import (
	"strconv"
	"unsafe"

	mapset "github.com/deckarep/golang-set"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
)

type IntObjectMapValueModel interface {
	MapValueModel
	setKey(key int)
	Key() int
}

type BaseIntObjectMapValue struct {
	baseMapValue
	key int
}

func (v *BaseIntObjectMapValue) Parent() BsonModel {
	return v.parent
}

func (v *BaseIntObjectMapValue) XPath() DotNotation {
	return v.parent.XPath().Resolve(strconv.Itoa(v.key))
}

func (v *BaseIntObjectMapValue) setParent(parent BsonModel) {
	v.parent = parent
}

func (v *BaseIntObjectMapValue) unbind() {
	v.parent = nil
	v.key = 0
}

func (v *BaseIntObjectMapValue) setKey(key int) {
	v.key = key
}

func (v *BaseIntObjectMapValue) Key() int {
	return v.key
}

func (v *BaseIntObjectMapValue) EmitUpdated() {
	if v.parent != nil {
		v.parent.(mapModel).emitUpdated(v.key)
	}
}

type IntObjectMapModel interface {
	mapModel
	Keys() []int
	Get(key int) IntObjectMapValueModel
	Put(key int, value IntObjectMapValueModel) IntObjectMapValueModel
	Remove(key int) bool
	SetUpdated(key string)
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

func (imap *intObjectMap) SetUpdated(key string) {
	imap.updatedKeys.Add(key)
}

func (imap *intObjectMap) ToBson() interface{} {
	return imap.ToDocument()
}

func (imap *intObjectMap) ToData() interface{} {
	data := make(map[int]interface{})
	for key, value := range imap.data {
		data[key] = value.ToData()
	}
	return data
}

func (imap *intObjectMap) Reset() {
	data := imap.data
	for _, k := range imap.updatedKeys.ToSlice() {
		data[k.(int)].Reset()
	}
	imap.updatedKeys.Clear()
	imap.removedKeys.Clear()
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
		dset := FixedEmbedded(updates, "$set")
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
		unset := FixedEmbedded(updates, "$unset")
		for _, uk := range removedKeys.ToSlice() {
			key := uk.(int)
			k := strconv.Itoa(key)
			name := imap.XPath().Resolve(k)
			unset[name.Value()] = ""
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
		v.setParent(imap)
		v.setKey(k)
	}
	return nil
}

func (imap *intObjectMap) ToSync() interface{} {
	sync := make(map[int]interface{})
	data := imap.data
	updatedKeys := imap.updatedKeys
	if updatedKeys.Cardinality() > 0 {
		for _, uk := range updatedKeys.ToSlice() {
			key := uk.(int)
			value := data[key]
			sync[key] = value.ToSync()
		}
	}
	return sync
}

func (imap *intObjectMap) ToDelete() interface{} {
	delete := make(map[int]int)
	removedKeys := imap.removedKeys
	if removedKeys.Cardinality() > 0 {
		for _, uk := range removedKeys.ToSlice() {
			key := uk.(int)
			delete[key] = 1
		}
	}
	return delete
}

func (imap *intObjectMap) ToDataJson() (string, error) {
	return jsoniter.MarshalToString(imap.ToData())
}

func (imap *intObjectMap) ToSyncJson() (string, error) {
	return jsoniter.MarshalToString(imap.ToSync())
}

func (imap *intObjectMap) ToDeleteJson() (string, error) {
	return jsoniter.MarshalToString(imap.ToDelete())
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
	Key() string
	setKey(key string)
}

type BaseStringObjectMapValue struct {
	baseMapValue
	key string
}

func (v *BaseStringObjectMapValue) Parent() BsonModel {
	return v.parent
}

func (v *BaseStringObjectMapValue) XPath() DotNotation {
	return v.parent.XPath().Resolve(v.key)
}

func (v *BaseStringObjectMapValue) setParent(parent BsonModel) {
	v.parent = parent
}

func (v *BaseStringObjectMapValue) unbind() {
	v.parent = nil
	v.key = ""
}

func (v *BaseStringObjectMapValue) setKey(key string) {
	v.key = key
}

func (v *BaseStringObjectMapValue) Key() string {
	return v.key
}

func (v *BaseStringObjectMapValue) EmitUpdated() {
	if v.parent != nil {
		v.parent.(mapModel).emitUpdated(v.key)
	}
}

type StringObjectMapModel interface {
	mapModel
	Keys() []string
	Get(key string) StringObjectMapValueModel
	Put(key string, value StringObjectMapValueModel) StringObjectMapValueModel
	Remove(key string) bool
	SetUpdated(key string)
}

type StringObjectMapValueFactory func() StringObjectMapValueModel

type stringObjectMap struct {
	baseMap
	valueFactory StringObjectMapValueFactory
	data         map[string]StringObjectMapValueModel
}

func (smap *stringObjectMap) Size() int {
	return len(smap.data)
}

func (smap *stringObjectMap) Clear() {
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

func (smap *stringObjectMap) Keys() []string {
	data := smap.data
	keys := make([]string, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	return keys
}

func (smap *stringObjectMap) Get(key string) StringObjectMapValueModel {
	value, ok := smap.data[key]
	if ok {
		return value
	}
	return nil
}

func (smap *stringObjectMap) Put(key string, value StringObjectMapValueModel) StringObjectMapValueModel {
	data := smap.data
	old, ok := data[key]
	if ok {
		if old != value {
			data[key] = value
			smap.updatedKeys.Add(key)
			value.setKey(key)
			value.setParent(smap)
			value.SetFullyUpdate(true)
			old.unbind()
		}
		return old
	}
	data[key] = value
	smap.updatedKeys.Add(key)
	smap.removedKeys.Remove(key)
	value.setKey(key)
	value.setParent(smap)
	value.SetFullyUpdate(true)
	return nil
}

func (smap *stringObjectMap) Remove(key string) bool {
	data := smap.data
	old, ok := data[key]
	if ok {
		delete(data, key)
		smap.updatedKeys.Remove(key)
		smap.removedKeys.Add(key)
		old.unbind()
		return true
	}
	return false
}

func (smap *stringObjectMap) SetUpdated(key string) {
	smap.updatedKeys.Add(key)
}

func (smap *stringObjectMap) ToBson() interface{} {
	return smap.ToDocument()
}

func (smap *stringObjectMap) ToData() interface{} {
	data := make(map[string]interface{})
	for key, value := range smap.data {
		data[key] = value.ToData()
	}
	return data
}

func (smap *stringObjectMap) Reset() {
	data := smap.data
	for _, k := range smap.updatedKeys.ToSlice() {
		data[k.(string)].Reset()
	}
	smap.updatedKeys.Clear()
	smap.removedKeys.Clear()
}

func (smap *stringObjectMap) LoadJsoniter(any jsoniter.Any) error {
	smap.Reset()
	data := smap.data
	for k, v := range data {
		v.unbind()
		delete(data, k)
	}
	if any.ValueType() == jsoniter.ObjectValue {
		valueFactory := smap.valueFactory
		keys := any.Keys()
		for _, key := range keys {
			value := valueFactory()
			err := value.LoadJsoniter(any.Get(key))
			if err != nil {
				return err
			}
			data[key] = value
			value.setParent(smap)
			value.setKey(key)
		}
	}
	return nil
}

func (smap *stringObjectMap) AppendUpdates(updates bson.M) bson.M {
	data := smap.data
	updatedKeys := smap.updatedKeys
	if updatedKeys.Cardinality() > 0 {
		dset := FixedEmbedded(updates, "$set")
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
	removedKeys := smap.removedKeys
	if removedKeys.Cardinality() > 0 {
		unset := FixedEmbedded(updates, "$unset")
		for _, uk := range removedKeys.ToSlice() {
			key := uk.(string)
			name := smap.XPath().Resolve(key)
			unset[name.Value()] = ""
		}
	}
	return updates
}

func (smap *stringObjectMap) ToDocument() bson.M {
	doc := bson.M{}
	for k, v := range smap.data {
		value := v.ToBson()
		doc[k] = value
	}
	return doc
}

func (smap *stringObjectMap) LoadDocument(document bson.M) error {
	smap.Reset()
	data := smap.data
	for k := range data {
		delete(data, k)
	}
	valueFactory := smap.valueFactory
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
		v.setParent(smap)
		v.setKey(key)
	}
	return nil
}

func (smap *stringObjectMap) ToSync() interface{} {
	sync := make(map[string]interface{})
	data := smap.data
	updatedKeys := smap.updatedKeys
	if updatedKeys.Cardinality() > 0 {
		for _, uk := range updatedKeys.ToSlice() {
			key := uk.(string)
			value := data[key]
			sync[key] = value.ToSync()
		}
	}
	return sync
}

func (smap *stringObjectMap) ToDelete() interface{} {
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

func (smap *stringObjectMap) ToDataJson() (string, error) {
	return jsoniter.MarshalToString(smap.ToData())
}

func (smap *stringObjectMap) ToSyncJson() (string, error) {
	return jsoniter.MarshalToString(smap.ToSync())
}

func (smap *stringObjectMap) ToDeleteJson() (string, error) {
	return jsoniter.MarshalToString(smap.ToDelete())
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

type intObjectMapEncoder struct{}

func (codec *intObjectMapEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	imap := ((*intObjectMap)(ptr))
	return len(imap.data) == 0
}

func (codec *intObjectMapEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	imap := ((*intObjectMap)(ptr))
	stream.WriteVal(imap.data)
}

type stringObjectMapEncoder struct{}

func (codec *stringObjectMapEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	smap := ((*stringObjectMap)(ptr))
	return len(smap.data) == 0
}

func (codec *stringObjectMapEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	smap := ((*stringObjectMap)(ptr))
	stream.WriteVal(smap.data)
}

func init() {
	jsoniter.RegisterTypeEncoder("bsonmodel.intObjectMap", &intObjectMapEncoder{})
	jsoniter.RegisterTypeEncoder("bsonmodel.stringObjectMap", &stringObjectMapEncoder{})
}
