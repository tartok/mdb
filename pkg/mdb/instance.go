package mdb

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Mdb struct {
	client *mongo.Client
	dbName string
}

var Db *Mdb

func Connect(ctx context.Context, url string, auth *options.Credential, dbName string) error {
	if Db != nil {
		if Db.client != nil {
			_ = Db.client.Disconnect(ctx)
		}
		Db = nil
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
	Db = &Mdb{
		client: client,
		dbName: dbName,
	}
	return nil
}
func Coll(dbName string, collName Collection) *mongo.Collection {
	if dbName == "" {
		dbName = Db.dbName
	}
	return Db.client.Database(dbName).Collection(string(collName))
}
