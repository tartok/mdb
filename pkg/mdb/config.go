package mdb

import "go.mongodb.org/mongo-driver/v2/mongo/options"

type Config struct {
	Url        string              `json:"url" yaml:"url"`
	Credential *options.Credential `json:"credential" yaml:"credential"`
	DbName     string              `json:"dbName" yaml:"dbName"`
}
