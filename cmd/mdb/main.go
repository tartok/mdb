package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/tartok/mdb/pkg/mdb"
	"github.com/tartok/mdb/pkg/mdbs"
	"go.mongodb.org/mongo-driver/v2/bson"
	"gopkg.in/yaml.v3"
)

func main() {
	ctx := struct {
		mdb.Context
	}{}
	ctx.LoginId = mdb.Ref(bson.NewObjectID())
	//{
	//	c, cancel := context.WithCancel(context.Background())
	//	cancel()
	//	ctx.Context = *mdb.NewContext(c)
	//}
	f, err := os.Open("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	r := yaml.NewDecoder(f)
	var conf map[string]mdb.Config
	err = r.Decode(&conf)
	if err != nil {
		log.Fatal(err)
	}

	{
		err = mdb.Connect(ctx.Ctx, conf["auth"].Url, conf["auth"].Credential, conf["auth"].DbName)
		if err != nil {
			log.Fatal(err)
		}

		testCollection := mdb.New("auth.ede_auth", "sysUsers")
		var res []struct {
			Name string `bson:"name"`
		}
		err := testCollection.ReadMany(ctx, &res)
		fmt.Println(res, err)
	}
	err = mdb.Connect(ctx.Ctx, conf["default"].Url, conf["default"].Credential, conf["default"].DbName)
	if err != nil {
		log.Fatal(err)
	}

	testCollection := mdb.New("", "cTest")
	id, err := testCollection.Create(ctx, struct {
		Name string
	}{Name: "John"})
	fmt.Println(id, err)
	//err = testCollection.Update(ctx, mdb.FilterId(*id), mdb.UpdateStruct{Set: bson.M{"text": "text2"}}, nil)
	//if err != nil {
	//	log.Fatal(err)
	//}
	var res []interface{}
	err = testCollection.ReadMany(ctx, &res, mdbs.Match(bson.M{"text": "text2"}))
	fmt.Println(res, err)
	b, err := json.Marshal(res)
	fmt.Println(string(b), err)
}
