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
	if doc[BnamePlayerWallet] == nil {
		t.Error("The value expected not be nil")
	} else {
		wlt := doc[BnamePlayerWallet].(bson.M)
		if 5000 != wlt[BnameWalletCoinTotal] {
			t.Errorf("The value expected <%v> but was <%v>", 5000, wlt[BnameWalletCoinTotal])
		}
		if 2000 != wlt[BnameWalletCoinUsed] {
			t.Errorf("The value expected <%v> but was <%v>", 2000, wlt[BnameWalletCoinUsed])
		}
		if 10 != wlt[BnameWalletDiamond] {
			t.Errorf("The value expected <%v> but was <%v>", 10, wlt[BnameWalletDiamond])
		}
	}
	if doc[BnamePlayerEquipments] == nil {
		t.Error("The value expected not be nil")
	} else {
		eqm := doc[BnamePlayerEquipments].(bson.M)
		if 2 != len(eqm) {
			t.Errorf("The value expected <%v> but was <%v>", 2, len(eqm))
		}
		if eqm["12345678-1234-5678-9abc-123456789abc"] == nil {
			t.Error("The value expected not be nil")
		} else {
			eq0 := eqm["12345678-1234-5678-9abc-123456789abc"].(bson.M)
			if "12345678-1234-5678-9abc-123456789abc" != eq0[BnameEquipmentId] {
				t.Errorf("The value expected <%v> but was <%v>", "12345678-1234-5678-9abc-123456789abc", eq0[BnameEquipmentId])
			}
			if 1001 != eq0[BnameEquipmentRefId] {
				t.Errorf("The value expected <%v> but was <%v>", 1001, eq0[BnameEquipmentRefId])
			}
			if 12 != eq0[BnameEquipmentAtk] {
				t.Errorf("The value expected <%v> but was <%v>", 12, eq0[BnameEquipmentAtk])
			}
			if 0 != eq0[BnameEquipmentDef] {
				t.Errorf("The value expected <%v> but was <%v>", 0, eq0[BnameEquipmentDef])
			}
			if 0 != eq0[BnameEquipmentHp] {
				t.Errorf("The value expected <%v> but was <%v>", 0, eq0[BnameEquipmentHp])
			}
		}
		if eqm["11111111-1111-1111-1111-111111111111"] == nil {
			t.Error("The value expected not be nil")
		} else {
			eq1 := eqm["11111111-1111-1111-1111-111111111111"].(bson.M)
			if "11111111-1111-1111-1111-111111111111" != eq1[BnameEquipmentId] {
				t.Errorf("The value expected <%v> but was <%v>", "11111111-1111-1111-1111-111111111111", eq1[BnameEquipmentId])
			}
			if 1101 != eq1[BnameEquipmentRefId] {
				t.Errorf("The value expected <%v> but was <%v>", 1101, eq1[BnameEquipmentRefId])
			}
			if 0 != eq1[BnameEquipmentAtk] {
				t.Errorf("The value expected <%v> but was <%v>", 0, eq1[BnameEquipmentAtk])
			}
			if 6 != eq1[BnameEquipmentDef] {
				t.Errorf("The value expected <%v> but was <%v>", 6, eq1[BnameEquipmentDef])
			}
			if 12 != eq1[BnameEquipmentHp] {
				t.Errorf("The value expected <%v> but was <%v>", 12, eq1[BnameEquipmentHp])
			}
		}
		eqn := eqm["no such equipment"]
		if eqn != nil {
			t.Errorf("The value expected nil but was <%v>", eqn)
		}
	}

	if doc[BnamePlayerItems] == nil {
		t.Error("The value expected not be nil")
	} else {
		itm := doc[BnamePlayerItems].(bson.M)
		if 2 != len(itm) {
			t.Errorf("The value expected <%v> but was <%v>", 2, len(itm))
		}
		if 10 != itm["2001"] {
			t.Errorf("The value expected <%v> but was <%v>", 10, itm["2001"])
		}
		if 1 != itm["2002"] {
			t.Errorf("The value expected <%v> but was <%v>", 1, itm["2002"])
		}
	}

	if doc[BnamePlayerCash] == nil {
		t.Error("The value expected not be nil")
		return
	}
	cs := doc[BnamePlayerCash].(bson.M)
	if 3 != len(cs) {
		t.Errorf("The value expected <%v> but was <%v>", 3, len(cs))
	}
	if cs[BnameCashInfoStages] == nil {
		t.Error("The value expected not be nil")
	} else {
		stg := cs[BnameCashInfoStages].(bson.M)
		if 2 != len(stg) {
			t.Errorf("The value expected <%v> but was <%v>", 2, len(stg))
		}
		if 2 != stg["1"] {
			t.Errorf("The value expected <%v> but was <%v>", 2, stg["1"])
		}
		if 1 != stg["2"] {
			t.Errorf("The value expected <%v> but was <%v>", 1, stg["2"])
		}
	}
	if cs[BnameCashInfoCards] == nil {
		t.Error("The value expected not be nil")
	} else {
		cds := cs[BnameCashInfoCards].(bson.A)
		if 2 != len(cds) {
			t.Errorf("The value expected <%v> but was <%v>", 2, len(cds))
		} else {
			if 1 != cds[0] {
				t.Errorf("The value expected <%v> but was <%v>", 1, cds[0])
			}
			if 2 != cds[1] {
				t.Errorf("The value expected <%v> but was <%v>", 2, cds[1])
			}
		}
	}
	if cs[BnameCashInfoOrderIds] == nil {
		t.Error("The value expected not be nil")
	} else {
		ois := cs[BnameCashInfoOrderIds].(bson.A)
		if 2 != len(ois) {
			t.Errorf("The value expected <%v> but was <%v>", 2, len(ois))
		} else {
			if "order-0" != ois[0] {
				t.Errorf("The value expected <%v> but was <%v>", "order-0", ois[0])
			}
			if "order-1" != ois[1] {
				t.Errorf("The value expected <%v> but was <%v>", "order-1", ois[1])
			}
		}
	}
}

