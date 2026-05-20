package main

import (
	"context"
	"fmt"
	"log"
	"mdb/pkg/mdb"
	"os"

	"go.mongodb.org/mongo-driver/v2/bson"
	"gopkg.in/yaml.v3"
)

func main() {
	ctx := mdb.Context{
		Ctx:     context.Background(),
		LoginId: nil,
	}
	f, err := os.Open("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	r := yaml.NewDecoder(f)
	var conf mdb.Config
	err = r.Decode(&conf)
	if err != nil {
		log.Fatal(err)
	}

	testCollection := mdb.New("plant", "cTest")
	err = mdb.Connect(ctx.Ctx, conf.Url, conf.Credential)
	if err != nil {
		log.Fatal(err)
	}

	id, err := testCollection.Create(ctx, struct {
		Name string
	}{Name: "John"})
	fmt.Println(id, err)
	err = testCollection.Update(ctx, bson.M{"_id": id}, mdb.UpdateStruct{Set: bson.M{"text": "text"}}, nil)
	if err != nil {
		log.Fatal(err)
	}
	var res []interface{}
	err = testCollection.ReadMany(ctx, &res, nil)
	fmt.Println(r, err)
}
