package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/fmjsjx/bson-model-go/bsonmodel"
	"github.com/fmjsjx/bson-model-go/example"
	jsoniter "github.com/json-iterator/go"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Player struct {
}

func main() {
	tnow := time.Now()
	fmt.Printf("time now: %v\n", tnow)
	dzt := time.Unix(time.Now().Unix()-2*86400, 0)
	fmt.Printf("time dzt: %v\n", dzt)

	tdt := dzt.Sub(tnow)
	fmt.Printf("time tdt: %v\n", tdt)
	today := time.Date(tnow.Year(), tnow.Month(), tnow.Day(), 0, 0, 0, 0, tnow.Location())
	dzd := time.Date(dzt.Year(), dzt.Month(), dzt.Day(), 0, 0, 0, 0, dzt.Location())
	tdd := today.Sub(dzd)
	fmt.Printf("time tdd: %v\n", tdd)
	tddays := int(tdd / time.Hour / 24)
	fmt.Printf("time tddays: %v\n", tddays)

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://10.7.125.140:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	databases, err := client.ListDatabaseNames(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Databases: %v\n", databases)

	collection := client.Database("cowboy").Collection("player")
	var result bson.M
	err = collection.FindOne(ctx, bson.D{{Key: "_id", Value: 1}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return
		}
		log.Fatal(err)
	}
	fmt.Printf("found document: %v\n", result)
	ct, ok := result["_ct"].(primitive.DateTime)
	if ok {
		createTime := ct.Time()
		fmt.Printf("createTime: %v\n", createTime)
		fmt.Printf("createTime.Unix: %d\n", createTime.Unix())
		fmt.Printf("createTime.UnixMilli: %d\n", createTime.UnixMilli())
	}

	imap := bsonmodel.NewIntSimpleMapModel(nil, "test", bsonmodel.IntValueType())
	imap.Put(1, 101)
	imap.Put(2, 102)
	imap.Put(3, 103)
	fmt.Printf("imap: %v\n", imap)
	// api := jsoniter.ConfigCompatibleWithStandardLibrary
	json, err := jsoniter.MarshalToString(imap)
	if err == nil {
		fmt.Printf("imap: %s\n", json)
	} else {
		fmt.Printf("failed: %e\n", err)
	}

	smap := bsonmodel.NewStringSimpleMapModel(nil, "test2", bsonmodel.IntValueType())
	smap.Put("a", 101)
	smap.Put("b", 102)
	smap.Put("c", 103)
	fmt.Printf("smap: %v\n", smap)
	json, err = jsoniter.MarshalToString(smap)
	if err == nil {
		fmt.Printf("smap: %s\n", json)
	} else {
		fmt.Printf("failed: %e\n", err)
	}
	now := time.Now()
	doc := bson.M{"_id": int32(123),
		"wlt": bson.M{"ct": int32(5000), "cu": int32(2000), "d": int32(10)},
		"eqm": bson.M{"12345678-1234-5678-9abc-123456789abc": bson.M{"id": "12345678-1234-5678-9abc-123456789abc", "rid": int32(1), "atk": int32(12), "def": int32(2), "hp": int32(100)}},
		"itm": bson.M{"2001": int32(10), "2002": int32(1)},
		"cs":  bson.M{"stg": bson.M{"1": int32(1), "2": int32(1)}, "cs": bson.A{int32(1), int32(2)}},
		"_uv": int32(1), "_ct": primitive.NewDateTimeFromTime(now), "_ut": primitive.NewDateTimeFromTime(now)}
	player, err := example.LoadPlayerFromDocument(doc)
	if err != nil {
		fmt.Printf("Load player failed: %e\n", err)
	}
	fmt.Printf("Player => %v\n", player)
	json, err = player.MarshalToJsonString()
	if err == nil {
		fmt.Printf("player: %s\n", json)
	} else {
		fmt.Printf("failed: %e\n", err)
	}
	fmt.Printf("update => %v\n", player.ToUpdate())
	player.Wallet().SetCoinUsed(2200)
	player.Wallet().SetDiamond(11)
	equipment := player.Equipment("12345678-1234-5678-9abc-123456789abc")
	equipment.SetAtk(13)
	player.Items().Put(3001, 1)
	player.Items().Remove(2002)
	player.Cash().Stages().Put(3, 1)
	player.Cash().SetCards([]int{1, 2, 3})
	player.IncreaseUpdateVersion()
	player.SetUpdateTime(time.Now())
	fmt.Printf("update => %v\n", player.ToUpdate())
	equipment.SetFullyUpdate(true)
	fmt.Printf("update => %v\n", player.ToUpdate())
	json, err = player.ToSyncJson()
	if err == nil {
		fmt.Printf("player sync data: %s\n", json)
	} else {
		fmt.Printf("failed: %e\n", err)
	}
	json, err = player.ToDeleteJson()
	if err == nil {
		fmt.Printf("player delete data: %s\n", json)
	} else {
		fmt.Printf("failed: %e\n", err)
	}
	player.Reset()
	fmt.Printf("update => %v\n", player.ToUpdate())

	json, err = player.MarshalToJsonString()
	if err == nil {
		fmt.Printf("player: %s\n", json)
	} else {
		fmt.Printf("failed: %e\n", err)
	}

	json, err = player.ToDataJson()
	if err == nil {
		fmt.Printf("player data: %s\n", json)
	} else {
		fmt.Printf("failed: %e\n", err)
	}
	any := jsoniter.Get([]byte(json))
	fmt.Printf("any => %s\n", any.ToString())
	player, err = example.LoadPlayerFromJsoniter(any)
	if err != nil {
		fmt.Printf("Load player failed: %e\n", err)
	}
	json, err = player.MarshalToJsonString()
	if err == nil {
		fmt.Printf("player: %s\n", json)
	} else {
		fmt.Printf("failed: %e\n", err)
	}
}
