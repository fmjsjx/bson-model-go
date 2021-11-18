package example

import (
	"unsafe"

	"github.com/bits-and-blooms/bitset"
	"github.com/fmjsjx/bson-model-go/bsonmodel"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
)

type CashInfo interface {
	bsonmodel.ObjectModel
	Stages() bsonmodel.IntSimpleMapModel
	Cards() []int
	SetCards(cards []int)
	OrderIds() []string
	SetOrderIds(orderIds []string)
}

const (
	BnameCashInfoStages   = "stg"
	BnameCashInfoCards    = "cs"
	BnameCashInfoOrderIds = "ois"
)

var xpathCashInfo bsonmodel.DotNotation = bsonmodel.PathOfNames("cs")

type defaultCashInfo struct {
	updatedFields *bitset.BitSet
	parent        Player
	stages        bsonmodel.IntSimpleMapModel
	cards         []int
	orderIds      []string
}

func (self *defaultCashInfo) ToBson() interface{} {
	return self.ToDocument()
}

func (self *defaultCashInfo) ToData() interface{} {
	data := make(map[string]interface{})
	data["stg"] = self.stages.ToData()
	if self.cards != nil {
		data["cs"] = self.cards
	}
	if self.orderIds != nil {
		data["ois"] = self.orderIds
	}
	return data
}

func (self *defaultCashInfo) LoadJsoniter(any jsoniter.Any) error {
	if any.ValueType() != jsoniter.ObjectValue {
		return nil
	}
	stages := any.Get("stg")
	if stages.ValueType() == jsoniter.ObjectValue {
		err := self.stages.LoadJsoniter(stages)
		if err != nil {
			return err
		}
	} else {
		self.stages.Reset()
	}
	cards, err := bsonmodel.AnyIntArrayValue(any.Get("cs"))
	if err != nil {
		return err
	}
	self.cards = cards
	orderIds, err := bsonmodel.AnyStringArrayValue(any.Get("ois"))
	if err != nil {
		return err
	}
	self.orderIds = orderIds
	return nil
}

func (self *defaultCashInfo) Reset() {
	self.stages.Reset()
	self.updatedFields.ClearAll()
}

func (self *defaultCashInfo) AnyUpdated() bool {
	return self.updatedFields.Any() || self.stages.AnyUpdated()
}

func (self *defaultCashInfo) AnyDeleted() bool {
	return self.DeletedSize() > 0
}

func (self *defaultCashInfo) Parent() bsonmodel.BsonModel {
	return self.parent
}

func (self *defaultCashInfo) XPath() bsonmodel.DotNotation {
	return xpathCashInfo
}

func (self *defaultCashInfo) AppendUpdates(updates bson.M) bson.M {
	dset := bsonmodel.FixedEmbedded(updates, "$set")
	xpath := self.XPath()
	if self.FullyUpdate() {
		dset[xpath.Value()] = self.ToDocument()
	} else {
		updatedFields := self.updatedFields
		if self.stages.AnyUpdated() {
			self.stages.AppendUpdates(updates)
		}
		if updatedFields.Test(2) {
			if self.cards == nil {
				bsonmodel.FixedEmbedded(updates, "$unset")[xpath.Resolve("cs").Value()] = ""
			} else {
				cardsArray := bson.A{}
				for _, v := range self.cards {
					cardsArray = append(cardsArray, v)
				}
				dset[xpath.Resolve("cs").Value()] = cardsArray
			}
		}
		if updatedFields.Test(3) {
			if self.orderIds == nil {
				bsonmodel.FixedEmbedded(updates, "$unset")[xpath.Resolve("ois").Value()] = ""
			} else {
				orderIdsArray := bson.A{}
				for _, v := range self.orderIds {
					orderIdsArray = append(orderIdsArray, v)
				}
				dset[xpath.Resolve("ois").Value()] = orderIdsArray
			}
		}
	}
	return updates
}

func (self *defaultCashInfo) ToDocument() bson.M {
	doc := bson.M{}
	doc["stg"] = self.stages.ToBson()
	if self.cards != nil {
		cardsArray := bson.A{}
		for _, v := range self.cards {
			cardsArray = append(cardsArray, v)
		}
		doc["cs"] = cardsArray
	}
	if self.orderIds != nil {
		orderIdsArray := bson.A{}
		for _, v := range self.orderIds {
			orderIdsArray = append(orderIdsArray, v)
		}
		doc["ois"] = orderIdsArray
	}
	return doc
}

