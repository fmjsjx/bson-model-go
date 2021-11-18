package example

import (
	"fmt"
	"testing"
	"time"

	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestLoadPlayerFromDocument(t *testing.T) {
	now := time.Now().Truncate(time.Millisecond)
	doc := bson.M{}
	doc["_id"] = int32(123)
	doc["_uv"] = int32(1)
	doc["wlt"] = bson.M{"ct": int32(5000), "cu": int32(2000), "d": int32(10)}
	eqm := bson.M{}
	eq0 := bson.M{"id": "12345678-1234-5678-9abc-123456789abc"}
	eq0["rid"] = int32(1001)
	eq0["atk"] = int32(12)
	eq0["def"] = int32(0)
	eq0["hp"] = int32(0)
	eqm[eq0["id"].(string)] = eq0
	eq1 := bson.M{"id": "11111111-1111-1111-1111-111111111111"}
	eq1["rid"] = int32(1101)
	eq1["atk"] = int32(0)
	eq1["def"] = int32(6)
	eq1["hp"] = int32(12)
	eqm[eq1["id"].(string)] = eq1
	doc["eqm"] = eqm
	doc["itm"] = bson.M{"2001": int32(10), "2002": int32(1)}
	cs := bson.M{}
	cs["stg"] = bson.M{"1": int32(2), "2": int32(1)}
	cs["cs"] = bson.A{int32(1), int32(2)}
	cs["ois"] = bson.A{"order-0", "order-1"}
	doc["cs"] = cs
	doc["_ct"] = primitive.NewDateTimeFromTime(now)
	doc["_ut"] = primitive.NewDateTimeFromTime(now)
	player, err := LoadPlayerFromDocument(doc)
	if err != nil {
		t.Errorf("Unexpected error occurs: %e", err)
	}
	if 123 != player.Uid() {
		t.Errorf("The value expected <%v> but was <%v>", 123, player.Uid())
	}
	if 1 != player.UpdateVersion() {
		t.Errorf("The value expected <%v> but was <%v>", 1, player.UpdateVersion())
	}
	if now != player.CreateTime() {
		t.Errorf("The value expected <%v> but was <%v>", now, player.CreateTime())
	}
	if now != player.UpdateTime() {
		t.Errorf("The value expected <%v> but was <%v>", now, player.UpdateTime())
	}
	// wallet
	wallet := player.Wallet()
	if 5000 != wallet.CoinTotal() {
		t.Errorf("The value expected <%v> but was <%v>", 5000, wallet.CoinTotal())
	}
	if 2000 != wallet.CoinUsed() {
		t.Errorf("The value expected <%v> but was <%v>", 2000, wallet.CoinUsed())
	}
	if 3000 != wallet.Coin() {
		t.Errorf("The value expected <%v> but was <%v>", 3000, wallet.Coin())
	}
	if 10 != wallet.Diamond() {
		t.Errorf("The value expected <%v> but was <%v>", 10, wallet.Diamond())
	}
	// equipments
	equipments := player.Equipments()
	if 2 != equipments.Size() {
		t.Errorf("The value expected <%v> but was <%v>", 2, equipments.Size())
	}
	equipment0 := player.Equipment("12345678-1234-5678-9abc-123456789abc")
	if equipment0 == nil {
		t.Error("The value expected not be nil")
	} else {
		e := equipments.Get("12345678-1234-5678-9abc-123456789abc")
		if equipment0 != e {
			t.Errorf("The value expected <%v> but was <%v>", equipment0, e)
		}
		if 1001 != equipment0.RefId() {
			t.Errorf("The value expected <%v> but was <%v>", 1001, equipment0.RefId())
		}
		if 12 != equipment0.Atk() {
			t.Errorf("The value expected <%v> but was <%v>", 12, equipment0.Atk())
		}
		if 0 != equipment0.Def() {
			t.Errorf("The value expected <%v> but was <%v>", 0, equipment0.Def())
		}
		if 0 != equipment0.Hp() {
			t.Errorf("The value expected <%v> but was <%v>", 0, equipment0.Hp())
		}
	}
	equipment1 := player.Equipment("11111111-1111-1111-1111-111111111111")
	if equipment1 == nil {
		t.Error("The value expected not be nil")
	} else {
		e := equipments.Get("11111111-1111-1111-1111-111111111111")
		if equipment1 != e {
			t.Errorf("The value expected <%v> but was <%v>", equipment1, e)
		}
		if 1101 != equipment1.RefId() {
			t.Errorf("The value expected <%v> but was <%v>", 1101, equipment1.RefId())
		}
		if 0 != equipment1.Atk() {
			t.Errorf("The value expected <%v> but was <%v>", 0, equipment1.Atk())
		}
		if 6 != equipment1.Def() {
			t.Errorf("The value expected <%v> but was <%v>", 6, equipment1.Def())
		}
		if 12 != equipment1.Hp() {
			t.Errorf("The value expected <%v> but was <%v>", 12, equipment1.Hp())
		}
	}
	if equipments.Get("no such equipment") != nil {
		t.Errorf("The value expected nil but was <%v>", equipments.Get("no such equipment"))
	}
	if player.Equipment("no such equipment") != nil {
		t.Errorf("The value expected nil but was <%v>", player.Equipment("no such equipment"))
	}
	// items
	items := player.Items()
	if 2 != items.Size() {
		t.Errorf("The value expected <%v> but was <%v>", 2, items.Size())
	}
	i2001 := items.Get(2001)
	if 10 != i2001 {
		t.Errorf("The value expected <%v> but was <%v>", 10, i2001)
	}
	i2002 := items.Get(2002)
	if 1 != i2002 {
		t.Errorf("The value expected <%v> but was <%v>", 1, i2002)
	}
	if items.Get(12345) != nil {
		t.Errorf("The value expected nil but was <%v>", items.Get(12345))
	}
	// cash
	cash := player.Cash()
	stages := cash.Stages()
	if 2 != stages.Size() {
		t.Errorf("The value expected <%v> but was <%v>", 2, stages.Size())
	}
	if 2 != stages.Get(1) {
		t.Errorf("The value expected <%v> but was <%v>", 2, stages.Get(1))
	}
	if 1 != stages.Get(2) {
		t.Errorf("The value expected <%v> but was <%v>", 2, stages.Get(2))
	}
	if stages.Get(999) != nil {
		t.Errorf("The value expected nil but was <%v>", stages.Get(999))
	}
	cards := cash.Cards()
	if cards == nil {
		t.Error("The value expected not be nil")
	} else {
		if 2 != len(cards) {
			t.Errorf("The value expected <%v> but was <%v>", 2, len(cards))
		}
		if 1 != cards[0] {
			t.Errorf("The value expected <%v> but was <%v>", 1, cards[0])
		}
		if 2 != cards[1] {
			t.Errorf("The value expected <%v> but was <%v>", 2, cards[1])
		}
	}
	orderIds := cash.OrderIds()
	if orderIds == nil {
		t.Error("The value expected not be nil")
	} else {
		if 2 != len(orderIds) {
			t.Errorf("The value expected <%v> but was <%v>", 2, len(orderIds))
		}
		if "order-0" != orderIds[0] {
			t.Errorf("The value expected <%v> but was <%v>", "order-0", orderIds[0])
		}
		if "order-1" != orderIds[1] {
			t.Errorf("The value expected <%v> but was <%v>", "order-1", orderIds[1])
		}
	}
}

func TestLoadPlayerFromJsoniter(t *testing.T) {
	now := time.Now().Truncate(time.Millisecond)
	jsonfmt := `{"wlt":{"ct":5000,"cu":2000,"d":10},"eqm":{"12345678-1234-5678-9abc-123456789abc":{"id":"12345678-1234-5678-9abc-123456789abc","rid":1001,"atk":12,"def":0,"hp":0},"11111111-1111-1111-1111-111111111111":{"id":"11111111-1111-1111-1111-111111111111","rid":1101,"atk":0,"def":6,"hp":12}},"itm":{"2001":10,"2002":1},"cs":{"cs":[1,2],"ois":["order-0","order-1"],"stg":{"1":2,"2":1}},"_uv":1,"_ct":%d,"_ut":%d,"_id":123}`
	json := fmt.Sprintf(jsonfmt, now.UnixMilli(), now.UnixMilli())
	any := jsoniter.Get([]byte(json))
	player, err := LoadPlayerFromJsoniter(any)
	if err != nil {
		t.Errorf("Unexpected error occurs: %e", err)
	}
	if 123 != player.Uid() {
		t.Errorf("The value expected <%v> but was <%v>", 123, player.Uid())
	}
	if 1 != player.UpdateVersion() {
		t.Errorf("The value expected <%v> but was <%v>", 1, player.UpdateVersion())
	}
	if now != player.CreateTime() {
		t.Errorf("The value expected <%v> but was <%v>", now, player.CreateTime())
	}
	if now != player.UpdateTime() {
		t.Errorf("The value expected <%v> but was <%v>", now, player.UpdateTime())
	}
	// wallet
	wallet := player.Wallet()
	if 5000 != wallet.CoinTotal() {
		t.Errorf("The value expected <%v> but was <%v>", 5000, wallet.CoinTotal())
	}
	if 2000 != wallet.CoinUsed() {
		t.Errorf("The value expected <%v> but was <%v>", 2000, wallet.CoinUsed())
	}
	if 3000 != wallet.Coin() {
		t.Errorf("The value expected <%v> but was <%v>", 3000, wallet.Coin())
	}
	if 10 != wallet.Diamond() {
		t.Errorf("The value expected <%v> but was <%v>", 10, wallet.Diamond())
	}
	// equipments
	equipments := player.Equipments()
	if 2 != equipments.Size() {
		t.Errorf("The value expected <%v> but was <%v>", 2, equipments.Size())
	}
	equipment0 := player.Equipment("12345678-1234-5678-9abc-123456789abc")
	if equipment0 == nil {
		t.Error("The value expected not be nil")
	} else {
		e := equipments.Get("12345678-1234-5678-9abc-123456789abc")
		if equipment0 != e {
			t.Errorf("The value expected <%v> but was <%v>", equipment0, e)
		}
		if 1001 != equipment0.RefId() {
			t.Errorf("The value expected <%v> but was <%v>", 1001, equipment0.RefId())
		}
		if 12 != equipment0.Atk() {
			t.Errorf("The value expected <%v> but was <%v>", 12, equipment0.Atk())
		}
		if 0 != equipment0.Def() {
			t.Errorf("The value expected <%v> but was <%v>", 0, equipment0.Def())
		}
		if 0 != equipment0.Hp() {
			t.Errorf("The value expected <%v> but was <%v>", 0, equipment0.Hp())
		}
	}
	equipment1 := player.Equipment("11111111-1111-1111-1111-111111111111")
	if equipment1 == nil {
		t.Error("The value expected not be nil")
	} else {
		e := equipments.Get("11111111-1111-1111-1111-111111111111")
		if equipment1 != e {
			t.Errorf("The value expected <%v> but was <%v>", equipment1, e)
		}
		if 1101 != equipment1.RefId() {
			t.Errorf("The value expected <%v> but was <%v>", 1101, equipment1.RefId())
		}
		if 0 != equipment1.Atk() {
			t.Errorf("The value expected <%v> but was <%v>", 0, equipment1.Atk())
		}
		if 6 != equipment1.Def() {
			t.Errorf("The value expected <%v> but was <%v>", 6, equipment1.Def())
		}
		if 12 != equipment1.Hp() {
			t.Errorf("The value expected <%v> but was <%v>", 12, equipment1.Hp())
		}
	}
	if equipments.Get("no such equipment") != nil {
		t.Errorf("The value expected nil but was <%v>", equipments.Get("no such equipment"))
	}
	if player.Equipment("no such equipment") != nil {
		t.Errorf("The value expected nil but was <%v>", player.Equipment("no such equipment"))
	}
	// items
	items := player.Items()
	if 2 != items.Size() {
		t.Errorf("The value expected <%v> but was <%v>", 2, items.Size())
	}
	i2001 := items.Get(2001)
	if 10 != i2001 {
		t.Errorf("The value expected <%v> but was <%v>", 10, i2001)
	}
	i2002 := items.Get(2002)
	if 1 != i2002 {
		t.Errorf("The value expected <%v> but was <%v>", 1, i2002)
	}
	if items.Get(12345) != nil {
		t.Errorf("The value expected nil but was <%v>", items.Get(12345))
	}
	// cash
	cash := player.Cash()
	stages := cash.Stages()
	if 2 != stages.Size() {
		t.Errorf("The value expected <%v> but was <%v>", 2, stages.Size())
	}
	if 2 != stages.Get(1) {
		t.Errorf("The value expected <%v> but was <%v>", 2, stages.Get(1))
	}
	if 1 != stages.Get(2) {
		t.Errorf("The value expected <%v> but was <%v>", 2, stages.Get(2))
	}
	if stages.Get(999) != nil {
		t.Errorf("The value expected nil but was <%v>", stages.Get(999))
	}
	cards := cash.Cards()
	if cards == nil {
		t.Error("The value expected not be nil")
	} else {
		if 2 != len(cards) {
			t.Errorf("The value expected <%v> but was <%v>", 2, len(cards))
		}
		if 1 != cards[0] {
			t.Errorf("The value expected <%v> but was <%v>", 1, cards[0])
		}
		if 2 != cards[1] {
			t.Errorf("The value expected <%v> but was <%v>", 2, cards[1])
		}
	}
	orderIds := cash.OrderIds()
	if orderIds == nil {
		t.Error("The value expected not be nil")
	} else {
		if 2 != len(orderIds) {
			t.Errorf("The value expected <%v> but was <%v>", 2, len(orderIds))
		}
		if "order-0" != orderIds[0] {
			t.Errorf("The value expected <%v> but was <%v>", "order-0", orderIds[0])
		}
		if "order-1" != orderIds[1] {
			t.Errorf("The value expected <%v> but was <%v>", "order-1", orderIds[1])
		}
	}
}

func TestToDocument(t *testing.T) {
	now := time.Now().Truncate(time.Millisecond)
	player := NewPlayer()
	player.SetUid(123)
	player.Wallet().SetCoinTotal(5000)
	player.Wallet().SetCoinUsed(2000)
	player.Wallet().SetDiamond(10)
	equipment0 := NewEquipment()
	equipment0.SetId("12345678-1234-5678-9abc-123456789abc")
	equipment0.SetRefId(1001)
	equipment0.SetAtk(12)
	player.Equipments().Put(equipment0.Id(), equipment0)
	equipment1 := NewEquipment()
	equipment1.SetId("11111111-1111-1111-1111-111111111111")
	equipment1.SetRefId(1101)
	equipment1.SetDef(6)
	equipment1.SetHp(12)
	player.Equipments().Put(equipment1.Id(), equipment1)
	player.Items().Put(2001, 10)
	player.Items().Put(2002, 1)
	player.Cash().Stages().Put(1, 2)
	player.Cash().Stages().Put(2, 1)
	player.Cash().SetCards([]int{1, 2})
	player.Cash().SetOrderIds([]string{"order-0", "order-1"})
	player.SetUpdateVersion(1)
	player.SetCreateTime(now)
	player.SetUpdateTime(now)
	player.Reset()

	doc := player.ToDocument()
	if doc == nil {
		t.Error("The value expected not be nil")
		return
	}
	if 8 != len(doc) {
		t.Errorf("The value expected <%v> but was <%v>", 8, len(doc))
	}
	if 123 != doc[BnamePlayerUid] {
		t.Errorf("The value expected <%v> but was <%v>", 1, doc[BnamePlayerUid])
	}
	if 1 != doc[BnamePlayerUpdateVersion] {
		t.Errorf("The value expected <%v> but was <%v>", 1, doc[BnamePlayerUpdateVersion])
	}
	if primitive.NewDateTimeFromTime(now) != doc[BnamePlayerCreateTime] {
		t.Errorf("The value expected <%v> but was <%v>", primitive.NewDateTimeFromTime(now), doc[BnamePlayerCreateTime])
	}
	if primitive.NewDateTimeFromTime(now) != doc[BnamePlayerUpdateTime] {
		t.Errorf("The value expected <%v> but was <%v>", primitive.NewDateTimeFromTime(now), doc[BnamePlayerUpdateTime])
	}
	wlt := doc[BnamePlayerWallet]
	if wlt == nil {
		t.Error("The value expected not be nil")
	} else {
		if 5000 != wlt.(bson.M)[BnameWalletCoinTotal] {
			t.Errorf("The value expected <%v> but was <%v>", 5000, wlt.(bson.M)[BnameWalletCoinTotal])
		}
		if 2000 != wlt.(bson.M)[BnameWalletCoinUsed] {
			t.Errorf("The value expected <%v> but was <%v>", 2000, wlt.(bson.M)[BnameWalletCoinUsed])
		}
		if 10 != wlt.(bson.M)[BnameWalletDiamond] {
			t.Errorf("The value expected <%v> but was <%v>", 10, wlt.(bson.M)[BnameWalletDiamond])
		}
	}
	eqm := doc[BnamePlayerEquipments]
	if eqm == nil {
		t.Error("The value expected not be nil")
	} else {
		if 2 != len(eqm.(bson.M)) {
			t.Errorf("The value expected <%v> but was <%v>", 2, len(eqm.(bson.M)))
		}
		eq0 := eqm.(bson.M)["12345678-1234-5678-9abc-123456789abc"]
		if eq0 == nil {
			t.Error("The value expected not be nil")
		} else {
			if "12345678-1234-5678-9abc-123456789abc" != eq0.(bson.M)[BnameEquipmentId] {
				t.Errorf("The value expected <%v> but was <%v>", "12345678-1234-5678-9abc-123456789abc", eq0.(bson.M)[BnameEquipmentId])
			}
			if 1001 != eq0.(bson.M)[BnameEquipmentRefId] {
				t.Errorf("The value expected <%v> but was <%v>", 1001, eq0.(bson.M)[BnameEquipmentRefId])
			}
			if 12 != eq0.(bson.M)[BnameEquipmentAtk] {
				t.Errorf("The value expected <%v> but was <%v>", 12, eq0.(bson.M)[BnameEquipmentAtk])
			}
			if 0 != eq0.(bson.M)[BnameEquipmentDef] {
				t.Errorf("The value expected <%v> but was <%v>", 0, eq0.(bson.M)[BnameEquipmentDef])
			}
			if 0 != eq0.(bson.M)[BnameEquipmentHp] {
				t.Errorf("The value expected <%v> but was <%v>", 0, eq0.(bson.M)[BnameEquipmentHp])
			}
		}
		eq1 := eqm.(bson.M)["11111111-1111-1111-1111-111111111111"]
		if eq1 == nil {
			t.Error("The value expected not be nil")
		} else {
			if "11111111-1111-1111-1111-111111111111" != eq1.(bson.M)[BnameEquipmentId] {
				t.Errorf("The value expected <%v> but was <%v>", "11111111-1111-1111-1111-111111111111", eq0.(bson.M)[BnameEquipmentId])
			}
			if 1101 != eq1.(bson.M)[BnameEquipmentRefId] {
				t.Errorf("The value expected <%v> but was <%v>", 1101, eq1.(bson.M)[BnameEquipmentRefId])
			}
			if 0 != eq1.(bson.M)[BnameEquipmentAtk] {
				t.Errorf("The value expected <%v> but was <%v>", 0, eq1.(bson.M)[BnameEquipmentAtk])
			}
			if 6 != eq1.(bson.M)[BnameEquipmentDef] {
				t.Errorf("The value expected <%v> but was <%v>", 6, eq1.(bson.M)[BnameEquipmentDef])
			}
			if 12 != eq1.(bson.M)[BnameEquipmentHp] {
				t.Errorf("The value expected <%v> but was <%v>", 12, eq1.(bson.M)[BnameEquipmentHp])
			}
		}
		eqn := eqm.(bson.M)["no such equipment"]
		if eqn != nil {
			t.Errorf("The value expected nil but was <%v>", eqn)
		}
	}
	itm := doc[BnamePlayerItems]
	if itm == nil {
		t.Error("The value expected not be nil")
	} else {
		if 2 != len(itm.(bson.M)) {
			t.Errorf("The value expected <%v> but was <%v>", 2, len(itm.(bson.M)))
		}
		if 10 != itm.(bson.M)["2001"] {
			t.Errorf("The value expected <%v> but was <%v>", 10, itm.(bson.M)["2001"])
		}
		if 1 != itm.(bson.M)["2002"] {
			t.Errorf("The value expected <%v> but was <%v>", 1, itm.(bson.M)["2002"])
		}
	}
	cs := doc[BnamePlayerCash]
	if cs == nil {
		t.Error("The value expected not be nil")
		return
	}
	if 3 != len(cs.(bson.M)) {
		t.Errorf("The value expected <%v> but was <%v>", 3, len(cs.(bson.M)))
	}
	stg := cs.(bson.M)[BnameCashInfoStages]
	if stg == nil {
		t.Error("The value expected not be nil")
	} else {
		if 2 != len(stg.(bson.M)) {
			t.Errorf("The value expected <%v> but was <%v>", 2, len(stg.(bson.M)))
		}
		if 2 != stg.(bson.M)["1"] {
			t.Errorf("The value expected <%v> but was <%v>", 2, stg.(bson.M)["1"])
		}
		if 1 != stg.(bson.M)["2"] {
			t.Errorf("The value expected <%v> but was <%v>", 1, stg.(bson.M)["2"])
		}
	}
	cds := cs.(bson.M)[BnameCashInfoCards]
	if cds == nil {
		t.Error("The value expected not be nil")
	} else {
		if 2 != len(cds.(bson.A)) {
			t.Errorf("The value expected <%v> but was <%v>", 2, len(cds.(bson.A)))
		} else {
			if 1 != cds.(bson.A)[0] {
				t.Errorf("The value expected <%v> but was <%v>", 1, cds.(bson.A)[0])
			}
			if 2 != cds.(bson.A)[1] {
				t.Errorf("The value expected <%v> but was <%v>", 2, cds.(bson.A)[1])
			}
		}
	}
	ois := cs.(bson.M)[BnameCashInfoOrderIds]
	if ois == nil {
		t.Error("The value expected not be nil")
	} else {
		if 2 != len(ois.(bson.A)) {
			t.Errorf("The value expected <%v> but was <%v>", 2, len(ois.(bson.A)))
		} else {
			if "order-0" != ois.(bson.A)[0] {
				t.Errorf("The value expected <%v> but was <%v>", "order-0", ois.(bson.A)[0])
			}
			if "order-1" != ois.(bson.A)[1] {
				t.Errorf("The value expected <%v> but was <%v>", "order-1", ois.(bson.A)[1])
			}
		}
	}
}
