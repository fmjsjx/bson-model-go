package bsonmodel

import (
	mapset "github.com/deckarep/golang-set"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
)

// The top interface for BSON model.
type BsonModel interface {
	ToBson() interface{}
	ToData() interface{}
	LoadJsoniter(any jsoniter.Any) error
	Reset()
	AnyUpdated() bool
	AnyDeleted() bool
	Parent() BsonModel
	XPath() DotNotation
	AppendUpdates(updates bson.M) bson.M
}

type DocumentModel interface {
	BsonModel
	ToDocument() bson.M
	LoadDocument(document bson.M) error
	DeletedSize() int
}

type ObjectModel interface {
	DocumentModel
	FullyUpdate() bool
	SetFullyUpdate(fullyUpdate bool)
}

type RootModel interface {
	ObjectModel
	ToUpdate() bson.M
	MarshalToJsonString() (string, error)
}

type MapValueModel interface {
	ObjectModel
	setParent(parent BsonModel)
	unbind()
}

type baseMapValue struct {
	parent BsonModel
}

func (v *baseMapValue) setParent(parent BsonModel) {
	v.parent = parent
}

type mapModel interface {
	DocumentModel
	Size() int
	Clear()
}

type baseMap struct {
	parent      BsonModel
	name        string
	updatedKeys mapset.Set
	removedKeys mapset.Set
}

func (smap *baseMap) AnyUpdated() bool {
	return smap.updatedKeys.Cardinality() > 0 || smap.AnyDeleted()
}

func (smap *baseMap) AnyDeleted() bool {
	return smap.DeletedSize() > 0
}

func (smap *baseMap) Parent() BsonModel {
	return smap.parent
}

func (smap *baseMap) XPath() DotNotation {
	return smap.parent.XPath().Resolve(smap.name)
}

func (smap *baseMap) DeletedSize() int {
	return smap.removedKeys.Cardinality()
}
