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

type Player interface {
	bsonmodel.RootModel
	Uid() int
	SetUid(uid int)
	Wallet() Wallet
	Equipments() bsonmodel.StringObjectMapModel
	Equipment(id string) Equipment
	Items() bsonmodel.IntSimpleMapModel
	Cash() CashInfo
	UpdateVersion() int
	SetUpdateVersion(updateVersion int)
	IncreaseUpdateVersion() int
	CreateTime() time.Time
	SetCreateTime(createTime time.Time)
	UpdateTime() time.Time
	SetUpdateTime(updateTime time.Time)
}

const (
	BnamePlayerUid           = "_id"
	BnamePlayerWallet        = "wlt"
	BnamePlayerEquipments    = "eqm"
	BnamePlayerItems         = "itm"
	BnamePlayerCash          = "cs"
	BnamePlayerUpdateVersion = "_uv"
	BnamePlayerCreateTime    = "_ct"
	BnamePlayerUpdateTime    = "_ut"
)

type defaultPlayer struct {
	updatedFields *bitset.BitSet
	uid           int
	wallet        Wallet
	equipments    bsonmodel.StringObjectMapModel
	items         bsonmodel.IntSimpleMapModel
	cash          CashInfo
	updateVersion int
	createTime    time.Time
	updateTime    time.Time
}

func (self *defaultPlayer) ToBson() interface{} {
	return self.ToDocument()
}

func (self *defaultPlayer) ToData() interface{} {
	data := make(map[string]interface{})
	data["_id"] = self.uid
	data["wlt"] = self.wallet.ToData()
	data["eqm"] = self.equipments.ToData()
	data["itm"] = self.items.ToData()
	data["cs"] = self.cash.ToData()
	data["_uv"] = self.updateVersion
	data["_ct"] = self.createTime.UnixMilli()
	data["_ut"] = self.updateTime.UnixMilli()
	return data
}

func (self *defaultPlayer) LoadJsoniter(any jsoniter.Any) error {
	if any.ValueType() != jsoniter.ObjectValue {
		self.Reset()
		return nil
	}
	uid, err := bsonmodel.AnyIntValue(any.Get("_id"), 0)
	if err != nil {
		return err
	}
	self.uid = uid
	wallet := any.Get("wlt")
	if wallet.ValueType() == jsoniter.ObjectValue {
		err = self.wallet.LoadJsoniter(wallet)
		if err != nil {
			return err
		}
	}
	equipments := any.Get("eqm")
	if equipments.ValueType() == jsoniter.ObjectValue {
		err = self.equipments.LoadJsoniter(equipments)
		if err != nil {
			return err
		}
	} else {
		self.equipments.Clear()
	}
	items := any.Get("itm")
	if items.ValueType() == jsoniter.ObjectValue {
		err = self.items.LoadJsoniter(items)
		if err != nil {
			return err
		}
	} else {
		self.items.Reset()
	}
	cash := any.Get("cs")
	if cash.ValueType() == jsoniter.ObjectValue {
		err = self.cash.LoadJsoniter(cash)
		if err != nil {
			return err
		}
	}
	updateVersion, err := bsonmodel.AnyIntValue(any.Get("_uv"), 0)
	if err != nil {
		return err
	}
	self.updateVersion = updateVersion
	createTime, err := bsonmodel.AnyDateTimeValue(any.Get("_ct"))
	if err != nil {
		return err
	}
	self.createTime = createTime
	updateTime, err := bsonmodel.AnyDateTimeValue(any.Get("_ut"))
	if err != nil {
		return err
	}
	self.updateTime = updateTime
	self.Reset()
	return nil
}

func (self *defaultPlayer) Reset() {
	self.wallet.Reset()
	self.equipments.Reset()
	self.items.Reset()
	self.cash.Reset()
	self.updatedFields.ClearAll()
}

func (self *defaultPlayer) AnyUpdated() bool {
	return self.updatedFields.Any() || self.wallet.AnyUpdated() || self.equipments.AnyUpdated() || self.items.AnyUpdated() || self.cash.AnyUpdated()
}

func (self *defaultPlayer) AnyDeleted() bool {
	return self.DeletedSize() > 0
}