func TestToData(t *testing.T) {
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

	data := player.ToData().(map[string]interface{})
	if 8 != len(data) {
		t.Errorf("The value expected <%v> but was <%v>", 8, len(data))
	}
	if 123 != data[BnamePlayerUid] {
		t.Errorf("The value expected <%v> but was <%v>", 1, data[BnamePlayerUid])
	}
	if 1 != data[BnamePlayerUpdateVersion] {
		t.Errorf("The value expected <%v> but was <%v>", 1, data[BnamePlayerUpdateVersion])
	}
	if now.UnixMilli() != data[BnamePlayerCreateTime] {
		t.Errorf("The value expected <%v> but was <%v>", now.UnixMilli(), data[BnamePlayerCreateTime])
	}
	if now.UnixMilli() != data[BnamePlayerUpdateTime] {
		t.Errorf("The value expected <%v> but was <%v>", now.UnixMilli(), data[BnamePlayerUpdateTime])
	}
	if data[BnamePlayerWallet] == nil {
		t.Error("The value expected not be nil")
	} else {
		wlt := data[BnamePlayerWallet].(map[string]interface{})
		if 3 != len(wlt) {
			t.Errorf("The value expected <%v> but was <%v>", 3, len(wlt))
		}
		if 5000 != wlt[BnameWalletCoinTotal] {
			t.Errorf("The value expected <%v> but was <%v>", 5000, wlt[BnameWalletCoinTotal])
		}
		if 2000 != wlt[BnameWalletCoinUsed] {
			t.Errorf("The value expected <%v> but was <%v>", 2000, wlt[BnameWalletCoinUsed])
		}
		if 10 != wlt[BnameWalletDiamond] {
			t.Errorf("The value expected <%v> but was <%v>", 10, wlt[BnameWalletDiamond])
		}
	}
	if data[BnamePlayerEquipments] == nil {
		t.Error("The value expected not be nil")
	} else {
		eqm := data[BnamePlayerEquipments].(map[string]interface{})
		if 2 != len(eqm) {
			t.Errorf("The value expected <%v> but was <%v>", 2, len(eqm))
		}
		if eqm["12345678-1234-5678-9abc-123456789abc"] == nil {
			t.Error("The value expected not be nil")
		} else {
			eq0 := eqm["12345678-1234-5678-9abc-123456789abc"].(map[string]interface{})
			if "12345678-1234-5678-9abc-123456789abc" != eq0[BnameEquipmentId] {
				t.Errorf("The value expected <%v> but was <%v>", "12345678-1234-5678-9abc-123456789abc", eq0[BnameEquipmentId])
			}
			if 1001 != eq0[BnameEquipmentRefId] {
				t.Errorf("The value expected <%v> but was <%v>", 1001, eq0[BnameEquipmentRefId])
			}
			if 12 != eq0[BnameEquipmentAtk] {
				t.Errorf("The value expected <%v> but was <%v>", 12, eq0[BnameEquipmentAtk])
			}
			if 0 != eq0[BnameEquipmentDef] {
				t.Errorf("The value expected <%v> but was <%v>", 0, eq0[BnameEquipmentDef])
			}
			if 0 != eq0[BnameEquipmentHp] {
				t.Errorf("The value expected <%v> but was <%v>", 0, eq0[BnameEquipmentHp])
			}
		}
		if eqm["11111111-1111-1111-1111-111111111111"] == nil {
			t.Error("The value expected not be nil")
		} else {
			eq1 := eqm["11111111-1111-1111-1111-111111111111"].(map[string]interface{})
			if "11111111-1111-1111-1111-111111111111" != eq1[BnameEquipmentId] {
				t.Errorf("The value expected <%v> but was <%v>", "11111111-1111-1111-1111-111111111111", eq1[BnameEquipmentId])
			}
			if 1101 != eq1[BnameEquipmentRefId] {
				t.Errorf("The value expected <%v> but was <%v>", 1101, eq1[BnameEquipmentRefId])
			}
			if 0 != eq1[BnameEquipmentAtk] {
				t.Errorf("The value expected <%v> but was <%v>", 0, eq1[BnameEquipmentAtk])
			}
			if 6 != eq1[BnameEquipmentDef] {
				t.Errorf("The value expected <%v> but was <%v>", 6, eq1[BnameEquipmentDef])
			}
			if 12 != eq1[BnameEquipmentHp] {
				t.Errorf("The value expected <%v> but was <%v>", 12, eq1[BnameEquipmentHp])
			}
		}
		if eqm["no such equipment"] != nil {
			t.Errorf("The value expected nil but was <%v>", eqm["no such equipment"])
		}
	}
	if data[BnamePlayerItems] == nil {
		t.Error("The value expected not be nil")
	} else {
		itm := data[BnamePlayerItems].(map[int]interface{})
		if 2 != len(itm) {
			t.Errorf("The value expected <%v> but was <%v>", 2, len(itm))
		}
		if 10 != itm[2001] {
			t.Errorf("The value expected <%v> but was <%v>", 10, itm[2001])
		}
		if 1 != itm[2002] {
			t.Errorf("The value expected <%v> but was <%v>", 1, itm[2002])
		}
	}
	cs := data[BnamePlayerCash].(map[string]interface{})
	if 3 != len(cs) {
		t.Errorf("The value expected <%v> but was <%v>", 3, len(cs))
	}
	if cs[BnameCashInfoStages] == nil {
		t.Error("The value expected not be nil")
	} else {
		stg := cs[BnameCashInfoStages].(map[int]interface{})
		if 2 != len(stg) {
			t.Errorf("The value expected <%v> but was <%v>", 2, len(stg))
		}
		if 2 != stg[1] {
			t.Errorf("The value expected <%v> but was <%v>", 2, stg[1])
		}
		if 1 != stg[2] {
			t.Errorf("The value expected <%v> but was <%v>", 1, stg[2])
		}
	}
	if cs[BnameCashInfoCards] == nil {
		t.Error("The value expected not be nil")
	} else {
		cds := cs[BnameCashInfoCards].([]int)
		if 2 != len(cds) {
			t.Errorf("The value expected <%v> but was <%v>", 2, len(cds))
		} else {
			if 1 != cds[0] {
				t.Errorf("The value expected <%v> but was <%v>", 1, cds[0])
			}
			if 2 != cds[1] {
				t.Errorf("The value expected <%v> but was <%v>", 2, cds[1])
			}
		}
	}
	if cs[BnameCashInfoOrderIds] == nil {
		t.Error("The value expected not be nil")
	} else {
		ois := cs[BnameCashInfoOrderIds].([]string)
		if 2 != len(ois) {
			t.Errorf("The value expected <%v> but was <%v>", 2, len(ois))
		} else {
			if "order-0" != ois[0] {
				t.Errorf("The value expected <%v> but was <%v>", "order-0", ois[0])
			}
			if "order-1" != ois[1] {
				t.Errorf("The value expected <%v> but was <%v>", "order-1", ois[1])
			}
		}
	}
}

