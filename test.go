package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bits-and-blooms/bitset"
	"github.com/fmjsjx/bson-model-go/bsonmodel"
	jsoniter "github.com/json-iterator/go"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Player struct {
}

func main() {
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
	doc := make(map[int]int)
	for i := 1; i < 10; i++ {
		doc[i] = i * i
	}
	for k, v := range doc {
		fmt.Printf("k: %d, v: %d\n", k, v)
	}
	delete(doc, 1)

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

	b := &bitset.BitSet{}
	fmt.Printf("b.Len() => %v\n", b.Len())
	fmt.Printf("b => %v\n", b)
	fmt.Printf("b.0 => %v\n", b.Test(0))
	b.Set(0)
	fmt.Printf("b.Len() => %v\n", b.Len())
	fmt.Printf("b => %v\n", b)
	fmt.Printf("b.0 => %v\n", b.Test(0))
	b.DeleteAt(0)
	fmt.Printf("b.Len() => %v\n", b.Len())
	fmt.Printf("b => %v\n", b)
	fmt.Printf("b.0 => %v\n", b.Test(0))
	b.Set(2)
	b.Set(3)
	fmt.Printf("b.Len() => %v\n", b.Len())
	fmt.Printf("b => %v\n", b)
	b.ClearAll()
	fmt.Printf("b.None() => %v\n", b.None())
	fmt.Printf("b.Len() => %v\n", b.Len())
	fmt.Printf("b => %v\n", b)
	fmt.Printf("b.Len() => %v\n", b.Len())
	b.DeleteAt(0)
	fmt.Printf("b.Len() => %v\n", b.Len())
	iii := bson.M{"cu": int32(1), "xxx": "yyy"}
	coinTotal := iii["ct"]
	if coinTotal == nil {
		fmt.Println(0)
	}

}