func (self *defaultPlayer) Parent() bsonmodel.BsonModel {
	return nil
}

func (self *defaultPlayer) XPath() bsonmodel.DotNotation {
	return bsonmodel.RootPath()
}

func (self *defaultPlayer) AppendUpdates(updates bson.M) bson.M {
	dset := bsonmodel.FixedEmbedded(updates, "$set")
	updatedFields := self.updatedFields
	if updatedFields.Test(1) {
		dset["_id"] = self.uid
	}
	if self.wallet.AnyUpdated() {
		self.wallet.AppendUpdates(updates)
	}
	if self.equipments.AnyUpdated() {
		self.equipments.AppendUpdates(updates)
	}
	if self.items.AnyUpdated() {
		self.items.AppendUpdates(updates)
	}
	if self.cash.AnyUpdated() {
		self.cash.AppendUpdates(updates)
	}
	if updatedFields.Test(6) {
		dset["_uv"] = self.updateVersion
	}
	if updatedFields.Test(7) {
		dset["_ct"] = primitive.NewDateTimeFromTime(self.createTime)
	}
	if updatedFields.Test(8) {
		dset["_ut"] = primitive.NewDateTimeFromTime(self.updateTime)
	}
	return updates
}

func (self *defaultPlayer) ToDocument() bson.M {
	doc := bson.M{}
	doc["_id"] = self.uid
	doc["wlt"] = self.wallet.ToBson()
	doc["eqm"] = self.equipments.ToBson()
	doc["itm"] = self.items.ToBson()
	doc["cs"] = self.cash.ToBson()
	doc["_uv"] = self.updateVersion
	doc["_ct"] = primitive.NewDateTimeFromTime(self.createTime)
	doc["_ut"] = primitive.NewDateTimeFromTime(self.updateTime)
	return doc
}

func (self *defaultPlayer) LoadDocument(document bson.M) error {
	uid, err := bsonmodel.IntValue(document, "_id", 0)
	if err != nil {
		return err
	}
	self.uid = uid
	wallet, err := bsonmodel.EmbeddedValue(document, "wlt")
	if err != nil {
		return err
	}
	if wallet != nil {
		err = self.wallet.LoadDocument(wallet)
		if err != nil {
			return err
		}
	}
	equipments, err := bsonmodel.EmbeddedValue(document, "eqm")
	if err != nil {
		return err
	}
	if equipments != nil {
		err = self.equipments.LoadDocument(equipments)
		if err != nil {
			return err
		}
	} else {
		self.equipments.Clear()
	}
	items, err := bsonmodel.EmbeddedValue(document, "itm")
	if err != nil {
		return err
	}
	if items != nil {
		err = self.items.LoadDocument(items)
		if err != nil {
			return err
		}
	} else {
		self.items.Clear()
	}
	cash, err := bsonmodel.EmbeddedValue(document, "cs")
	if err != nil {
		return err
	}
	if cash != nil {
		err = self.cash.LoadDocument(cash)
		if err != nil {
			return err
		}
	}
	updateVersion, err := bsonmodel.IntValue(document, "_uv", 0)
	if err != nil {
		return err
	}
	self.updateVersion = updateVersion
	createTime, err := bsonmodel.DateTimeValue(document, "_ct")
	if err != nil {
		return err
	}
	self.createTime = createTime
	updateTime, err := bsonmodel.DateTimeValue(document, "_ut")
	if err != nil {
		return err
	}
	self.updateTime = updateTime
	self.Reset()
	return nil
}

func (self *defaultPlayer) DeletedSize() int {
	n := 0
	if self.wallet.AnyDeleted() {
		n += 1
	}
	if self.equipments.AnyDeleted() {
		n += 1
	}
	if self.items.AnyDeleted() {
		n += 1
	}
	if self.cash.AnyDeleted() {
		n += 1
	}
	return n
}

func (self *defaultPlayer) FullyUpdate() bool {
	return false
}

func (self *defaultPlayer) SetFullyUpdate(fullyUpdate bool) {
	// no effect
}

