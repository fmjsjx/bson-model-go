package example

import (
	"time"
	"unsafe"

	"github.com/bits-and-blooms/bitset"
	"github.com/fmjsjx/bson-model-go/bsonmodel"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	BnamePlayerUid           = "_id"
	BnamePlayerWallet        = "wt"
	BnamePlayerEquipments    = "eqm"
	BnamePlayerItems         = "itm"
	BnamePlayerUpdateVersion = "_uv"
	BnamePlayerCreateTime    = "_ct"
	BnamePlayerUpdateTime    = "_ut"
)

type Player interface {
	bsonmodel.RootModel
	Uid() int
	SetUid(uid int)
	Wallet() Wallet
	Equipments() bsonmodel.StringObjectMapModel
	Equipment(id string) Equipment
	Items() bsonmodel.IntSimpleMapModel
	UpdateVersion() int
	SetUpdateVersion(updateVersion int)
	IncreaseUpdateVersion() int
	CreateTime() time.Time
	SetCreateTime(createTime time.Time)
	UpdateTime() time.Time
	SetUpdateTime(updateTime time.Time)
}

type defaultPlayer struct {
	updatedFields *bitset.BitSet
	uid           int
	wallet        Wallet
	equipments    bsonmodel.StringObjectMapModel
	items         bsonmodel.IntSimpleMapModel
	updateVersion int
	createTime    time.Time
	updateTime    time.Time
}

func (e *defaultPlayer) ToBson() interface{} {
	return e.ToDocument()
}

func (p *defaultPlayer) ToData() interface{} {
	data := make(map[string]interface{})
	data["_id"] = p.uid
	data["wt"] = p.wallet.ToData()
	data["eqm"] = p.equipments.ToData()
	data["itm"] = p.items.ToData()
	data["_uv"] = p.updateVersion
	data["_ct"] = p.createTime.UnixMilli()
	data["_ut"] = p.updateTime.UnixMilli()
	return data
}

func (p *defaultPlayer) LoadJsoniter(any jsoniter.Any) error {
	if any.ValueType() != jsoniter.ObjectValue {
		p.Reset()
		return nil
	}
	uid, err := bsonmodel.AnyIntValue(any.Get("_id"), 0)
	if err != nil {
		return err
	}
	p.uid = uid
	wallet := any.Get("wt")
	if wallet.ValueType() == jsoniter.ObjectValue {
		err = p.wallet.LoadJsoniter(wallet)
		if err != nil {
			return err
		}
	}
	equipments := any.Get("eqm")
	if equipments.ValueType() == jsoniter.ObjectValue {
		err = p.equipments.LoadJsoniter(equipments)
		if err != nil {
			return err
		}
	} else {
		p.equipments.Clear()
	}
	items := any.Get("itm")
	if items.ValueType() == jsoniter.ObjectValue {
		err = p.items.LoadJsoniter(items)
		if err != nil {
			return err
		}
	} else {
		p.items.Clear()
	}
	updateVersion, err := bsonmodel.AnyIntValue(any.Get("_uv"), 0)
	if err != nil {
		return err
	}
	p.updateVersion = updateVersion
	createTime, err := bsonmodel.AnyDateTimeValue(any.Get("_ct"))
	if err != nil {
		return err
	}
	p.createTime = createTime
	updateTime, err := bsonmodel.AnyDateTimeValue(any.Get("_ut"))
	if err != nil {
		return err
	}
	p.updateTime = updateTime
	p.Reset()
	return nil
}

func (p *defaultPlayer) Reset() {
	p.wallet.Reset()
	p.equipments.Reset()
	p.items.Reset()
	p.updatedFields.ClearAll()
}

func (p *defaultPlayer) AnyUpdated() bool {
	return p.updatedFields.Any() || p.wallet.AnyUpdated() || p.equipments.AnyUpdated() || p.items.AnyUpdated()
}

func (p *defaultPlayer) AnyDeleted() bool {
	return p.DeletedSize() > 0
}

func (p *defaultPlayer) Parent() bsonmodel.BsonModel {
	return nil
}

func (p *defaultPlayer) XPath() bsonmodel.DotNotation {
	return bsonmodel.RootPath()
}

func (p *defaultPlayer) AppendUpdates(updates bson.M) bson.M {
	dset := bsonmodel.FixedEmbedded(updates, "$set")
	updatedFields := p.updatedFields
	if updatedFields.Test(1) {
		dset["id"] = p.uid
	}
	if p.wallet.AnyUpdated() {
		p.wallet.AppendUpdates(updates)
	}
	if p.equipments.AnyUpdated() {
		p.equipments.AppendUpdates(updates)
	}
	if p.items.AnyUpdated() {
		p.items.AppendUpdates(updates)
	}
	if updatedFields.Test(5) {
		dset["_uv"] = p.updateVersion
	}
	if updatedFields.Test(6) {
		dset["_ct"] = primitive.NewDateTimeFromTime(p.createTime)
	}
	if updatedFields.Test(7) {
		dset["_ut"] = primitive.NewDateTimeFromTime(p.updateTime)
	}
	return updates
}

func (p *defaultPlayer) ToDocument() bson.M {
	doc := bson.M{}
	doc["_id"] = p.uid
	doc["wt"] = p.wallet.ToBson()
	doc["eqm"] = p.equipments.ToBson()
	doc["itm"] = p.items.ToBson()
	doc["_uv"] = p.updateVersion
	doc["_ct"] = primitive.NewDateTimeFromTime(p.createTime)
	doc["_ut"] = primitive.NewDateTimeFromTime(p.updateTime)
	return doc
}

