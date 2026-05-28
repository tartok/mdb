package mdb

import (
	"context"
	"strings"
	"sync"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Mdb struct {
	client *mongo.Client
	dbName string
}

var dbList = map[string]*Mdb{}
var dbMutex sync.Mutex

func Connect(ctx context.Context, url string, auth *options.Credential, dbName string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	instance := ""
	ss := strings.SplitN(dbName, ".", 2)
	if len(ss) > 1 {
		instance = ss[0]
		dbName = ss[1]
	}
	db, ok := dbList[instance]
	if !ok {
		db = &Mdb{
			dbName: dbName,
		}
		dbList[instance] = db
	}
	if db.client != nil {
		_ = db.client.Disconnect(ctx)
		db.client = nil
	}
	clientOptions := options.Client().ApplyURI(url)
	if auth != nil {
		clientOptions = clientOptions.SetAuth(*auth)
	}

	registry := bson.NewRegistry()
	//registry.RegisterTypeMapEntry(bson.TypeEmbeddedDocument, reflect.TypeOf(primitive.M{}))
	clientOptions.SetRegistry(registry)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}
	db.client = client
	return nil
}
func Coll(dbName string, collName Collection) *mongo.Collection {
	db := dbList[""]
	if db != nil && dbName == "" {
		dbName = db.dbName
	}
	if dbName != "" {
		ss := strings.SplitN(dbName, ".", 2)
		if len(ss) > 1 {
			db = dbList[ss[0]]
			dbName = ss[1]
		}
	}
	if db == nil || db.client == nil {
		return nil
	}
	return db.client.Database(dbName).Collection(string(collName))
}