func TestToUpdate(t *testing.T) {
	createTime := time.Now().Add(-1 * time.Hour).Truncate(time.Millisecond)
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
	player.SetCreateTime(createTime)
	player.SetUpdateTime(createTime)
	player.Reset()

	player.Wallet().SetCoinTotal(5200)
	player.Wallet().SetCoinUsed(2100)
	player.Wallet().SetDiamond(11)
	player.Equipments().Remove("12345678-1234-5678-9abc-123456789abc")
	player.Equipment("11111111-1111-1111-1111-111111111111").SetHp(20)
	equipment2 := NewEquipment()
	equipment2.SetId("22222222-2222-2222-2222-222222222222")
	equipment2.SetRefId(1201)
	equipment2.SetAtk(2)
	equipment2.SetDef(2)
	equipment2.SetHp(2)
	player.Equipments().Put(equipment2.Id(), equipment2)
	player.Items().Put(2001, 12)
	player.Items().Put(2002, 0)
	player.Items().Put(2003, 1)
	player.Cash().Stages().Remove(1)
	player.Cash().SetCards(nil)
	player.Cash().SetOrderIds([]string{"order-0", "order-1", "order-2"})
	player.IncreaseUpdateVersion()
	now := time.Now().Truncate(time.Millisecond)
	player.SetUpdateTime(now)

	update := player.ToUpdate()
	if 2 != len(update) {
		t.Errorf("The value expected <%v> but was <%v>", 2, len(update))
	}
	if update["$set"] == nil {
		t.Error("The value expected not be nil")
		t.FailNow()
	}
	dset := update["$set"].(bson.M)
	if 11 != len(dset) {
		t.Errorf("The value expected <%v> but was <%v>", 11, len(dset))
	}
	if 5200 != dset["wlt.ct"] {
		t.Errorf("The value expected <%v> but was <%v>", 5200, dset["wlt.ct"])
	}
	if 2100 != dset["wlt.cu"] {
		t.Errorf("The value expected <%v> but was <%v>", 2100, dset["wlt.cu"])
	}
	if 11 != dset["wlt.d"] {
		t.Errorf("The value expected <%v> but was <%v>", 11, dset["wlt.d"])
	}
	if 20 != dset["eqm.11111111-1111-1111-1111-111111111111.hp"] {
		t.Errorf("The value expected <%v> but was <%v>", 20, dset["eqm.11111111-1111-1111-1111-111111111111.hp"])
	}
	if dset["eqm.22222222-2222-2222-2222-222222222222"] == nil {
		t.Error("The value expected not be nil")
	} else {
		eq2 := dset["eqm.22222222-2222-2222-2222-222222222222"].(bson.M)
		if "22222222-2222-2222-2222-222222222222" != eq2["id"] {
			t.Errorf("The value expected <%v> but was <%v>", "22222222-2222-2222-2222-222222222222", eq2["id"])
		}
		if 1201 != eq2["rid"] {
			t.Errorf("The value expected <%v> but was <%v>", 1201, eq2["rid"])
		}
		if 2 != eq2["atk"] {
			t.Errorf("The value expected <%v> but was <%v>", 2, eq2["atk"])
		}
		if 2 != eq2["def"] {
			t.Errorf("The value expected <%v> but was <%v>", 2, eq2["def"])
		}
		if 2 != eq2["hp"] {
			t.Errorf("The value expected <%v> but was <%v>", 2, eq2["hp"])
		}
	}
	if 12 != dset["itm.2001"] {
		t.Errorf("The value expected <%v> but was <%v>", 12, dset["itm.2001"])
	}
	if 0 != dset["itm.2002"] {
		t.Errorf("The value expected <%v> but was <%v>", 0, dset["itm.2002"])
	}
	if 1 != dset["itm.2003"] {
		t.Errorf("The value expected <%v> but was <%v>", 1, dset["itm.2003"])
	}
	if dset["cs.ois"] == nil {
		t.Error("The value expected not be nil")
	} else {
		ois := dset["cs.ois"].(bson.A)
		if 3 != len(ois) {
			t.Errorf("The value expected <%v> but was <%v>", 3, len(ois))
		}
		if "order-0" != ois[0] {
			t.Errorf("The value expected <%v> but was <%v>", "order-0", ois[0])
		}
		if "order-1" != ois[1] {
			t.Errorf("The value expected <%v> but was <%v>", "order-1", ois[1])
		}
		if "order-2" != ois[2] {
			t.Errorf("The value expected <%v> but was <%v>", "order-2", ois[2])
		}
	}
	if 2 != dset["_uv"] {
		t.Errorf("The value expected <%v> but was <%v>", 2, dset["_uv"])
	}
	if primitive.NewDateTimeFromTime(now) != dset["_ut"] {
		t.Errorf("The value expected <%v> but was <%v>", primitive.NewDateTimeFromTime(now), dset["_ut"])
	}
	if update["$unset"] == nil {
		t.Error("The value expected not be nil")
		t.FailNow()
	}
	unset := update["$unset"].(bson.M)
	if 3 != len(unset) {
		t.Errorf("The value expected <%v> but was <%v>", 3, len(unset))
	}
	if "" != unset["eqm.12345678-1234-5678-9abc-123456789abc"] {
		t.Errorf("The value expected <%v> but was <%v>", "", unset["eqm.12345678-1234-5678-9abc-123456789abc"])
	}
	if "" != unset["cs.stg.1"] {
		t.Errorf("The value expected <%v> but was <%v>", "", unset["cs.stg.1"])
	}
	if "" != unset["cs.cs"] {
		t.Errorf("The value expected <%v> but was <%v>", "", unset["cs.cs"])
	}
}