func (p *defaultPlayer) LoadDocument(document bson.M) error {
	uid, err := bsonmodel.IntValue(document, "_id", 0)
	if err != nil {
		return err
	}
	p.uid = uid
	wallet, err := bsonmodel.EmbeddedValue(document, "wt")
	if err != nil {
		return err
	}
	if wallet != nil {
		err = p.wallet.LoadDocument(wallet)
		if err != nil {
			return err
		}
	}
	equipments, err := bsonmodel.EmbeddedValue(document, "eqm")
	if err != nil {
		return err
	}
	if equipments != nil {
		err = p.equipments.LoadDocument(equipments)
		if err != nil {
			return err
		}
	} else {
		p.equipments.Clear()
	}
	items, err := bsonmodel.EmbeddedValue(document, "itm")
	if err != nil {
		return err
	}
	if items != nil {
		err = p.items.LoadDocument(items)
		if err != nil {
			return err
		}
	} else {
		p.items.Clear()
	}
	updateVersion, err := bsonmodel.IntValue(document, "_uv", 0)
	if err != nil {
		return err
	}
	p.updateVersion = updateVersion
	createTime, err := bsonmodel.DateTimeValue(document, "_ct")
	if err != nil {
		return err
	}
	p.createTime = createTime
	updateTime, err := bsonmodel.DateTimeValue(document, "_ut")
	if err != nil {
		return err
	}
	p.updateTime = updateTime
	p.Reset()
	return nil
}

func (p *defaultPlayer) DeletedSize() int {
	n := 0
	if p.wallet.AnyDeleted() {
		n += 1
	}
	if p.equipments.AnyDeleted() {
		n += 1
	}
	if p.items.AnyDeleted() {
		n += 1
	}
	return n
}

func (p *defaultPlayer) FullyUpdate() bool {
	return false
}

func (p *defaultPlayer) SetFullyUpdate(fullyUpdate bool) {
	// no effect
}

func (p *defaultPlayer) ToUpdate() bson.M {
	if p.AnyUpdated() {
		return p.AppendUpdates(bson.M{})
	}
	return bson.M{}
}

func (p *defaultPlayer) MarshalToJsonString() (string, error) {
	return jsoniter.MarshalToString(p)
}

func (p *defaultPlayer) Uid() int {
	return p.uid
}

func (p *defaultPlayer) SetUid(uid int) {
	if p.uid != uid {
		p.uid = uid
		p.updatedFields.Set(1)
	}
}

func (p *defaultPlayer) Wallet() Wallet {
	return p.wallet
}

func (p *defaultPlayer) Equipments() bsonmodel.StringObjectMapModel {
	return p.equipments
}

func (p *defaultPlayer) Equipment(id string) Equipment {
	value := p.equipments.Get(id)
	if value == nil {
		return nil
	}
	return value.(Equipment)
}

func (p *defaultPlayer) Items() bsonmodel.IntSimpleMapModel {
	return p.items
}

func (p *defaultPlayer) UpdateVersion() int {
	return p.updateVersion
}

func (p *defaultPlayer) SetUpdateVersion(updateVersion int) {
	if p.updateVersion != updateVersion {
		p.updateVersion = updateVersion
		p.updatedFields.Set(5)
	}
}

func (p *defaultPlayer) IncreaseUpdateVersion() int {
	updateVersion := p.updateVersion + 1
	p.updateVersion = updateVersion
	p.updatedFields.Set(5)
	return updateVersion
}

func (p *defaultPlayer) CreateTime() time.Time {
	return p.createTime
}

func (p *defaultPlayer) SetCreateTime(createTime time.Time) {
	if p.createTime != createTime {
		p.createTime = createTime
		p.updatedFields.Set(6)
	}
}

func (p *defaultPlayer) UpdateTime() time.Time {
	return p.updateTime
}

func (p *defaultPlayer) SetUpdateTime(updateTime time.Time) {
	if p.updateTime != updateTime {
		p.updateTime = updateTime
		p.updatedFields.Set(7)
	}
}

func NewPlayer() Player {
	player := &defaultPlayer{updatedFields: &bitset.BitSet{}}
	player.wallet = NewWallet(player)
	player.equipments = bsonmodel.NewStringObjectMapModel(player, "eqm", EquipmentFactory())
	player.items = bsonmodel.NewIntSimpleMapModel(player, "itm", bsonmodel.IntValueType())
	return player
}

func LoadPlayerFromDocument(m bson.M) (player Player, err error) {
	player = NewPlayer()
	err = player.LoadDocument(m)
	return
}

func LoadPlayerFromJsoniter(any jsoniter.Any) (player Player, err error) {
	player = NewPlayer()
	err = player.LoadJsoniter(any)
	return
}

type playerEncoder struct{}

func (codec *playerEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

func (codec *playerEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	p := ((*defaultPlayer)(ptr))
	stream.WriteObjectStart()
	stream.WriteObjectField("uid")
	stream.WriteInt(p.uid)
	stream.WriteMore()
	stream.WriteObjectField("wallet")
	stream.WriteVal(p.wallet)
	stream.WriteMore()
	stream.WriteObjectField("equipments")
	stream.WriteVal(p.equipments)
	stream.WriteMore()
	stream.WriteObjectField("items")
	stream.WriteVal(p.items)
	stream.WriteObjectEnd()
}

func init() {
	jsoniter.RegisterTypeEncoder("example.defaultPlayer", &playerEncoder{})
}
