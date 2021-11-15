package example

import (
	"unsafe"

	"github.com/bits-and-blooms/bitset"
	"github.com/fmjsjx/bson-model-go/bsonmodel"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
)

type Equipment interface {
	bsonmodel.StringObjectMapValueModel
	Id() string
	SetId(id string)
	RefId() int
	SetRefId(refId int)
	Atk() int
	SetAtk(atk int)
	Def() int
	SetDef(def int)
	Hp() int
	SetHp(hp int)
}

const (
	BnameEquipmentId    = "id"
	BnameEquipmentRefId = "rid"
	BnameEquipmentAtk   = "atk"
	BnameEquipmentDef   = "def"
	BnameEquipmentHp    = "hp"
)

type defaultEquipment struct {
	bsonmodel.BaseStringObjectMapValue
	updatedFields *bitset.BitSet
	id            string
	refId         int
	atk           int
	def           int
	hp            int
}

func (self *defaultEquipment) ToBson() interface{} {
	return self.ToDocument()
}

func (self *defaultEquipment) ToData() interface{} {
	data := make(map[string]interface{})
	data["id"] = self.id
	data["rid"] = self.refId
	data["atk"] = self.atk
	data["def"] = self.def
	data["hp"] = self.hp
	return data
}

func (self *defaultEquipment) LoadJsoniter(any jsoniter.Any) error {
	if any.ValueType() != jsoniter.ObjectValue {
		return nil
	}
	id, err := bsonmodel.AnyStringValue(any.Get("id"), "")
	if err != nil {
		return err
	}
	self.id = id
	refId, err := bsonmodel.AnyIntValue(any.Get("rid"), 0)
	if err != nil {
		return err
	}
	self.refId = refId
	atk, err := bsonmodel.AnyIntValue(any.Get("atk"), 0)
	if err != nil {
		return err
	}
	self.atk = atk
	def, err := bsonmodel.AnyIntValue(any.Get("def"), 0)
	if err != nil {
		return err
	}
	self.def = def
	hp, err := bsonmodel.AnyIntValue(any.Get("hp"), 0)
	if err != nil {
		return err
	}
	self.hp = hp
	return nil
}

func (self *defaultEquipment) Reset() {
	self.updatedFields.ClearAll()
}

func (self *defaultEquipment) AnyUpdated() bool {
	return self.updatedFields.Any()
}

func (self *defaultEquipment) AnyDeleted() bool {
	return self.DeletedSize() > 0
}

func (self *defaultEquipment) AppendUpdates(updates bson.M) bson.M {
	dset := bsonmodel.FixedEmbedded(updates, "$set")
	xpath := self.XPath()
	if self.FullyUpdate() {
		dset[xpath.Value()] = self.ToDocument()
	} else {
		updatedFields := self.updatedFields
		if updatedFields.Test(1) {
			dset[xpath.Resolve("id").Value()] = self.id
		}
		if updatedFields.Test(2) {
			dset[xpath.Resolve("rid").Value()] = self.refId
		}
		if updatedFields.Test(3) {
			dset[xpath.Resolve("atk").Value()] = self.atk
		}
		if updatedFields.Test(4) {
			dset[xpath.Resolve("def").Value()] = self.def
		}
		if updatedFields.Test(5) {
			dset[xpath.Resolve("hp").Value()] = self.hp
		}
	}
	return updates
}

func (self *defaultEquipment) ToDocument() bson.M {
	doc := bson.M{}
	doc["id"] = self.id
	doc["rid"] = self.refId
	doc["atk"] = self.atk
	doc["def"] = self.def
	doc["hp"] = self.hp
	return doc
}

func (self *defaultEquipment) LoadDocument(document bson.M) error {
	id, err := bsonmodel.StringValue(document, "id", "")
	if err != nil {
		return err
	}
	self.id = id
	refId, err := bsonmodel.IntValue(document, "rid", 0)
	if err != nil {
		return err
	}
	self.refId = refId
	atk, err := bsonmodel.IntValue(document, "atk", 0)
	if err != nil {
		return err
	}
	self.atk = atk
	def, err := bsonmodel.IntValue(document, "def", 0)
	if err != nil {
		return err
	}
	self.def = def
	hp, err := bsonmodel.IntValue(document, "hp", 0)
	if err != nil {
		return err
	}
	self.hp = hp
	return nil
}

