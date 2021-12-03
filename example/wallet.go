package example

import (
	"unsafe"

	"github.com/bits-and-blooms/bitset"
	"github.com/fmjsjx/bson-model-go/bsonmodel"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
)

type Wallet interface {
	bsonmodel.ObjectModel
	CoinTotal() int
	SetCoinTotal(coinTotal int)
	CoinUsed() int
	SetCoinUsed(coinUsed int)
	Coin() int
	Diamond() int
	SetDiamond(diamond int)
}

const (
	BnameWalletCoinTotal = "ct"
	BnameWalletCoinUsed  = "cu"
	BnameWalletDiamond   = "d"
)

var xpathWallet bsonmodel.DotNotation = bsonmodel.PathOfNames("wlt")

type defaultWallet struct {
	updatedFields *bitset.BitSet
	parent        Player
	coinTotal     int
	coinUsed      int
	diamond       int
}

func (self *defaultWallet) ToBson() interface{} {
	return self.ToDocument()
}

func (self *defaultWallet) ToData() interface{} {
	data := make(map[string]interface{})
	data["ct"] = self.coinTotal
	data["cu"] = self.coinUsed
	data["d"] = self.diamond
	return data
}

func (self *defaultWallet) LoadJsoniter(any jsoniter.Any) error {
	if any.ValueType() != jsoniter.ObjectValue {
		return nil
	}
	coinTotal, err := bsonmodel.AnyIntValue(any.Get("ct"), 0)
	if err != nil {
		return err
	}
	self.coinTotal = coinTotal
	coinUsed, err := bsonmodel.AnyIntValue(any.Get("cu"), 0)
	if err != nil {
		return err
	}
	self.coinUsed = coinUsed
	diamond, err := bsonmodel.AnyIntValue(any.Get("d"), 0)
	if err != nil {
		return err
	}
	self.diamond = diamond
	return nil
}

func (self *defaultWallet) Reset() {
	self.updatedFields.ClearAll()
}

func (self *defaultWallet) AnyUpdated() bool {
	return self.updatedFields.Any()
}

func (self *defaultWallet) AnyDeleted() bool {
	return self.DeletedSize() > 0
}

func (self *defaultWallet) Parent() bsonmodel.BsonModel {
	return self.parent
}

func (self *defaultWallet) XPath() bsonmodel.DotNotation {
	return xpathWallet
}

func (self *defaultWallet) AppendUpdates(updates bson.M) bson.M {
	dset := bsonmodel.FixedEmbedded(updates, "$set")
	xpath := self.XPath()
	if self.FullyUpdate() {
		dset[xpath.Value()] = self.ToDocument()
	} else {
		updatedFields := self.updatedFields
		if updatedFields.Test(1) {
			dset[xpath.Resolve("ct").Value()] = self.coinTotal
		}
		if updatedFields.Test(2) {
			dset[xpath.Resolve("cu").Value()] = self.coinUsed
		}
		if updatedFields.Test(4) {
			dset[xpath.Resolve("d").Value()] = self.diamond
		}
	}
	return updates
}

func (self *defaultWallet) ToDocument() bson.M {
	doc := bson.M{}
	doc["ct"] = self.coinTotal
	doc["cu"] = self.coinUsed
	doc["d"] = self.diamond
	return doc
}

func (self *defaultWallet) LoadDocument(document bson.M) error {
	coinTotal, err := bsonmodel.IntValue(document, "ct", 0)
	if err != nil {
		return err
	}
	self.coinTotal = coinTotal
	coinUsed, err := bsonmodel.IntValue(document, "cu", 0)
	if err != nil {
		return err
	}
	self.coinUsed = coinUsed
	diamond, err := bsonmodel.IntValue(document, "d", 0)
	if err != nil {
		return err
	}
	self.diamond = diamond
	return nil
}

func (self *defaultWallet) DeletedSize() int {
	return 0
}

func (self *defaultWallet) FullyUpdate() bool {
	return self.updatedFields.Test(0)
}

func (self *defaultWallet) SetFullyUpdate(fullyUpdate bool) {
	if fullyUpdate {
		self.updatedFields.Set(0)
	} else {
		self.updatedFields.DeleteAt(0)
	}
}

func (self *defaultWallet) ToSync() interface{} {
	if self.FullyUpdate() {
		return self
	}
	sync := make(map[string]interface{})
	updatedFields := self.updatedFields
	if updatedFields.Test(1) {
		sync["coinTotal"] = self.coinTotal
	}
	if updatedFields.Test(3) {
		sync["coin"] = self.Coin()
	}
	if updatedFields.Test(4) {
		sync["diamond"] = self.diamond
	}
	return sync
}

func (self *defaultWallet) ToDelete() interface{} {
	delete := make(map[string]interface{})
	return delete
}

func (self *defaultWallet) ToDataJson() (string, error) {
	return jsoniter.MarshalToString(self.ToData())
}

func (self *defaultWallet) ToSyncJson() (string, error) {
	return jsoniter.MarshalToString(self.ToSync())
}

func (self *defaultWallet) ToDeleteJson() (string, error) {
	return jsoniter.MarshalToString(self.ToDelete())
}

func (self *defaultWallet) MarshalJSON() ([]byte, error) {
	return jsoniter.Marshal(self)
}

func (self *defaultWallet) CoinTotal() int {
	return self.coinTotal
}

func (self *defaultWallet) SetCoinTotal(coinTotal int) {
	if self.coinTotal != coinTotal {
		self.coinTotal = coinTotal
		self.updatedFields.Set(1)
		self.updatedFields.Set(3)
	}
}

func (self *defaultWallet) CoinUsed() int {
	return self.coinUsed
}

func (self *defaultWallet) SetCoinUsed(coinUsed int) {
	if self.coinUsed != coinUsed {
		self.coinUsed = coinUsed
		self.updatedFields.Set(2)
		self.updatedFields.Set(3)
	}
}

func (self *defaultWallet) Coin() int {
	return self.coinTotal - self.coinUsed
}

func (self *defaultWallet) Diamond() int {
	return self.diamond
}

func (self *defaultWallet) SetDiamond(diamond int) {
	if self.diamond != diamond {
		self.diamond = diamond
		self.updatedFields.Set(4)
	}
}

func NewWallet(parent Player) Wallet {
	self := &defaultWallet{updatedFields: &bitset.BitSet{}, parent: parent}
	return self
}

type walletEncoder struct{}

func (codec *walletEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

func (codec *walletEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	p := ((*defaultWallet)(ptr))
	stream.WriteObjectStart()
	stream.WriteObjectField("coinTotal")
	stream.WriteInt(p.coinTotal)
	stream.WriteMore()
	stream.WriteObjectField("coin")
	stream.WriteInt(p.Coin())
	stream.WriteMore()
	stream.WriteObjectField("diamond")
	stream.WriteInt(p.diamond)
	stream.WriteObjectEnd()
}

func init() {
	jsoniter.RegisterTypeEncoder("example.defaultWallet", &walletEncoder{})
}