func TestToSync(t *testing.T) {
	createTime := time.Now().Add(-1 * time.Hour).Truncate(time.Millisecond)
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
	player.SetCreateTime(createTime)
	player.SetUpdateTime(createTime)
	player.Reset()

	player.Wallet().SetCoinTotal(5200)
	player.Wallet().SetCoinUsed(2100)
	player.Wallet().SetDiamond(11)
	player.Equipments().Remove("12345678-1234-5678-9abc-123456789abc")
	player.Equipment("11111111-1111-1111-1111-111111111111").SetHp(20)
	equipment2 := NewEquipment()
	equipment2.SetId("22222222-2222-2222-2222-222222222222")
	equipment2.SetRefId(1201)
	equipment2.SetAtk(2)
	equipment2.SetDef(2)
	equipment2.SetHp(2)
	player.Equipments().Put(equipment2.Id(), equipment2)
	player.Items().Put(2001, 12)
	player.Items().Put(2002, 0)
	player.Items().Put(2003, 1)
	player.Cash().Stages().Remove(1)
	player.Cash().Stages().Put(3, 1)
	player.Cash().SetCards(nil)
	player.Cash().SetOrderIds([]string{"order-0", "order-1", "order-2"})
	player.IncreaseUpdateVersion()
	now := time.Now().Truncate(time.Millisecond)
	player.SetUpdateTime(now)

	sync := player.ToSync().(map[string]interface{})
	if 4 != len(sync) {
		t.Errorf("The value expected <%v> but was <%v>", 4, len(sync))
	}
	if sync["wallet"] == nil {
		t.Error("The value expected not be nil")
	} else {
		wallet := sync["wallet"].(map[string]interface{})
		if 3 != len(wallet) {
			t.Errorf("The value expected <%v> but was <%v>", 3, len(wallet))
		}
		if 5200 != wallet["coinTotal"] {
			t.Errorf("The value expected <%v> but was <%v>", 5200, wallet["coinTotal"])
		}
		if 3100 != wallet["coin"] {
			t.Errorf("The value expected <%v> but was <%v>", 3100, wallet["coin"])
		}
		if 11 != wallet["diamond"] {
			t.Errorf("The value expected <%v> but was <%v>", 11, wallet["diamond"])
		}
	}
	if sync["equipments"] == nil {
		t.Error("The value expected not be nil")
	} else {
		equipments := sync["equipments"].(map[string]interface{})
		if 2 != len(equipments) {
			t.Errorf("The value expected <%v> but was <%v>", 2, len(equipments))
		}
		if equipments["11111111-1111-1111-1111-111111111111"] == nil {
			t.Error("The value expected not be nil")
		} else {
			eq1 := equipments["11111111-1111-1111-1111-111111111111"].(map[string]interface{})
			if 1 != len(eq1) {
				t.Errorf("The value expected <%v> but was <%v>", 1, len(eq1))
			}
			if 20 != eq1["hp"] {
				t.Errorf("The value expected <%v> but was <%v>", 20, eq1["hp"])
			}
		}
		if equipments["22222222-2222-2222-2222-222222222222"] == nil {
			t.Error("The value expected not be nil")
		} else {
			eq2 := equipments["22222222-2222-2222-2222-222222222222"].(Equipment)
			if equipment2 != eq2 {
				t.Errorf("The value expected <%v> but was <%v>", equipment2, eq2)
			}
		}
	}
	if sync["items"] == nil {
		t.Error("The value expected not be nil")
	} else {
		items := sync["items"].(map[int]interface{})
		if 3 != len(items) {
			t.Errorf("The value expected <%v> but was <%v>", 3, len(items))
		}
		if 12 != items[2001] {
			t.Errorf("The value expected <%v> but was <%v>", 12, items[2001])
		}
		if 0 != items[2002] {
			t.Errorf("The value expected <%v> but was <%v>", 0, items[2002])
		}
		if 1 != items[2003] {
			t.Errorf("The value expected <%v> but was <%v>", 1, items[2003])
		}
	}
	if sync["cash"] == nil {
		t.Error("The value expected not be nil")
	} else {
		cash := sync["cash"].(map[string]interface{})
		if 2 != len(cash) {
			t.Errorf("The value expected <%v> but was <%v>", 2, len(cash))
		}
		if cash["stages"] == nil {
			t.Error("The value expected not be nil")
		} else {
			stages := cash["stages"].(map[int]interface{})
			if 1 != len(stages) {
				t.Errorf("The value expected <%v> but was <%v>", 1, len(stages))
			}
			if 1 != stages[3] {
				t.Errorf("The value expected <%v> but was <%v>", 1, stages[3])
			}
		}
		if cash["orderIds"] == nil {
			t.Error("The value expected not be nil")
		} else {
			orderIds := cash["orderIds"].([]string)
			if 3 != len(orderIds) {
				t.Errorf("The value expected <%v> but was <%v>", 3, len(orderIds))
			}
			if "order-0" != orderIds[0] {
				t.Errorf("The value expected <%v> but was <%v>", "order-0", orderIds[0])
			}
			if "order-1" != orderIds[1] {
				t.Errorf("The value expected <%v> but was <%v>", "order-1", orderIds[1])
			}
			if "order-2" != orderIds[2] {
				t.Errorf("The value expected <%v> but was <%v>", "order-2", orderIds[2])
			}
		}
	}
}