func (self *defaultEquipment) DeletedSize() int {
	return 0
}

func (self *defaultEquipment) FullyUpdate() bool {
	return self.updatedFields.Test(0)
}

func (self *defaultEquipment) SetFullyUpdate(fullyUpdate bool) {
	if fullyUpdate {
		self.updatedFields.Set(0)
	} else {
		self.updatedFields.DeleteAt(0)
	}
}

func (self *defaultEquipment) ToSync() interface{} {
	if self.FullyUpdate() {
		return self
	}
	sync := make(map[string]interface{})
	updatedFields := self.updatedFields
	if updatedFields.Test(1) {
		sync["id"] = self.id
	}
	if updatedFields.Test(2) {
		sync["refId"] = self.refId
	}
	if updatedFields.Test(3) {
		sync["atk"] = self.atk
	}
	if updatedFields.Test(4) {
		sync["def"] = self.def
	}
	if updatedFields.Test(5) {
		sync["hp"] = self.hp
	}
	return sync
}

func (self *defaultEquipment) ToDelete() interface{} {
	delete := make(map[string]interface{})
	return delete
}

func (self *defaultEquipment) ToDataJson() (string, error) {
	return jsoniter.MarshalToString(self.ToData())
}

func (self *defaultEquipment) ToSyncJson() (string, error) {
	return jsoniter.MarshalToString(self.ToSync())
}

func (self *defaultEquipment) ToDeleteJson() (string, error) {
	return jsoniter.MarshalToString(self.ToDelete())
}

func (self *defaultEquipment) Id() string {
	return self.id
}

func (self *defaultEquipment) SetId(id string) {
	if self.id != id {
		self.id = id
		self.updatedFields.Set(1)
	}
}

func (self *defaultEquipment) RefId() int {
	return self.refId
}

func (self *defaultEquipment) SetRefId(refId int) {
	if self.refId != refId {
		self.refId = refId
		self.updatedFields.Set(2)
	}
}

func (self *defaultEquipment) Atk() int {
	return self.atk
}

func (self *defaultEquipment) SetAtk(atk int) {
	if self.atk != atk {
		self.atk = atk
		self.updatedFields.Set(3)
	}
}

func (self *defaultEquipment) Def() int {
	return self.def
}

func (self *defaultEquipment) SetDef(def int) {
	if self.def != def {
		self.def = def
		self.updatedFields.Set(4)
	}
}

func (self *defaultEquipment) Hp() int {
	return self.hp
}

func (self *defaultEquipment) SetHp(hp int) {
	if self.hp != hp {
		self.hp = hp
		self.updatedFields.Set(5)
	}
}

func NewEquipment() Equipment {
	self := &defaultEquipment{updatedFields: &bitset.BitSet{}}
	return self
}

var equipmentFactory bsonmodel.StringObjectMapValueFactory = func() bsonmodel.StringObjectMapValueModel {
	return NewEquipment()
}

func EquipmentFactory() bsonmodel.StringObjectMapValueFactory {
	return equipmentFactory
}

type equipmentEncoder struct{}

func (codec *equipmentEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

func (codec *equipmentEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	p := ((*defaultEquipment)(ptr))
	stream.WriteObjectStart()
	stream.WriteObjectField("id")
	stream.WriteString(p.id)
	stream.WriteMore()
	stream.WriteObjectField("refId")
	stream.WriteInt(p.refId)
	stream.WriteMore()
	stream.WriteObjectField("atk")
	stream.WriteInt(p.atk)
	stream.WriteMore()
	stream.WriteObjectField("def")
	stream.WriteInt(p.def)
	stream.WriteMore()
	stream.WriteObjectField("hp")
	stream.WriteInt(p.hp)
	stream.WriteObjectEnd()
}

func init() {
	jsoniter.RegisterTypeEncoder("example.defaultEquipment", &equipmentEncoder{})
}