func (self *defaultCashInfo) LoadDocument(document bson.M) error {
	stages, err := bsonmodel.EmbeddedValue(document, "stg")
	if err != nil {
		return err
	}
	if stages != nil {
		err = self.stages.LoadDocument(stages)
		if err != nil {
			return err
		}
	} else {
		self.stages.Clear()
	}
	cards, err := bsonmodel.IntArrayValue(document, "cs")
	if err != nil {
		return err
	}
	self.cards = cards
	orderIds, err := bsonmodel.StringArrayValue(document, "ois")
	if err != nil {
		return err
	}
	self.orderIds = orderIds
	return nil
}

func (self *defaultCashInfo) DeletedSize() int {
	n := 0
	if self.stages.AnyDeleted() {
		n += 1
	}
	if self.updatedFields.Test(2) && self.cards == nil {
		n += 1
	}
	if self.updatedFields.Test(3) && self.orderIds == nil {
		n += 1
	}
	return n
}

func (self *defaultCashInfo) FullyUpdate() bool {
	return self.updatedFields.Test(0)
}

func (self *defaultCashInfo) SetFullyUpdate(fullyUpdate bool) {
	if fullyUpdate {
		self.updatedFields.Set(0)
	} else {
		self.updatedFields.DeleteAt(0)
	}
}

func (self *defaultCashInfo) ToSync() interface{} {
	if self.FullyUpdate() {
		return self
	}
	sync := make(map[string]interface{})
	updatedFields := self.updatedFields
	if self.stages.AnyUpdated() {
		sync["stages"] = self.stages.ToSync()
	}
	if updatedFields.Test(2) {
		if self.cards != nil {
			sync["cards"] = self.cards
		}
	}
	if updatedFields.Test(3) {
		if self.orderIds != nil {
			sync["orderIds"] = self.orderIds
		}
	}
	return sync
}

func (self *defaultCashInfo) ToDelete() interface{} {
	delete := make(map[string]interface{})
	if self.stages.AnyDeleted() {
		delete["stages"] = self.stages.ToDelete()
	}
	if self.updatedFields.Test(2) {
		delete["cards"] = 1
	}
	if self.updatedFields.Test(3) {
		delete["orderIds"] = 1
	}
	return delete
}

func (self *defaultCashInfo) ToDataJson() (string, error) {
	return jsoniter.MarshalToString(self.ToData())
}

func (self *defaultCashInfo) ToSyncJson() (string, error) {
	return jsoniter.MarshalToString(self.ToSync())
}

func (self *defaultCashInfo) ToDeleteJson() (string, error) {
	return jsoniter.MarshalToString(self.ToDelete())
}

func (self *defaultCashInfo) Stages() bsonmodel.IntSimpleMapModel {
	return self.stages
}

func (self *defaultCashInfo) Cards() []int {
	return self.cards
}

func (self *defaultCashInfo) SetCards(cards []int) {
	self.cards = cards
	self.updatedFields.Set(2)
}

func (self *defaultCashInfo) OrderIds() []string {
	return self.orderIds
}

func (self *defaultCashInfo) SetOrderIds(orderIds []string) {
	self.orderIds = orderIds
	self.updatedFields.Set(3)
}

func NewCashInfo(parent Player) CashInfo {
	self := &defaultCashInfo{updatedFields: &bitset.BitSet{}, parent: parent}
	self.stages = bsonmodel.NewIntSimpleMapModel(self, "stg", bsonmodel.IntValueType())
	return self
}

type cashInfoEncoder struct{}

func (codec *cashInfoEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

func (codec *cashInfoEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	p := ((*defaultCashInfo)(ptr))
	stream.WriteObjectStart()
	stream.WriteObjectField("stages")
	stream.WriteVal(p.stages)
	stream.WriteMore()
	stream.WriteObjectField("cards")
	stream.WriteVal(p.cards)
	stream.WriteMore()
	stream.WriteObjectField("orderIds")
	stream.WriteVal(p.orderIds)
	stream.WriteObjectEnd()
}

func init() {
	jsoniter.RegisterTypeEncoder("example.defaultCashInfo", &cashInfoEncoder{})
}
