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

func (e *defaultEquipment) ToBson() interface{} {
	return e.ToDocument()
}

func (e *defaultEquipment) ToData() interface{} {
	data := make(map[string]interface{})
	data["id"] = e.id
	data["rid"] = e.refId
	data["atk"] = e.atk
	data["def"] = e.def
	data["hp"] = e.hp
	return data
}

func (e *defaultEquipment) LoadJsoniter(any jsoniter.Any) error {
	if any.ValueType() != jsoniter.ObjectValue {
		return nil
	}
	id, err := bsonmodel.AnyStringValue(any.Get("id"), "")
	if err != nil {
		return err
	}
	e.id = id
	refId, err := bsonmodel.AnyIntValue(any.Get("rid"), 0)
	if err != nil {
		return err
	}
	e.refId = refId
	atk, err := bsonmodel.AnyIntValue(any.Get("atk"), 0)
	if err != nil {
		return err
	}
	e.atk = atk
	def, err := bsonmodel.AnyIntValue(any.Get("def"), 0)
	if err != nil {
		return err
	}
	e.def = def
	hp, err := bsonmodel.AnyIntValue(any.Get("hp"), 0)
	if err != nil {
		return err
	}
	e.hp = hp
	return nil
}

func (e *defaultEquipment) Reset() {
	e.updatedFields.ClearAll()
}

func (e *defaultEquipment) AnyUpdated() bool {
	return e.updatedFields.Any()
}

func (e *defaultEquipment) AnyDeleted() bool {
	return e.DeletedSize() > 0
}

func (e *defaultEquipment) AppendUpdates(updates bson.M) bson.M {
	dset := bsonmodel.FixedEmbedded(updates, "$set")
	xpath := e.XPath()
	if e.FullyUpdate() {
		dset[xpath.Value()] = e.ToDocument()
	} else {
		updatedFields := e.updatedFields
		if updatedFields.Test(1) {
			dset[xpath.Resolve("id").Value()] = e.id
		}
		if updatedFields.Test(2) {
			dset[xpath.Resolve("rid").Value()] = e.refId
		}
		if updatedFields.Test(3) {
			dset[xpath.Resolve("atk").Value()] = e.atk
		}
		if updatedFields.Test(4) {
			dset[xpath.Resolve("def").Value()] = e.def
		}
		if updatedFields.Test(5) {
			dset[xpath.Resolve("hp").Value()] = e.hp
		}
	}
	return updates
}

func (e *defaultEquipment) ToDocument() bson.M {
	doc := bson.M{}
	doc["id"] = e.id
	doc["rid"] = e.refId
	doc["atk"] = e.atk
	doc["def"] = e.def
	doc["hp"] = e.hp
	return doc
}

func (e *defaultEquipment) LoadDocument(document bson.M) error {
	id, err := bsonmodel.StringValue(document, "id", "")
	if err != nil {
		return err
	}
	e.id = id
	refId, err := bsonmodel.IntValue(document, "rid", 0)
	if err != nil {
		return err
	}
	e.refId = refId
	atk, err := bsonmodel.IntValue(document, "atk", 0)
	if err != nil {
		return err
	}
	e.atk = atk
	def, err := bsonmodel.IntValue(document, "def", 0)
	if err != nil {
		return err
	}
	e.def = def
	hp, err := bsonmodel.IntValue(document, "hp", 0)
	if err != nil {
		return err
	}
	e.hp = hp
	return nil
}

func (e *defaultEquipment) DeletedSize() int {
	return 0
}

func (e *defaultEquipment) FullyUpdate() bool {
	return e.updatedFields.Test(0)
}

func (e *defaultEquipment) SetFullyUpdate(fullyUpdate bool) {
	if fullyUpdate {
		e.updatedFields.Set(0)
	} else {
		e.updatedFields.DeleteAt(0)
	}
}

func (e *defaultEquipment) Id() string {
	return e.id
}

func (e *defaultEquipment) SetId(id string) {
	if e.id != id {
		e.id = id
		e.updatedFields.Set(1)
		p := e.Parent()
		if p != nil {
			p.(bsonmodel.StringObjectMapModel).SetUpdated(e.Key())
		}
	}
}

func (e *defaultEquipment) RefId() int {
	return e.refId
}

func (e *defaultEquipment) SetRefId(refId int) {
	if e.refId != refId {
		e.refId = refId
		e.updatedFields.Set(2)
		p := e.Parent()
		if p != nil {
			p.(bsonmodel.StringObjectMapModel).SetUpdated(e.Key())
		}
	}
}

func (e *defaultEquipment) Atk() int {
	return e.atk
}

func (e *defaultEquipment) SetAtk(atk int) {
	if e.atk != atk {
		e.atk = atk
		e.updatedFields.Set(3)
		p := e.Parent()
		if p != nil {
			p.(bsonmodel.StringObjectMapModel).SetUpdated(e.Key())
		}
	}
}

func (e *defaultEquipment) Def() int {
	return e.def
}

func (e *defaultEquipment) SetDef(def int) {
	if e.def != def {
		e.def = def
		e.updatedFields.Set(4)
		p := e.Parent()
		if p != nil {
			p.(bsonmodel.StringObjectMapModel).SetUpdated(e.Key())
		}
	}
}

func (e *defaultEquipment) Hp() int {
	return e.hp
}

func (e *defaultEquipment) SetHp(hp int) {
	if e.hp != hp {
		e.hp = hp
		e.updatedFields.Set(5)
		p := e.Parent()
		if p != nil {
			p.(bsonmodel.StringObjectMapModel).SetUpdated(e.Key())
		}
	}
}

func NewEquipment() Equipment {
	equipment := &defaultEquipment{updatedFields: &bitset.BitSet{}}
	return equipment
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
	e := ((*defaultEquipment)(ptr))
	stream.WriteObjectStart()
	stream.WriteObjectField("id")
	stream.WriteString(e.id)
	stream.WriteMore()
	stream.WriteObjectField("refId")
	stream.WriteInt(e.refId)
	stream.WriteMore()
	stream.WriteObjectField("atk")
	stream.WriteInt(e.atk)
	stream.WriteMore()
	stream.WriteObjectField("def")
	stream.WriteInt(e.def)
	stream.WriteMore()
	stream.WriteObjectField("hp")
	stream.WriteInt(e.hp)
	stream.WriteObjectEnd()
}

func init() {
	jsoniter.RegisterTypeEncoder("example.defaultEquipment", &equipmentEncoder{})
}