func TestToDelete(t *testing.T) {
	createTime := time.Now().Add(-1 * time.Hour).Truncate(time.Millisecond)
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
	player.SetCreateTime(createTime)
	player.SetUpdateTime(createTime)
	player.Reset()

	player.Wallet().SetCoinTotal(5200)
	player.Wallet().SetCoinUsed(2100)
	player.Wallet().SetDiamond(11)
	player.Equipments().Remove("12345678-1234-5678-9abc-123456789abc")
	player.Equipment("11111111-1111-1111-1111-111111111111").SetHp(20)
	equipment2 := NewEquipment()
	equipment2.SetId("22222222-2222-2222-2222-222222222222")
	equipment2.SetRefId(1201)
	equipment2.SetAtk(2)
	equipment2.SetDef(2)
	equipment2.SetHp(2)
	player.Equipments().Put(equipment2.Id(), equipment2)
	player.Items().Put(2001, 12)
	player.Items().Put(2002, 0)
	player.Items().Put(2003, 1)
	player.Cash().Stages().Remove(1)
	player.Cash().Stages().Put(3, 1)
	player.Cash().SetCards(nil)
	player.Cash().SetOrderIds([]string{"order-0", "order-1", "order-2"})
	player.IncreaseUpdateVersion()
	now := time.Now().Truncate(time.Millisecond)
	player.SetUpdateTime(now)

	delete := player.ToDelete().(map[string]interface{})
	if 2 != len(delete) {
		t.Errorf("The value expected <%v> but was <%v>", 2, len(delete))
	}
	if delete["equipments"] == nil {
		t.Error("The value expected not be nil")
	} else {
		equipments := delete["equipments"].(map[string]int)
		if 1 != len(equipments) {
			t.Errorf("The value expected <%v> but was <%v>", 1, len(equipments))
		}
		if 1 != equipments["12345678-1234-5678-9abc-123456789abc"] {
			t.Errorf("The value expected <%v> but was <%v>", 1, equipments["12345678-1234-5678-9abc-123456789abc"])
		}
	}
	if delete["cash"] == nil {
		t.Error("The value expected not be nil")
	} else {
		cash := delete["cash"].(map[string]interface{})
		if 2 != len(cash) {
			t.Errorf("The value expected <%v> but was <%v>", 2, len(cash))
		}
		if cash["stages"] == nil {
			t.Error("The value expected not be nil")
		} else {
			stages := cash["stages"].(map[int]int)
			if 1 != len(stages) {
				t.Errorf("The value expected <%v> but was <%v>", 1, len(stages))
			}
			if 1 != stages[1] {
				t.Errorf("The value expected <%v> but was <%v>", 1, stages[1])
			}
		}
		if 1 != cash["cards"] {
			t.Errorf("The value expected <%v> but was <%v>", 1, cash["cards"])
		}
	}
}
