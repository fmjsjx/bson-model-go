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
	SetCoinUsed(coinTotal int)
	Coin() int
	Diamond() int
	SetDiamond(diamond int)
}

const (
	BnameWalletCoinTotal = "ct"
	BnameWalletCoinUsed  = "cu"
	BnameWalletDiamond   = "d"
)

var xpathWallet bsonmodel.DotNotation = bsonmodel.PathOfNames("wt")

type defaultWallet struct {
	updatedFields *bitset.BitSet
	parent        Player
	coinTotal     int
	coinUsed      int
	diamond       int
}

func (w *defaultWallet) ToBson() interface{} {
	return w.ToDocument()
}

func (w *defaultWallet) ToData() interface{} {
	data := make(map[string]interface{})
	data["ct"] = w.coinTotal
	data["cu"] = w.coinUsed
	data["d"] = w.diamond
	return data
}

func (w *defaultWallet) LoadJsoniter(any jsoniter.Any) error {
	if any.ValueType() != jsoniter.ObjectValue {
		return nil
	}
	ct, err := bsonmodel.AnyIntValue(any.Get("ct"), 0)
	if err != nil {
		return err
	}
	w.coinTotal = ct
	cu, err := bsonmodel.AnyIntValue(any.Get("cu"), 0)
	if err != nil {
		return err
	}
	w.coinUsed = cu
	d, err := bsonmodel.AnyIntValue(any.Get("d"), 0)
	if err != nil {
		return err
	}
	w.diamond = d
	return nil
}

func (w *defaultWallet) Reset() {
	w.updatedFields.ClearAll()
}

func (w *defaultWallet) AnyUpdated() bool {
	return w.updatedFields.Any()
}

func (w *defaultWallet) AnyDeleted() bool {
	return w.DeletedSize() > 0
}

func (w *defaultWallet) Parent() bsonmodel.BsonModel {
	return w.parent
}

func (w *defaultWallet) XPath() bsonmodel.DotNotation {
	return xpathWallet
}

func (w *defaultWallet) AppendUpdates(updates bson.M) bson.M {
	dset := bsonmodel.FixedEmbedded(updates, "$set")
	xpath := w.XPath()
	if w.FullyUpdate() {
		dset[xpath.Value()] = w.ToDocument()
	} else {
		updatedFields := w.updatedFields
		if updatedFields.Test(1) {
			dset[xpath.Resolve("ct").Value()] = w.coinTotal
		}
		if updatedFields.Test(2) {
			dset[xpath.Resolve("cu").Value()] = w.coinUsed
		}
		if updatedFields.Test(4) {
			dset[xpath.Resolve("d").Value()] = w.diamond
		}
	}
	return updates
}

func (w *defaultWallet) ToDocument() bson.M {
	doc := bson.M{}
	doc["ct"] = w.coinTotal
	doc["cu"] = w.coinUsed
	doc["d"] = w.diamond
	return doc
}

func (w *defaultWallet) LoadDocument(document bson.M) error {
	coinTotal, err := bsonmodel.IntValue(document, "ct", 0)
	if err != nil {
		return err
	}
	w.coinTotal = coinTotal
	coinUsed, err := bsonmodel.IntValue(document, "cu", 0)
	if err != nil {
		return err
	}
	w.coinUsed = coinUsed
	diamond, err := bsonmodel.IntValue(document, "d", 0)
	if err != nil {
		return err
	}
	w.diamond = diamond
	return nil
}

func (w *defaultWallet) DeletedSize() int {
	return 0
}

func (w *defaultWallet) FullyUpdate() bool {
	return w.updatedFields.Test(0)
}

func (w *defaultWallet) SetFullyUpdate(fullyUpdate bool) {
	if fullyUpdate {
		w.updatedFields.Set(0)
	} else {
		w.updatedFields.DeleteAt(0)
	}
}

func (w *defaultWallet) CoinTotal() int {
	return w.coinTotal
}

func (w *defaultWallet) SetCoinTotal(coinTotal int) {
	if w.coinTotal != coinTotal {
		w.coinTotal = coinTotal
		w.updatedFields.Set(1)
	}
}

func (w *defaultWallet) CoinUsed() int {
	return w.coinUsed
}

func (w *defaultWallet) SetCoinUsed(coinUsed int) {
	if w.coinTotal != coinUsed {
		w.coinUsed = coinUsed
		w.updatedFields.Set(2)
	}
}

func (w *defaultWallet) Coin() int {
	return w.coinTotal - w.coinUsed
}

func (w *defaultWallet) Diamond() int {
	return w.diamond
}

func (w *defaultWallet) SetDiamond(diamond int) {
	if w.diamond != diamond {
		w.diamond = diamond
		w.updatedFields.Set(4)
	}
}

func NewWallet(parent Player) Wallet {
	wallet := &defaultWallet{updatedFields: &bitset.BitSet{}, parent: parent}
	return wallet
}

type walletEncoder struct{}

func (codec *walletEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

func (codec *walletEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	w := ((*defaultWallet)(ptr))
	stream.WriteObjectStart()
	stream.WriteObjectField("coinTotal")
	stream.WriteInt(w.coinTotal)
	stream.WriteMore()
	stream.WriteObjectField("coin")
	stream.WriteInt(w.Coin())
	stream.WriteMore()
	stream.WriteObjectField("diamond")
	stream.WriteInt(w.diamond)
	stream.WriteObjectEnd()
}

func init() {
	jsoniter.RegisterTypeEncoder("example.defaultWallet", &walletEncoder{})
}