func (self *defaultPlayer) ToSync() interface{} {
	sync := make(map[string]interface{})
	updatedFields := self.updatedFields
	if updatedFields.Test(1) {
		sync["uid"] = self.uid
	}
	if self.wallet.AnyUpdated() {
		sync["wallet"] = self.wallet.ToSync()
	}
	if self.equipments.AnyUpdated() {
		sync["equipments"] = self.equipments.ToSync()
	}
	if self.items.AnyUpdated() {
		sync["items"] = self.items.ToSync()
	}
	if self.cash.AnyUpdated() {
		sync["cash"] = self.cash.ToSync()
	}
	return sync
}

func (self *defaultPlayer) ToDelete() interface{} {
	delete := make(map[string]interface{})
	if self.wallet.AnyDeleted() {
		delete["wallet"] = self.wallet.ToDelete()
	}
	if self.equipments.AnyDeleted() {
		delete["equipments"] = self.equipments.ToDelete()
	}
	if self.items.AnyDeleted() {
		delete["items"] = self.items.ToDelete()
	}
	if self.cash.AnyDeleted() {
		delete["cash"] = self.cash.ToDelete()
	}
	return delete
}

func (self *defaultPlayer) ToDataJson() (string, error) {
	return jsoniter.MarshalToString(self.ToData())
}

func (self *defaultPlayer) ToSyncJson() (string, error) {
	return jsoniter.MarshalToString(self.ToSync())
}

func (self *defaultPlayer) ToDeleteJson() (string, error) {
	return jsoniter.MarshalToString(self.ToDelete())
}

func (self *defaultPlayer) ToUpdate() bson.M {
	if self.AnyUpdated() {
		return self.AppendUpdates(bson.M{})
	}
	return bson.M{}
}

func (self *defaultPlayer) MarshalToJsonString() (string, error) {
	return jsoniter.MarshalToString(self)
}

func (self *defaultPlayer) Uid() int {
	return self.uid
}

func (self *defaultPlayer) SetUid(uid int) {
	if self.uid != uid {
		self.uid = uid
		self.updatedFields.Set(1)
	}
}

func (self *defaultPlayer) Wallet() Wallet {
	return self.wallet
}

func (self *defaultPlayer) Equipments() bsonmodel.StringObjectMapModel {
	return self.equipments
}

func (self *defaultPlayer) Equipment(id string) Equipment {
	value := self.equipments.Get(id)
	if value == nil {
		return nil
	}
	return value.(Equipment)
}

func (self *defaultPlayer) Items() bsonmodel.IntSimpleMapModel {
	return self.items
}

func (self *defaultPlayer) Cash() CashInfo {
	return self.cash
}

func (self *defaultPlayer) UpdateVersion() int {
	return self.updateVersion
}

func (self *defaultPlayer) SetUpdateVersion(updateVersion int) {
	if self.updateVersion != updateVersion {
		self.updateVersion = updateVersion
		self.updatedFields.Set(6)
	}
}

func (self *defaultPlayer) IncreaseUpdateVersion() int {
	updateVersion := self.updateVersion + 1
	self.updateVersion = updateVersion
	self.updatedFields.Set(6)
	return updateVersion
}

func (self *defaultPlayer) CreateTime() time.Time {
	return self.createTime
}

func (self *defaultPlayer) SetCreateTime(createTime time.Time) {
	if self.createTime != createTime {
		self.createTime = createTime
		self.updatedFields.Set(7)
	}
}

func (self *defaultPlayer) UpdateTime() time.Time {
	return self.updateTime
}

func (self *defaultPlayer) SetUpdateTime(updateTime time.Time) {
	if self.updateTime != updateTime {
		self.updateTime = updateTime
		self.updatedFields.Set(8)
	}
}

func NewPlayer() Player {
	self := &defaultPlayer{updatedFields: &bitset.BitSet{}}
	self.wallet = NewWallet(self)
	self.equipments = bsonmodel.NewStringObjectMapModel(self, "eqm", EquipmentFactory())
	self.items = bsonmodel.NewIntSimpleMapModel(self, "itm", bsonmodel.IntValueType())
	self.cash = NewCashInfo(self)
	return self
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
	stream.WriteMore()
	stream.WriteObjectField("cash")
	stream.WriteVal(p.cash)
	stream.WriteObjectEnd()
}

func init() {
	jsoniter.RegisterTypeEncoder("example.defaultPlayer", &playerEncoder{})
}

